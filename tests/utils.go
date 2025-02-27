package tests

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/ethereum/go-ethereum/log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/host/peerstore/pstoremem"
	"github.com/multiformats/go-multiaddr"
	"github.com/waku-org/go-waku/waku/v2/peerstore"
	"github.com/waku-org/go-waku/waku/v2/protocol/pb"
)

// GetHostAddress returns the first listen address used by a host
func GetHostAddress(ha host.Host) multiaddr.Multiaddr {
	return ha.Addrs()[0]
}

// FindFreePort returns an available port number
func FindFreePort(t *testing.T, host string, maxAttempts int) (int, error) {
	t.Helper()

	if host == "" {
		host = "localhost"
	}

	for i := 0; i < maxAttempts; i++ {
		addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Logf("unable to resolve tcp addr: %v", err)
			continue
		}
		l, err := net.ListenTCP("tcp", addr)
		if err != nil {
			l.Close()
			t.Logf("unable to listen on addr %q: %v", addr, err)
			continue
		}

		port := l.Addr().(*net.TCPAddr).Port
		l.Close()
		return port, nil

	}

	return 0, fmt.Errorf("no free port found")
}

// FindFreePort returns an available port number
func FindFreeUDPPort(t *testing.T, host string, maxAttempts int) (int, error) {
	t.Helper()

	if host == "" {
		host = "localhost"
	}

	for i := 0; i < maxAttempts; i++ {
		addr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Logf("unable to resolve tcp addr: %v", err)
			continue
		}
		l, err := net.ListenUDP("udp", addr)
		if err != nil {
			l.Close()
			t.Logf("unable to listen on addr %q: %v", addr, err)
			continue
		}

		port := l.LocalAddr().(*net.UDPAddr).Port
		l.Close()
		return port, nil

	}

	return 0, fmt.Errorf("no free port found")
}

// MakeHost creates a Libp2p host with a random key on a specific port
func MakeHost(ctx context.Context, port int, randomness io.Reader) (host.Host, error) {
	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randomness)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port))
	if err != nil {
		return nil, err
	}

	ps, err := pstoremem.NewPeerstore()
	if err != nil {
		return nil, err
	}

	psWrapper := peerstore.NewWakuPeerstore(ps)
	if err != nil {
		return nil, err
	}

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	return libp2p.New(
		libp2p.Peerstore(psWrapper),
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
}

// CreateWakuMessage creates a WakuMessage protobuffer with default values and a custom contenttopic and timestamp
func CreateWakuMessage(contentTopic string, timestamp int64, optionalPayload ...string) *pb.WakuMessage {
	var payload []byte
	if len(optionalPayload) > 0 {
		payload = []byte(optionalPayload[0])
	} else {
		payload = []byte{1, 2, 3}
	}
	return &pb.WakuMessage{Payload: payload, ContentTopic: contentTopic, Version: 0, Timestamp: timestamp}
}

// RandomHex returns a random hex string of n bytes
func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
