package main

import (
	"pubsubroller/subscription"
	"pubsubroller/topic"
)

type pubsubClient interface {
	CreateTopic(topic topic.Topic) error
	DeleteTopic(topic topic.Topic) error
	CreateSubscription(subscription subscription.Subscription) error
	DeleteSubscription(subscription subscription.Subscription) error
}
