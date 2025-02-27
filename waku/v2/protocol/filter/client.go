package filter

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"net/http"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	libp2pProtocol "github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-msgio/pbio"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/waku-org/go-waku/logging"
	"github.com/waku-org/go-waku/waku/v2/peermanager"
	"github.com/waku-org/go-waku/waku/v2/protocol"
	"github.com/waku-org/go-waku/waku/v2/protocol/filter/pb"
	wpb "github.com/waku-org/go-waku/waku/v2/protocol/pb"
	"github.com/waku-org/go-waku/waku/v2/protocol/relay"
	"github.com/waku-org/go-waku/waku/v2/timesource"
	"go.uber.org/zap"
)

// FilterPushID_v20beta1 is the current Waku Filter protocol identifier used to allow
// filter service nodes to push messages matching registered subscriptions to this client.
const FilterPushID_v20beta1 = libp2pProtocol.ID("/vac/waku/filter-push/2.0.0-beta1")

var (
	ErrNoPeersAvailable = errors.New("no suitable remote peers")
)

type WakuFilterLightNode struct {
	*protocol.CommonService
	h             host.Host
	broadcaster   relay.Broadcaster //TODO: Move the broadcast functionality outside of relay client to a higher SDK layer.s
	timesource    timesource.Timesource
	metrics       Metrics
	log           *zap.Logger
	subscriptions *SubscriptionsMap
	pm            *peermanager.PeerManager
}

type ContentFilter struct {
	Topic         string
	ContentTopics []string
}

type WakuFilterPushResult struct {
	Err    error
	PeerID peer.ID
}

// NewWakuFilterLightnode returns a new instance of Waku Filter struct setup according to the chosen parameter and options
// Note that broadcaster is optional.
// Takes an optional peermanager if WakuFilterLightnode is being created along with WakuNode.
// If using libp2p host, then pass peermanager as nil
func NewWakuFilterLightNode(broadcaster relay.Broadcaster, pm *peermanager.PeerManager,
	timesource timesource.Timesource, reg prometheus.Registerer, log *zap.Logger) *WakuFilterLightNode {
	wf := new(WakuFilterLightNode)
	wf.log = log.Named("filterv2-lightnode")
	wf.broadcaster = broadcaster
	wf.timesource = timesource
	wf.pm = pm
	wf.CommonService = protocol.NewCommonService()
	wf.metrics = newMetrics(reg)

	return wf
}

// Sets the host to be able to mount or consume a protocol
func (wf *WakuFilterLightNode) SetHost(h host.Host) {
	wf.h = h
}

func (wf *WakuFilterLightNode) Start(ctx context.Context) error {
	return wf.CommonService.Start(ctx, wf.start)

}

func (wf *WakuFilterLightNode) start() error {
	wf.subscriptions = NewSubscriptionMap(wf.log)
	wf.h.SetStreamHandlerMatch(FilterPushID_v20beta1, protocol.PrefixTextMatch(string(FilterPushID_v20beta1)), wf.onRequest(wf.Context()))

	wf.log.Info("filter-push protocol started")
	return nil
}

// Stop unmounts the filter protocol
func (wf *WakuFilterLightNode) Stop() {
	wf.CommonService.Stop(func() {
		wf.h.RemoveStreamHandler(FilterPushID_v20beta1)
		res, err := wf.unsubscribeAll(wf.Context())
		if err != nil {
			wf.log.Warn("unsubscribing from full nodes", zap.Error(err))
		}

		for r := range res {
			if r.Err != nil {
				wf.log.Warn("unsubscribing from full nodes", zap.Error(r.Err), logging.HostID("peerID", r.PeerID))
			}

		}
		//
		wf.subscriptions.Clear()
	})
}

func (wf *WakuFilterLightNode) onRequest(ctx context.Context) func(s network.Stream) {
	return func(s network.Stream) {
		defer s.Close()
		logger := wf.log.With(logging.HostID("peer", s.Conn().RemotePeer()))

		if !wf.subscriptions.IsSubscribedTo(s.Conn().RemotePeer()) {
			logger.Warn("received message push from unknown peer", logging.HostID("peerID", s.Conn().RemotePeer()))
			wf.metrics.RecordError(unknownPeerMessagePush)
			return
		}

		reader := pbio.NewDelimitedReader(s, math.MaxInt32)

		messagePush := &pb.MessagePushV2{}
		err := reader.ReadMsg(messagePush)
		if err != nil {
			logger.Error("reading message push", zap.Error(err))
			wf.metrics.RecordError(decodeRPCFailure)
			return
		}
		//For now returning failure, this will get addressed with autosharding changes for filter.
		if messagePush.PubsubTopic == nil {
			logger.Error("empty pubsub topic")
			wf.metrics.RecordError(decodeRPCFailure)
			return
		}
		if !wf.subscriptions.Has(s.Conn().RemotePeer(), *messagePush.PubsubTopic, messagePush.WakuMessage.ContentTopic) {
			logger.Warn("received messagepush with invalid subscription parameters",
				logging.HostID("peerID", s.Conn().RemotePeer()), zap.String("topic", *messagePush.PubsubTopic),
				zap.String("contentTopic", messagePush.WakuMessage.ContentTopic))
			wf.metrics.RecordError(invalidSubscriptionMessage)
			return
		}

		wf.metrics.RecordMessage()

		wf.notify(s.Conn().RemotePeer(), *messagePush.PubsubTopic, messagePush.WakuMessage)

		logger.Info("received message push")
	}
}

func (wf *WakuFilterLightNode) notify(remotePeerID peer.ID, pubsubTopic string, msg *wpb.WakuMessage) {
	envelope := protocol.NewEnvelope(msg, wf.timesource.Now().UnixNano(), pubsubTopic)

	if wf.broadcaster != nil {
		// Broadcasting message so it's stored
		wf.broadcaster.Submit(envelope)
	}
	// Notify filter subscribers
	wf.subscriptions.Notify(remotePeerID, envelope)
}

func (wf *WakuFilterLightNode) request(ctx context.Context, params *FilterSubscribeParameters,
	reqType pb.FilterSubscribeRequest_FilterSubscribeType, contentFilter ContentFilter) error {
	conn, err := wf.h.NewStream(ctx, params.selectedPeer, FilterSubscribeID_v20beta1)
	if err != nil {
		wf.metrics.RecordError(dialFailure)
		return err
	}
	defer conn.Close()

	writer := pbio.NewDelimitedWriter(conn)
	reader := pbio.NewDelimitedReader(conn, math.MaxInt32)

	request := &pb.FilterSubscribeRequest{
		RequestId:           hex.EncodeToString(params.requestID),
		FilterSubscribeType: reqType,
		PubsubTopic:         &contentFilter.Topic,
		ContentTopics:       contentFilter.ContentTopics,
	}

	wf.log.Debug("sending FilterSubscribeRequest", zap.Stringer("request", request))
	err = writer.WriteMsg(request)
	if err != nil {
		wf.metrics.RecordError(writeRequestFailure)
		wf.log.Error("sending FilterSubscribeRequest", zap.Error(err))
		return err
	}

	filterSubscribeResponse := &pb.FilterSubscribeResponse{}
	err = reader.ReadMsg(filterSubscribeResponse)
	if err != nil {
		wf.log.Error("receiving FilterSubscribeResponse", zap.Error(err))
		wf.metrics.RecordError(decodeRPCFailure)
		return err
	}

	if filterSubscribeResponse.RequestId != request.RequestId {
		wf.log.Error("requestID mismatch", zap.String("expected", request.RequestId), zap.String("received", filterSubscribeResponse.RequestId))
		wf.metrics.RecordError(requestIDMismatch)
		err := NewFilterError(300, "request_id_mismatch")
		return &err
	}

	if filterSubscribeResponse.StatusCode != http.StatusOK {
		wf.metrics.RecordError(errorResponse)
		err := NewFilterError(int(filterSubscribeResponse.StatusCode), filterSubscribeResponse.StatusDesc)
		return &err
	}

	return nil
}

// Subscribe setups a subscription to receive messages that match a specific content filter
func (wf *WakuFilterLightNode) Subscribe(ctx context.Context, contentFilter ContentFilter, opts ...FilterSubscribeOption) (*SubscriptionDetails, error) {
	wf.RLock()
	defer wf.RUnlock()
	if err := wf.ErrOnNotRunning(); err != nil {
		return nil, err
	}

	if contentFilter.Topic == "" {
		return nil, errors.New("topic is required")
	}

	if len(contentFilter.ContentTopics) == 0 {
		return nil, errors.New("at least one content topic is required")
	}

	if len(contentFilter.ContentTopics) > MaxContentTopicsPerRequest {
		return nil, fmt.Errorf("exceeds maximum content topics: %d", MaxContentTopicsPerRequest)
	}

	params := new(FilterSubscribeParameters)
	params.log = wf.log
	params.host = wf.h
	params.pm = wf.pm

	optList := DefaultSubscriptionOptions()
	optList = append(optList, opts...)
	for _, opt := range optList {
		opt(params)
	}

	if params.selectedPeer == "" {
		wf.metrics.RecordError(peerNotFoundFailure)
		return nil, ErrNoPeersAvailable
	}

	err := wf.request(ctx, params, pb.FilterSubscribeRequest_SUBSCRIBE, contentFilter)
	if err != nil {
		return nil, err
	}
	return wf.subscriptions.NewSubscription(params.selectedPeer, contentFilter.Topic, contentFilter.ContentTopics), nil
}

// FilterSubscription is used to obtain an object from which you could receive messages received via filter protocol
func (wf *WakuFilterLightNode) FilterSubscription(peerID peer.ID, contentFilter ContentFilter) (*SubscriptionDetails, error) {
	wf.RLock()
	defer wf.RUnlock()
	if err := wf.ErrOnNotRunning(); err != nil {
		return nil, err
	}

	if !wf.subscriptions.Has(peerID, contentFilter.Topic, contentFilter.ContentTopics...) {
		return nil, errors.New("subscription does not exist")
	}

	return wf.subscriptions.NewSubscription(peerID, contentFilter.Topic, contentFilter.ContentTopics), nil
}

func (wf *WakuFilterLightNode) getUnsubscribeParameters(opts ...FilterUnsubscribeOption) (*FilterUnsubscribeParameters, error) {
	params := new(FilterUnsubscribeParameters)
	params.log = wf.log
	opts = append(DefaultUnsubscribeOptions(), opts...)
	for _, opt := range opts {
		opt(params)
	}

	return params, nil
}

func (wf *WakuFilterLightNode) Ping(ctx context.Context, peerID peer.ID) error {
	wf.RLock()
	defer wf.RUnlock()
	if err := wf.ErrOnNotRunning(); err != nil {
		return err
	}

	return wf.request(
		ctx,
		&FilterSubscribeParameters{selectedPeer: peerID},
		pb.FilterSubscribeRequest_SUBSCRIBER_PING,
		ContentFilter{})
}

func (wf *WakuFilterLightNode) IsSubscriptionAlive(ctx context.Context, subscription *SubscriptionDetails) error {
	wf.RLock()
	defer wf.RUnlock()
	if err := wf.ErrOnNotRunning(); err != nil {
		return err
	}

	return wf.Ping(ctx, subscription.PeerID)
}

func (wf *WakuFilterLightNode) Subscriptions() []*SubscriptionDetails {
	wf.RLock()
	defer wf.RUnlock()
	if err := wf.ErrOnNotRunning(); err != nil {
		return nil
	}

	wf.subscriptions.RLock()
	defer wf.subscriptions.RUnlock()

	var output []*SubscriptionDetails

	for _, peerSubscription := range wf.subscriptions.items {
		for _, subscriptionPerTopic := range peerSubscription.subscriptionsPerTopic {
			for _, subscriptionDetail := range subscriptionPerTopic {
				output = append(output, subscriptionDetail)
			}
		}
	}

	return output
}

func (wf *WakuFilterLightNode) cleanupSubscriptions(peerID peer.ID, contentFilter ContentFilter) {
	wf.subscriptions.Lock()
	defer wf.subscriptions.Unlock()

	peerSubscription, ok := wf.subscriptions.items[peerID]
	if !ok {
		return
	}

	subscriptionDetailList, ok := peerSubscription.subscriptionsPerTopic[contentFilter.Topic]
	if !ok {
		return
	}

	for subscriptionDetailID, subscriptionDetail := range subscriptionDetailList {
		subscriptionDetail.Remove(contentFilter.ContentTopics...)
		if len(subscriptionDetail.ContentTopics) == 0 {
			delete(subscriptionDetailList, subscriptionDetailID)
		} else {
			subscriptionDetailList[subscriptionDetailID] = subscriptionDetail
		}
	}

	if len(subscriptionDetailList) == 0 {
		delete(wf.subscriptions.items[peerID].subscriptionsPerTopic, contentFilter.Topic)
	} else {
		wf.subscriptions.items[peerID].subscriptionsPerTopic[contentFilter.Topic] = subscriptionDetailList
	}

}

// Unsubscribe is used to stop receiving messages from a peer that match a content filter
func (wf *WakuFilterLightNode) Unsubscribe(ctx context.Context, contentFilter ContentFilter, opts ...FilterUnsubscribeOption) (<-chan WakuFilterPushResult, error) {
	wf.RLock()
	defer wf.RUnlock()
	if err := wf.ErrOnNotRunning(); err != nil {
		return nil, err
	}

	if contentFilter.Topic == "" {
		return nil, errors.New("topic is required")
	}

	if len(contentFilter.ContentTopics) == 0 {
		return nil, errors.New("at least one content topic is required")
	}

	if len(contentFilter.ContentTopics) > MaxContentTopicsPerRequest {
		return nil, fmt.Errorf("exceeds maximum content topics: %d", MaxContentTopicsPerRequest)
	}

	params, err := wf.getUnsubscribeParameters(opts...)
	if err != nil {
		return nil, err
	}

	resultChan := make(chan WakuFilterPushResult, len(wf.subscriptions.items))
	for peerID := range wf.subscriptions.items {
		if params.selectedPeer != "" && peerID != params.selectedPeer {
			continue
		}

		subscriptions, ok := wf.subscriptions.items[peerID]
		if !ok || subscriptions == nil {
			continue
		}

		wf.cleanupSubscriptions(peerID, contentFilter)
		if len(subscriptions.subscriptionsPerTopic) == 0 {
			delete(wf.subscriptions.items, peerID)
		}

		if params.wg != nil {
			params.wg.Add(1)
		}

		go func(peerID peer.ID) {
			defer func() {
				if params.wg != nil {
					params.wg.Done()
				}
			}()

			err := wf.request(
				ctx,
				&FilterSubscribeParameters{selectedPeer: peerID, requestID: params.requestID},
				pb.FilterSubscribeRequest_UNSUBSCRIBE,
				contentFilter)
			if err != nil {
				ferr, ok := err.(*FilterError)
				if ok && ferr.Code == http.StatusNotFound {
					wf.log.Warn("peer does not have a subscription", logging.HostID("peerID", peerID), zap.Error(err))
				} else {
					wf.log.Error("could not unsubscribe from peer", logging.HostID("peerID", peerID), zap.Error(err))
					return
				}
			}

			if params.wg != nil {
				resultChan <- WakuFilterPushResult{
					Err:    err,
					PeerID: peerID,
				}
			}
		}(peerID)
	}

	if params.wg != nil {
		params.wg.Wait()
	}

	close(resultChan)

	return resultChan, nil
}

// Unsubscribe is used to stop receiving messages from a peer that match a content filter
func (wf *WakuFilterLightNode) UnsubscribeWithSubscription(ctx context.Context, sub *SubscriptionDetails, opts ...FilterUnsubscribeOption) (<-chan WakuFilterPushResult, error) {
	wf.RLock()
	defer wf.RUnlock()
	if err := wf.ErrOnNotRunning(); err != nil {
		return nil, err
	}

	var contentTopics []string
	for k := range sub.ContentTopics {
		contentTopics = append(contentTopics, k)
	}

	opts = append(opts, Peer(sub.PeerID))

	return wf.Unsubscribe(ctx, ContentFilter{Topic: sub.PubsubTopic, ContentTopics: contentTopics}, opts...)
}

func (wf *WakuFilterLightNode) unsubscribeAll(ctx context.Context, opts ...FilterUnsubscribeOption) (<-chan WakuFilterPushResult, error) {
	params, err := wf.getUnsubscribeParameters(opts...)
	if err != nil {
		return nil, err
	}

	wf.subscriptions.Lock()
	defer wf.subscriptions.Unlock()

	resultChan := make(chan WakuFilterPushResult, len(wf.subscriptions.items))

	for peerID := range wf.subscriptions.items {
		if params.selectedPeer != "" && peerID != params.selectedPeer {
			continue
		}

		delete(wf.subscriptions.items, peerID)

		if params.wg != nil {
			params.wg.Add(1)
		}

		go func(peerID peer.ID) {
			defer func() {
				if params.wg != nil {
					params.wg.Done()
				}
			}()

			err := wf.request(
				ctx,
				&FilterSubscribeParameters{selectedPeer: peerID, requestID: params.requestID},
				pb.FilterSubscribeRequest_UNSUBSCRIBE_ALL,
				ContentFilter{})
			if err != nil {
				wf.log.Error("could not unsubscribe from peer", logging.HostID("peerID", peerID), zap.Error(err))
			}
			if params.wg != nil {
				resultChan <- WakuFilterPushResult{
					Err:    err,
					PeerID: peerID,
				}
			}
		}(peerID)
	}

	if params.wg != nil {
		params.wg.Wait()
	}

	close(resultChan)

	return resultChan, nil
}

// UnsubscribeAll is used to stop receiving messages from peer(s). It does not close subscriptions
func (wf *WakuFilterLightNode) UnsubscribeAll(ctx context.Context, opts ...FilterUnsubscribeOption) (<-chan WakuFilterPushResult, error) {
	wf.RLock()
	defer wf.RUnlock()
	if err := wf.ErrOnNotRunning(); err != nil {
		return nil, err
	}

	return wf.unsubscribeAll(ctx, opts...)
}
