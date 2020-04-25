package main

import (
	config "pubsubroller/config"
	"pubsubroller/subscription"
	"pubsubroller/topic"
)

// mock structures for testing

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

var (
	mockConfig = config.Configuration{
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
)
