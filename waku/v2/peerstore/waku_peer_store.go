package peerstore

import (
	"sync"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
)

// Origin is used to determine how the peer is identified,
// either it is statically added or discovered via one of the discovery protocols
type Origin int64

const (
	Unknown Origin = iota
	Discv5
	Static
	PeerExchange
	DNSDiscovery
	Rendezvous
)

const peerOrigin = "origin"
const peerENR = "enr"
const peerDirection = "direction"
const peerPubSubTopics = "pubSubTopics"

// ConnectionFailures contains connection failure information towards all peers
type ConnectionFailures struct {
	sync.RWMutex
	failures map[peer.ID]int
}

// WakuPeerstoreImpl is a implementation of WakuPeerStore
type WakuPeerstoreImpl struct {
	peerStore    peerstore.Peerstore
	connFailures ConnectionFailures
}

// WakuPeerstore is an interface for implementing WakuPeerStore
type WakuPeerstore interface {
	SetOrigin(p peer.ID, origin Origin) error
	Origin(p peer.ID) (Origin, error)
	PeersByOrigin(origin Origin) peer.IDSlice
	SetENR(p peer.ID, enr *enode.Node) error
	ENR(p peer.ID) (*enode.Node, error)
	AddConnFailure(p peer.AddrInfo)
	ResetConnFailures(p peer.AddrInfo)
	ConnFailures(p peer.AddrInfo) int

	SetDirection(p peer.ID, direction network.Direction) error
	Direction(p peer.ID) (network.Direction, error)

	AddPubSubTopic(p peer.ID, topic string) error
	PubSubTopics(p peer.ID) ([]string, error)
	SetPubSubTopics(p peer.ID, topics []string) error
	PeersByPubSubTopic(pubSubTopic string) peer.IDSlice
}

// NewWakuPeerstore creates a new WakuPeerStore object
func NewWakuPeerstore(p peerstore.Peerstore) peerstore.Peerstore {
	return &WakuPeerstoreImpl{
		peerStore: p,
		connFailures: ConnectionFailures{
			failures: make(map[peer.ID]int),
		},
	}
}

// SetOrigin sets origin for a specific peer.
func (ps *WakuPeerstoreImpl) SetOrigin(p peer.ID, origin Origin) error {
	return ps.peerStore.Put(p, peerOrigin, origin)
}

// Origin fetches the origin for a specific peer.
func (ps *WakuPeerstoreImpl) Origin(p peer.ID) (Origin, error) {
	result, err := ps.peerStore.Get(p, peerOrigin)
	if err != nil {
		return Unknown, err
	}

	return result.(Origin), nil
}

// PeersByOrigin returns the list of peers for a specific origin
func (ps *WakuPeerstoreImpl) PeersByOrigin(expectedOrigin Origin) peer.IDSlice {
	var result peer.IDSlice
	for _, p := range ps.Peers() {
		actualOrigin, err := ps.Origin(p)
		if err == nil && actualOrigin == expectedOrigin {
			result = append(result, p)
		}
	}
	return result
}

// SetENR sets the ENR record a peer
func (ps *WakuPeerstoreImpl) SetENR(p peer.ID, enr *enode.Node) error {
	return ps.peerStore.Put(p, peerENR, enr)
}

// ENR fetches the ENR record for a peer
func (ps *WakuPeerstoreImpl) ENR(p peer.ID) (*enode.Node, error) {
	result, err := ps.peerStore.Get(p, peerENR)
	if err != nil {
		return nil, err
	}
	return result.(*enode.Node), nil
}

// AddConnFailure increments connectionFailures for a peer
func (ps *WakuPeerstoreImpl) AddConnFailure(p peer.AddrInfo) {
	ps.connFailures.Lock()
	defer ps.connFailures.Unlock()
	ps.connFailures.failures[p.ID]++
}

// ResetConnFailures resets connectionFailures for a peer to 0
func (ps *WakuPeerstoreImpl) ResetConnFailures(p peer.AddrInfo) {
	ps.connFailures.Lock()
	defer ps.connFailures.Unlock()
	ps.connFailures.failures[p.ID] = 0
}

// ConnFailures fetches connectionFailures for a peer
func (ps *WakuPeerstoreImpl) ConnFailures(p peer.AddrInfo) int {
	ps.connFailures.RLock()
	defer ps.connFailures.RUnlock()
	return ps.connFailures.failures[p.ID]
}

// SetDirection sets connection direction for a specific peer.
func (ps *WakuPeerstoreImpl) SetDirection(p peer.ID, direction network.Direction) error {
	return ps.peerStore.Put(p, peerDirection, direction)
}

// Direction fetches the connection direction (Inbound or outBound) for a specific peer
func (ps *WakuPeerstoreImpl) Direction(p peer.ID) (network.Direction, error) {
	result, err := ps.peerStore.Get(p, peerDirection)
	if err != nil {
		return network.DirUnknown, err
	}

	return result.(network.Direction), nil
}

// AddPubSubTopic adds a new pubSubTopic for a peer
func (ps *WakuPeerstoreImpl) AddPubSubTopic(p peer.ID, topic string) error {
	existingTopics, err := ps.PubSubTopics(p)
	if err != nil {
		return err
	}
	existingTopics = append(existingTopics, topic)
	return ps.peerStore.Put(p, peerPubSubTopics, existingTopics)
}

// SetPubSubTopics sets pubSubTopics for a peer, it also overrides existing ones that were set previously..
func (ps *WakuPeerstoreImpl) SetPubSubTopics(p peer.ID, topics []string) error {
	return ps.peerStore.Put(p, peerPubSubTopics, topics)
}

// PubSubTopics fetches list of pubSubTopics for a peer
func (ps *WakuPeerstoreImpl) PubSubTopics(p peer.ID) ([]string, error) {
	result, err := ps.peerStore.Get(p, peerPubSubTopics)
	if err != nil {
		return nil, err
	}
	return result.([]string), nil
}

// PeersByPubSubTopic Returns list of peers by pubSubTopic
func (ps *WakuPeerstoreImpl) PeersByPubSubTopic(pubSubTopic string) peer.IDSlice {
	var result peer.IDSlice
	for _, p := range ps.Peers() {
		topics, err := ps.PubSubTopics(p)
		if err == nil {
			for _, topic := range topics {
				if topic == pubSubTopic {
					result = append(result, p)
				}
			}
		} //Note: skipping a peer in case of an error as there would be others available.
	}
	return result
}
