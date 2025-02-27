package rest

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/waku-org/go-waku/waku/v2/node"
	"github.com/waku-org/go-waku/waku/v2/utils"
)

func TestWakuRest(t *testing.T) {
	options := node.WithWakuStore()
	n, err := node.New(options)
	require.NoError(t, err)

	rpc := NewWakuRest(n, "127.0.0.1", 8080, false, 10, utils.Logger())
	require.NotNil(t, rpc.server)
	require.Equal(t, rpc.server.Addr, "127.0.0.1:8080")
}
