package client

import (
	"context"

	"cloud.google.com/go/pubsub"
)

type TopicClient interface {
	Topic(id string) *pubsub.Topic
	CreateTopic(ctx context.Context, id string) (*pubsub.Topic, error)
}

type SubscriptionClient interface {
	Subscription(id string) *pubsub.Subscription
	CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error)
}

type PubsubClient interface {
	TopicClient
	SubscriptionClient
}
