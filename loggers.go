package main

import (
	"fmt"
	subscription "pubsubroller/subscription"
	topic "pubsubroller/topic"
)

type CreateTopicsLogger struct{}

func (_ CreateTopicsLogger) Initialized() {
	fmt.Println("Start creating topics...")
}

func (_ CreateTopicsLogger) Each(topic topic.Topic) {
	fmt.Printf("Topic created: %s\n", topic.Name())
}

func (_ CreateTopicsLogger) Finalized(done int, skipped int) {
	fmt.Printf("Topics created: %d, skipped: %d\n", done, skipped)
}

type CreateSubscriptionsLogger struct{}

func (_ CreateSubscriptionsLogger) Initialized() {
	fmt.Println("Start creating subscriptions...")
}

func (_ CreateSubscriptionsLogger) Each(subscription subscription.Subscription) {
	fmt.Printf("Subscription creatd: %s\n", subscription.Name())
}

func (_ CreateSubscriptionsLogger) Finalized(done int, skipped int) {
	fmt.Printf("Subscriptions created: %d, skipped: %d\n", done, skipped)
}

type DeleteTopicsLogger struct{}

func (_ DeleteTopicsLogger) Initialized() {
	fmt.Println("Start deleting topics...")
}

func (_ DeleteTopicsLogger) Each(topic topic.Topic) {
	fmt.Printf("Topic deleted: %s\n", topic.Name())
}

func (_ DeleteTopicsLogger) Finalized(done int, skipped int) {
	fmt.Printf("Topics deleted: %d, skipped: %d\n", done, skipped)
}

type DeleteSubscriptionLogger struct{}

func (_ DeleteSubscriptionLogger) Initialized() {
	fmt.Println("Start deleting subscriptions...")
}

func (_ DeleteSubscriptionLogger) Each(subscription subscription.Subscription) {
	fmt.Printf("Subscription deleted: %s\n", subscription.Name())
}

func (_ DeleteSubscriptionLogger) Finalized(done int, skipped int) {
	fmt.Printf("Subscriptions deleted: %d, skipped: %d\n", done, skipped)
}
