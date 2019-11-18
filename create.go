package main

import (
	"context"
	"fmt"
	config "pubsubroller/config"
	subscription "pubsubroller/subscription"
	topic "pubsubroller/topic"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func createSubscriptions(client *pubsub.Client, ctx context.Context, conf config.Configuration, opts Options) {
	egSubscriptions := errgroup.Group{}
	subscriptionSkippedCount := 0
	subscriptionCreatedCount := 0

	fmt.Println("Start creating subscriptions...")

	for _, sub := range subscription.FromConfig(conf, opts.Variables, client) {
		sub := sub

		egSubscriptions.Go(func() error {
			if !opts.IsDryRun {
				if err := sub.Create(client, ctx); err != nil {
					if errors.Cause(err) == subscription.SUBSCRIPTION_EXISTS_ERR {
						subscriptionSkippedCount += 1
						return nil
					} else {
						return err
					}
				}
			}

			subscriptionCreatedCount += 1
			return nil
		})
	}

	if err := egSubscriptions.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("Subscriptions created: %d, skipped: %d\n", subscriptionCreatedCount, subscriptionSkippedCount)
}

func createTopics(client *pubsub.Client, ctx context.Context, conf config.Configuration, opts Options) {
	egTopics := errgroup.Group{}
	topicSkippedCount := 0
	topicCreatedCount := 0

	fmt.Println("Start creating topics...")

	for _, tp := range topic.FromConfig(conf, opts.Variables, client) {
		tp := tp

		egTopics.Go(func() error {
			if !opts.IsDryRun {
				if err := tp.Create(client, ctx); err != nil {
					if errors.Cause(err) == topic.TOPIC_EXISTS_ERR {
						topicSkippedCount += 1
						return nil
					} else {
						return err
					}
				}
			}

			topicCreatedCount += 1
			return nil
		})
	}

	if err := egTopics.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("Topics created: %d, skipped: %d\n", topicCreatedCount, topicSkippedCount)
}
