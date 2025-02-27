package lightpush

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/waku-org/go-waku/tests"
	"github.com/waku-org/go-waku/waku/v2/utils"
)

func TestLightPushOption(t *testing.T) {
	port, err := tests.FindFreePort(t, "", 5)
	require.NoError(t, err)

	host, err := tests.MakeHost(context.Background(), port, rand.Reader)
	require.NoError(t, err)

	options := []Option{
		WithPeer("QmWLxGxG65CZ7vRj5oNXCJvbY9WkF9d9FxuJg8cg8Y7q3"),
		WithAutomaticPeerSelection(),
		WithFastestPeerSelection(context.Background()),
		WithRequestID([]byte("requestID")),
		WithAutomaticRequestID(),
	}

	params := new(lightPushParameters)
	params.host = host
	params.log = utils.Logger()

	for _, opt := range options {
		opt(params)
	}

	require.Equal(t, host, params.host)
	require.NotNil(t, params.selectedPeer)
	require.NotNil(t, params.requestID)
}
