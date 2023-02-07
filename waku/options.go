package waku

import (
	"crypto/ecdsa"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/multiformats/go-multiaddr"
	"github.com/urfave/cli/v2"
)

// DiscV5Options are settings to enable a modified version of Ethereum’s Node
// Discovery Protocol v5 as a means for ambient node discovery.
type DiscV5Options struct {
	Enable     bool
	Nodes      cli.StringSlice
	Port       uint
	AutoUpdate bool
}

// RelayOptions are settings to enable the relay protocol which is a pubsub
// approach to peer-to-peer messaging with a strong focus on privacy,
// censorship-resistance, security and scalability.
type RelayOptions struct {
	Enable                 bool
	Topics                 cli.StringSlice
	PeerExchange           bool
	MinRelayPeersToPublish int
}

type RLNRelayOptions struct {
	Enable                    bool
	CredentialsPath           string
	CredentialsPassword       string
	MembershipIndex           int
	PubsubTopic               string
	ContentTopic              string
	Dynamic                   bool
	ETHPrivateKey             *ecdsa.PrivateKey
	ETHClientAddress          string
	MembershipContractAddress common.Address
}

// FilterOptions are settings used to enable filter protocol. This is a protocol
// that enables subscribing to messages that a peer receives. This is a more
// lightweight version of WakuRelay specifically designed for bandwidth
// restricted devices.
type FilterOptions struct {
	Enable          bool
	UseV2           bool
	DisableFullNode bool
	Nodes           []multiaddr.Multiaddr
	Timeout         time.Duration
}

// LightpushOptions are settings used to enable the lightpush protocol. This is
// a lightweight protocol used to avoid having to run the relay protocol which
// is more resource intensive. With this protocol a message is pushed to a peer
// that supports both the lightpush protocol and relay protocol. That peer will
// broadcast the message and return a confirmation that the message was
// broadcasted
type LightpushOptions struct {
	Enable bool
	Nodes  []multiaddr.Multiaddr
}

// StoreOptions are settings used for enabling the store protocol, used to
// retrieve message history from other nodes as well as acting as a store
// node and provide message history to nodes that ask for it.
type StoreOptions struct {
	Enable               bool
	DatabaseURL          string
	RetentionTime        time.Duration
	RetentionMaxMessages int
	ResumeNodes          []multiaddr.Multiaddr
	Nodes                []multiaddr.Multiaddr
}

// SwapOptions are settings used for configuring the swap protocol
type SwapOptions struct {
	Enable              bool
	Mode                int
	PaymentThreshold    int
	DisconnectThreshold int
}

// DNSDiscoveryOptions are settings used for enabling DNS-based discovery
// protocol that stores merkle trees in DNS records which contain connection
// information for nodes. It's very useful for bootstrapping a p2p network.
type DNSDiscoveryOptions struct {
	Enable     bool
	URLs       cli.StringSlice
	Nameserver string
}

// MetricsOptions are settings used to start a prometheus server for obtaining
// useful node metrics to monitor the health of behavior of the go-waku node.
type MetricsOptions struct {
	Enable  bool
	Address string
	Port    int
}

// RPCServerOptions are settings used to start a json rpc server
type RPCServerOptions struct {
	Enable             bool
	Port               int
	Address            string
	Admin              bool
	Private            bool
	RelayCacheCapacity int
}

// RESTServerOptions are settings used to start a rest http server
type RESTServerOptions struct {
	Enable             bool
	Port               int
	Address            string
	Admin              bool
	Private            bool
	RelayCacheCapacity int
}

// WSOptions are settings used for enabling websockets and secure websockets
// support
type WSOptions struct {
	Enable   bool
	WSPort   int
	WSSPort  int
	Address  string
	Secure   bool
	KeyPath  string
	CertPath string
}

// PeerExchangeOptions are settings used with the peer exchange protocol
type PeerExchangeOptions struct {
	Enable bool
	Node   *multiaddr.Multiaddr
}

// Options contains all the available features and settings that can be
// configured via flags when executing go-waku as a service.
type Options struct {
	Port               int
	Address            string
	Dns4DomainName     string
	NodeKey            *ecdsa.PrivateKey
	KeyFile            string
	KeyPasswd          string
	GenerateKey        bool
	Overwrite          bool
	StaticNodes        []multiaddr.Multiaddr
	KeepAlive          time.Duration
	AdvertiseAddresses []multiaddr.Multiaddr
	ShowAddresses      bool
	LogLevel           string
	LogEncoding        string
	LogOutput          string
	NAT                string
	PersistPeers       bool
	UserAgent          string
	PProf              bool

	PeerExchange PeerExchangeOptions
	Websocket    WSOptions
	Relay        RelayOptions
	Store        StoreOptions
	Swap         SwapOptions
	Filter       FilterOptions
	LightPush    LightpushOptions
	RLNRelay     RLNRelayOptions
	DiscV5       DiscV5Options
	DNSDiscovery DNSDiscoveryOptions
	Metrics      MetricsOptions
	RPCServer    RPCServerOptions
	RESTServer   RESTServerOptions
}
