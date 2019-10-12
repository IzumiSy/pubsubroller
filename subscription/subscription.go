package subscription

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/pkg/errors"
)

type Subscription struct {
	topic 	 *pubsub.Topic
	name     string `yaml:"name"`
	endpoint string `yaml:"endpoint,omitempty"`
	pull     bool   `yaml:"pull,omitempty"`
}

func New(name, endpoint string, pull bool, topic *pubsub.Topic) Subscription {
	return Subscription{
		topic: topic,
		name: name,
		endpoint: endpoint,
		pull: pull,
	}
}

var (
	INTERNAL_ERR error = errors.New("Internal error")
	SUBSCRIPTION_EXISTS_ERR error = errors.New("Subscription already exists")
	NO_ENDPOINT_SPECIFIED_ERR error = errors.New("No endpoint specified")
)

func (subscription Subscription) Create(client *pubsub.Client, ctx context.Context) error {
	s := client.Subscription(subscription.name)
	exists, err := s.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, INTERNAL_ERR.Error())
	}

	if exists {
		return SUBSCRIPTION_EXISTS_ERR 
	}

	if subscription.endpoint == "" {
		return errors.WithMessage(NO_ENDPOINT_SPECIFIED_ERR, subscription.name)
	}

	var pushConfig pubsub.PushConfig
	if subscription.pull {
		pushConfig = pubsub.PushConfig{}
	} else {
		pushConfig = pubsub.PushConfig{Endpoint: subscription.endpoint}
	}

	_, err = client.CreateSubscription(
			ctx, 
			subscription.name, 
			pubsub.SubscriptionConfig{
				Topic: subscription.topic,
				PushConfig: pushConfig,
			},
		)
	if err != nil {
		return errors.Wrap(err, INTERNAL_ERR.Error())
	}

	return nil
}