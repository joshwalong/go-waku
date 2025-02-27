package filter

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/waku-org/go-waku/tests"
	"github.com/waku-org/go-waku/waku/v2/protocol"
	"github.com/waku-org/go-waku/waku/v2/utils"
)

func TestSubscriptionMapAppend(t *testing.T) {
	fmap := NewSubscriptionMap(utils.Logger())
	peerID := createPeerID(t)
	contentTopics := []string{"ct1", "ct2"}

	sub := fmap.NewSubscription(peerID, TOPIC, contentTopics)
	_, found := sub.ContentTopics["ct1"]
	require.True(t, found)
	_, found = sub.ContentTopics["ct2"]
	require.True(t, found)
	require.False(t, sub.Closed)
	require.Equal(t, sub.PeerID, peerID)
	require.Equal(t, sub.PubsubTopic, TOPIC)

	sub.Add("ct3")
	_, found = sub.ContentTopics["ct3"]
	require.True(t, found)

	sub.Remove("ct3")
	_, found = sub.ContentTopics["ct3"]
	require.False(t, found)

	err := sub.Close()
	require.NoError(t, err)
	require.True(t, sub.Closed)

	err = sub.Close()
	require.NoError(t, err)
}

func TestSubscriptionClear(t *testing.T) {
	fmap := NewSubscriptionMap(utils.Logger())
	contentTopics := []string{"ct1", "ct2"}

	var subscriptions = []*SubscriptionDetails{
		fmap.NewSubscription(createPeerID(t), TOPIC+"1", contentTopics),
		fmap.NewSubscription(createPeerID(t), TOPIC+"2", contentTopics),
		fmap.NewSubscription(createPeerID(t), TOPIC+"3", contentTopics),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(len(subscriptions))
	for _, s := range subscriptions {
		go func(s *SubscriptionDetails) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				t.Fail()
				return
			case <-s.C:
				return
			}
		}(s)
	}

	fmap.Clear()

	wg.Wait()

	require.True(t, subscriptions[0].Closed)
	require.True(t, subscriptions[1].Closed)
	require.True(t, subscriptions[2].Closed)
}

func TestSubscriptionsNotify(t *testing.T) {
	fmap := NewSubscriptionMap(utils.Logger())
	p1 := createPeerID(t)
	p2 := createPeerID(t)
	var subscriptions = []*SubscriptionDetails{
		fmap.NewSubscription(p1, TOPIC+"1", []string{"ct1", "ct2"}),
		fmap.NewSubscription(p2, TOPIC+"1", []string{"ct1"}),
		fmap.NewSubscription(p1, TOPIC+"2", []string{"ct1", "ct2"}),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	successChan := make(chan struct{}, 10)
	wg := sync.WaitGroup{}

	successOnReceive := func(ctx context.Context, i int) {
		defer wg.Done()

		if subscriptions[i].Closed {
			successChan <- struct{}{}
			return
		}

		select {
		case <-ctx.Done():
			panic("should have failed1")
		case c := <-subscriptions[i].C:
			if c == nil {
				panic("should have failed2")
			}
			successChan <- struct{}{}
			return
		}
	}

	failOnReceive := func(ctx context.Context, i int) {
		defer wg.Done()

		if subscriptions[i].Closed {
			successChan <- struct{}{}
			return
		}

		select {
		case <-ctx.Done():
			successChan <- struct{}{}
			return
		case c := <-subscriptions[i].C:
			if c != nil {
				panic("should have failed")
			}
			successChan <- struct{}{}
			return
		}
	}

	wg.Add(3)
	go successOnReceive(ctx, 0)
	go successOnReceive(ctx, 1)
	go failOnReceive(ctx, 2)
	time.Sleep(200 * time.Millisecond)

	envTopic1Ct1 := protocol.NewEnvelope(tests.CreateWakuMessage("ct1", 0), 0, TOPIC+"1")
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmap.Notify(p1, envTopic1Ct1)
		fmap.Notify(p2, envTopic1Ct1)
	}()

	<-successChan
	<-successChan
	cancel()
	wg.Wait()
	<-successChan

	//////////////////////////////////////

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	wg.Add(3)
	go successOnReceive(ctx, 0)
	go failOnReceive(ctx, 1)
	go failOnReceive(ctx, 2)
	time.Sleep(200 * time.Millisecond)

	envTopic1Ct2 := protocol.NewEnvelope(tests.CreateWakuMessage("ct2", 0), 0, TOPIC+"1")
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmap.Notify(p1, envTopic1Ct2)
		fmap.Notify(p2, envTopic1Ct2)
	}()

	<-successChan
	cancel()
	wg.Wait()
	<-successChan
	<-successChan

	//////////////////////////////////////

	// Testing after closing the subscription

	subscriptions[0].Close()
	time.Sleep(200 * time.Millisecond)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

	wg.Add(3)
	go failOnReceive(ctx, 0)
	go successOnReceive(ctx, 1)
	go failOnReceive(ctx, 2)
	time.Sleep(200 * time.Millisecond)

	envTopic1Ct1_2 := protocol.NewEnvelope(tests.CreateWakuMessage("ct1", 1), 1, TOPIC+"1")

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmap.Notify(p1, envTopic1Ct1_2)
		fmap.Notify(p2, envTopic1Ct1_2)
	}()

	<-successChan // One of these successes is for closing the subscription
	<-successChan
	cancel()
	wg.Wait()
	<-successChan
}
