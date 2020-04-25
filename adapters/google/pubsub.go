package google

import (
	"cloud.google.com/go/pubsub"
	"pubsubroller/subscription"
	"pubsubroller/topic"
)

type PubsubClient struct {
	Client *pubsub.Client
}

func (c PubsubClient) CreateTopic(tp topic.Topic) error {
	return nil
}

func (c PubsubClient) DeleteTopic(tp topic.Topic) error {
	return nil
}

func (c PubsubClient) CreateSubscription(sub subscription.Subscription) error {
	return nil
}

func (c PubsubClient) DeleteSubscription(sub subscription.Subscription) error {
	return nil
}
