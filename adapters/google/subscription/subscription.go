package subscription

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/pkg/errors"
	"pubsubroller/subscription"
	"pubsubroller/topic"
)

type Subscription struct {
	topic    topic.Topic
	name     string
	endpoint string
	pull     bool
}

func New(subscription subscription.Subscription) Subscription {
	return Subscription{
		topic:    subscription.Topic,
		name:     subscription.Name,
		endpoint: subscription.Endpoint,
		pull:     subscription.Pull,
	}
}

func (sub Subscription) Create(c *pubsub.Client, ctx context.Context) error {
	s := c.Subscription(sub.name)
	exists, err := s.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, subscription.INTERNAL_ERR.Error())
	}

	if exists {
		return subscription.SUBSCRIPTION_EXISTS_ERR
	}

	var pushConfig pubsub.PushConfig
	if sub.pull {
		pushConfig = pubsub.PushConfig{}
	} else {
		if sub.endpoint == "" {
			return errors.WithMessage(subscription.NO_ENDPOINT_SPECIFIED_ERR, sub.name)
		}
		pushConfig = pubsub.PushConfig{Endpoint: sub.endpoint}
	}

	_, err = c.CreateSubscription(
		ctx,
		sub.name,
		pubsub.SubscriptionConfig{
			Topic:      c.Topic(sub.topic.Name),
			PushConfig: pushConfig,
		},
	)

	return errors.Wrap(err, subscription.INTERNAL_ERR.Error())
}

func (sub Subscription) Delete(c *pubsub.Client, ctx context.Context) error {
	s := c.Subscription(sub.name)
	exists, err := s.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, subscription.INTERNAL_ERR.Error())
	}

	if !exists {
		return subscription.SUBSCRIPTION_NOT_FOUND_ERR
	}

	return errors.Wrap(s.Delete(ctx), subscription.INTERNAL_ERR.Error())
}
