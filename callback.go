package main

import (
	subscription "pubsubroller/subscription"
	topic "pubsubroller/topic"
)

type SubscriptionCallbacks interface {
	Initialized()
	Each(subscription subscription.Subscription)
	Finalized(done int, skipped int)
}

type TopicCallbacks interface {
	Initialized()
	Each(topic topic.Topic)
	Finalized(done int, skipped int)
}
