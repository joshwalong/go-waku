package peer_exchange

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	libp2pProtocol "github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-msgio/pbio"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/waku-org/go-waku/logging"
	"github.com/waku-org/go-waku/waku/v2/discv5"
	"github.com/waku-org/go-waku/waku/v2/peermanager"
	"github.com/waku-org/go-waku/waku/v2/protocol"
	"github.com/waku-org/go-waku/waku/v2/protocol/enr"
	"github.com/waku-org/go-waku/waku/v2/protocol/peer_exchange/pb"
	"go.uber.org/zap"
)

// PeerExchangeID_v20alpha1 is the current Waku Peer Exchange protocol identifier
const PeerExchangeID_v20alpha1 = libp2pProtocol.ID("/vac/waku/peer-exchange/2.0.0-alpha1")
const MaxCacheSize = 1000

var (
	ErrNoPeersAvailable = errors.New("no suitable remote peers")
	ErrInvalidID        = errors.New("invalid request id")
)

// PeerConnector will subscribe to a channel containing the information for all peers found by this discovery protocol
type PeerConnector interface {
	Subscribe(context.Context, <-chan peermanager.PeerData)
}

type WakuPeerExchange struct {
	h       host.Host
	disc    *discv5.DiscoveryV5
	pm      *peermanager.PeerManager
	metrics Metrics
	log     *zap.Logger

	*protocol.CommonService

	peerConnector PeerConnector
	enrCache      *enrCache
}

// NewWakuPeerExchange returns a new instance of WakuPeerExchange struct
// Takes an optional peermanager if WakuPeerExchange is being created along with WakuNode.
// If using libp2p host, then pass peermanager as nil
func NewWakuPeerExchange(disc *discv5.DiscoveryV5, peerConnector PeerConnector, pm *peermanager.PeerManager, reg prometheus.Registerer, log *zap.Logger) (*WakuPeerExchange, error) {
	newEnrCache, err := newEnrCache(MaxCacheSize)
	if err != nil {
		return nil, err
	}
	wakuPX := new(WakuPeerExchange)
	wakuPX.disc = disc
	wakuPX.metrics = newMetrics(reg)
	wakuPX.log = log.Named("wakupx")
	wakuPX.enrCache = newEnrCache
	wakuPX.peerConnector = peerConnector
	wakuPX.pm = pm
	wakuPX.CommonService = protocol.NewCommonService()

	return wakuPX, nil
}

// SetHost sets the host to be able to mount or consume a protocol
func (wakuPX *WakuPeerExchange) SetHost(h host.Host) {
	wakuPX.h = h
}

// Start inits the peer exchange protocol
func (wakuPX *WakuPeerExchange) Start(ctx context.Context) error {
	return wakuPX.CommonService.Start(ctx, wakuPX.start)
}

func (wakuPX *WakuPeerExchange) start() error {
	wakuPX.h.SetStreamHandlerMatch(PeerExchangeID_v20alpha1, protocol.PrefixTextMatch(string(PeerExchangeID_v20alpha1)), wakuPX.onRequest())

	wakuPX.WaitGroup().Add(1)
	go wakuPX.runPeerExchangeDiscv5Loop(wakuPX.Context())
	wakuPX.log.Info("Peer exchange protocol started")
	return nil
}

func (wakuPX *WakuPeerExchange) onRequest() func(s network.Stream) {
	return func(s network.Stream) {
		defer s.Close()
		logger := wakuPX.log.With(logging.HostID("peer", s.Conn().RemotePeer()))
		requestRPC := &pb.PeerExchangeRPC{}
		reader := pbio.NewDelimitedReader(s, math.MaxInt32)
		err := reader.ReadMsg(requestRPC)
		if err != nil {
			logger.Error("reading request", zap.Error(err))
			wakuPX.metrics.RecordError(decodeRPCFailure)
			return
		}

		if requestRPC.Query != nil {
			logger.Info("request received")

			records, err := wakuPX.enrCache.getENRs(int(requestRPC.Query.NumPeers))
			if err != nil {
				logger.Error("obtaining enrs from cache", zap.Error(err))
				wakuPX.metrics.RecordError(pxFailure)
				return
			}

			responseRPC := &pb.PeerExchangeRPC{}
			responseRPC.Response = new(pb.PeerExchangeResponse)
			responseRPC.Response.PeerInfos = records

			writer := pbio.NewDelimitedWriter(s)
			err = writer.WriteMsg(responseRPC)
			if err != nil {
				logger.Error("writing response", zap.Error(err))
				wakuPX.metrics.RecordError(pxFailure)
				return
			}
		}
	}
}

// Stop unmounts the peer exchange protocol
func (wakuPX *WakuPeerExchange) Stop() {
	wakuPX.CommonService.Stop(func() {
		wakuPX.h.RemoveStreamHandler(PeerExchangeID_v20alpha1)
	})
}

func (wakuPX *WakuPeerExchange) iterate(ctx context.Context) error {
	iterator, err := wakuPX.disc.PeerIterator()
	if err != nil {
		return fmt.Errorf("obtaining iterator: %w", err)
	}
	// Closing iterator
	defer iterator.Close()

	for iterator.Next() {
		_, addresses, err := enr.Multiaddress(iterator.Node())
		if err != nil {
			wakuPX.log.Error("extracting multiaddrs from enr", zap.Error(err))
			continue
		}

		if len(addresses) == 0 {
			continue
		}

		wakuPX.log.Debug("Discovered px peers via discv5")
		wakuPX.enrCache.updateCache(iterator.Node())

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
	return nil
}

func (wakuPX *WakuPeerExchange) runPeerExchangeDiscv5Loop(ctx context.Context) {
	defer wakuPX.WaitGroup().Done()

	// Runs a discv5 loop adding new peers to the px peer cache
	if wakuPX.disc == nil {
		wakuPX.log.Warn("trying to run discovery v5 (for PX) while it's disabled")
		return
	}

	for {
		err := wakuPX.iterate(ctx)
		if err != nil {
			wakuPX.log.Debug("iterating peer exchange", zap.Error(err))
			time.Sleep(2 * time.Second)
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
