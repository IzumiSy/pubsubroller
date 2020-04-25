package main

import (
	"fmt"
	"pubsubroller/subscription"
	"pubsubroller/topic"
)

// topic

type createTopicsLogger struct{}

func (_ createTopicsLogger) Initialized() {
	fmt.Println("Start creating topics...")
}

func (_ createTopicsLogger) Each(topic topic.Topic) {
	fmt.Printf("Topic created: %s\n", topic.Name)
}

func (_ createTopicsLogger) Finalized(done int, skipped int) {
	fmt.Printf("Topics created: %d, skipped: %d\n", done, skipped)
}

type deleteTopicsLogger struct{}

func (_ deleteTopicsLogger) Initialized() {
	fmt.Println("Start deleting topics...")
}

func (_ deleteTopicsLogger) Each(topic topic.Topic) {
	fmt.Printf("Topic deleted: %s\n", topic.Name)
}

func (_ deleteTopicsLogger) Finalized(done int, skipped int) {
	fmt.Printf("Topics deleted: %d, skipped: %d\n", done, skipped)
}

// subscription

type createSubscriptionsLogger struct{}

func (_ createSubscriptionsLogger) Initialized() {
	fmt.Println("Start creating subscriptions...")
}

func (_ createSubscriptionsLogger) Each(subscription subscription.Subscription) {
	fmt.Printf("Subscription creatd: %s\n", subscription.Name)
}

func (_ createSubscriptionsLogger) Finalized(done int, skipped int) {
	fmt.Printf("Subscriptions created: %d, skipped: %d\n", done, skipped)
}

type deleteSubscriptionLogger struct{}

func (_ deleteSubscriptionLogger) Initialized() {
	fmt.Println("Start deleting subscriptions...")
}

func (_ deleteSubscriptionLogger) Each(subscription subscription.Subscription) {
	fmt.Printf("Subscription deleted: %s\n", subscription.Name)
}

func (_ deleteSubscriptionLogger) Finalized(done int, skipped int) {
	fmt.Printf("Subscriptions deleted: %d, skipped: %d\n", done, skipped)
}
