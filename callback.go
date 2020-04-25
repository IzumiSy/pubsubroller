package main

import (
	"pubsubroller/subscription"
	"pubsubroller/topic"
)

type countable interface {
	Result() (int, int)
}

type SubscriptionCallbacks interface {
	Initialized()
	Each(subscription subscription.Subscription)
	Finalized(counter countable)
}

type TopicCallbacks interface {
	Initialized()
	Each(topic topic.Topic)
	Finalized(counter countable)
}
