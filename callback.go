package main

type SubscriptionCallbacks interface {
	Initialized()
	Each(subscription namable)
	Finalized(done int, skipped int)
}

type TopicCallbacks interface {
	Initialized()
	Each(topic namable)
	Finalized(done int, skipped int)
}
