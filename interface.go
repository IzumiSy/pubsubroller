package main

type topicLike interface {
	Delete() error
	Exists() (bool, error)
}

type subscriptionLike interface {
	Delete() error
	Exists() (bool, error)
}

type pubsubClient interface {
	CreateTopic(name string) error
	CreateSubscription(name string) error
}
