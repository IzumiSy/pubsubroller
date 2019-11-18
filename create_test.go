package main

import (
	"context"
	"errors"
	config "pubsubroller/config"
	subscription "pubsubroller/subscription"

	"testing"

	"cloud.google.com/go/pubsub"
)

type FakeCallbacks struct {
	IsInitialized bool
	IsFinazlied   bool
	Calls         int
}

func (f *FakeCallbacks) Initialized() {
	f.IsInitialized = true
}

func (f *FakeCallbacks) Each(subscription subscription.Subscription) {
	f.Calls++
}

func (f *FakeCallbacks) Finalized(done int, skipped int) {
	f.IsFinazlied = true
}

type FakeClient struct{}

func (_ FakeClient) Subscription(id string) *pubsub.Subscription {
	return nil
}

func (_ FakeClient) CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	return nil, errors.New("mocked")
}

func (_ FakeClient) Topic(id string) *pubsub.Topic {
	return nil
}

func (_ FakeClient) CreateTopic(ctx context.Context, id string) (*pubsub.Topic, error) {
	return nil, errors.New("mocked")
}

func TestCreateSubscription(t *testing.T) {
	ctx := context.Background()
	conf := config.Configuration{
		Internal_Topics_: map[string]config.Topic{
			"topic1": config.Topic{
				Internal_Subscriptions_: []config.Subscription{
					config.Subscription{
						Name: "subscription1",
					},
				},
			},
			"topic2": config.Topic{
				Internal_Subscriptions_: []config.Subscription{},
			},
		},
	}
	opts := Options{}
	cb := FakeCallbacks{}

	createSubscriptions(FakeClient{}, &cb, ctx, conf, opts)

	if !cb.IsInitialized {
		t.Error("It must be initialized")
	}

	println(cb.Calls)

	if !cb.IsFinazlied {
		t.Error("It must be finalized")
	}
}
