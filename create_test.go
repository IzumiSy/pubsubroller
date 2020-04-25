package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	config "pubsubroller/config"
	"testing"
)

type fakeCallbacks struct {
	IsInitialized bool
	IsFinazlied   bool
	Calls         int
}

func (f *fakeCallbacks) Initialized() {
	f.IsInitialized = true
}

func (f *fakeCallbacks) Each(subscription namable) {
	f.Calls++
}

func (f *fakeCallbacks) Finalized(done int, skipped int) {
	f.IsFinazlied = true
}

type fakeClient struct{}

func (_ fakeClient) Subscription(id string) *pubsub.Subscription {
	return nil
}

func (_ fakeClient) CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	return nil, errors.New("mocked")
}

func (_ fakeClient) Topic(id string) *pubsub.Topic {
	return nil
}

func (_ fakeClient) CreateTopic(ctx context.Context, id string) (*pubsub.Topic, error) {
	return nil, errors.New("mocked")
}

func TestCreateSubscription(t *testing.T) {
	ctx := context.Background()
	conf := config.Configuration{
		Internal_Topics_: map[string]config.Topic{
			"topic1": config.Topic{
				Internal_Subscriptions_: []config.Subscription{
					config.Subscription{Name: "subscription1"},
					config.Subscription{Name: "subscription2"},
					config.Subscription{Name: "subscription3"},
				},
			},
		},
	}
	opts := Options{}
	cb := fakeCallbacks{}

	createSubscriptions(fakeClient{}, &cb, ctx, conf, opts)

	if !cb.IsInitialized {
		t.Error("It must be initialized")
	}

	if !cb.IsFinazlied {
		t.Error("It must be finalized")
	}
}
