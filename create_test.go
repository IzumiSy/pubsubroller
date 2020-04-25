package main

import (
	"context"
	config "pubsubroller/config"
	"pubsubroller/subscription"
	"pubsubroller/topic"
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

func (f *fakeCallbacks) Each(_ subscription.Subscription) {
	f.Calls++
}

func (f *fakeCallbacks) Finalized(done int, skipped int) {
	f.IsFinazlied = true
}

type fakeClient struct{}

func (_ fakeClient) CreateSubscription(_ subscription.Subscription) error {
	return nil
}

func (_ fakeClient) DeleteSubscription(_ subscription.Subscription) error {
	return nil
}

func (_ fakeClient) CreateTopic(_ topic.Topic) error {
	return nil
}

func (_ fakeClient) DeleteTopic(_ topic.Topic) error {
	return nil
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
