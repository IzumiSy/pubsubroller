package google

import (
	"cloud.google.com/go/pubsub"
	"context"
	gcpSubscription "pubsubroller/adapters/google/subscription"
	gcpTopic "pubsubroller/adapters/google/topic"
	"pubsubroller/subscription"
	"pubsubroller/topic"
)

type PubsubClient struct {
	Client *pubsub.Client
	Ctx    context.Context
}

func (c PubsubClient) CreateTopic(tp topic.Topic) error {
	return gcpTopic.New(tp).Create(c.Client, c.Ctx)
}

func (c PubsubClient) DeleteTopic(tp topic.Topic) error {
	return gcpTopic.New(tp).Delete(c.Client, c.Ctx)
}

func (c PubsubClient) CreateSubscription(sub subscription.Subscription) error {
	return gcpSubscription.New(sub).Create(c.Client, c.Ctx)
}

func (c PubsubClient) DeleteSubscription(sub subscription.Subscription) error {
	return gcpSubscription.New(sub).Delete(c.Client, c.Ctx)
}
