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

func deleteSubscriptions(client *pubsub.Client, ctx context.Context, conf config.Configuration, opts Options) {
	egSubscriptions := errgroup.Group{}
	subscriptionSkippedCount := 0
	subscriptionDeletedCount := 0

	fmt.Println("Start deleting subscriptions...")

	for _, sub := range subscription.FromConfig(conf, opts.Variables, client) {
		sub := sub

		egSubscriptions.Go(func() error {
			if !opts.IsDryRun {
				if err := sub.Delete(client, ctx); err != nil {
					if errors.Cause(err) == subscription.SUBSCRIPTION_NOT_FOUND_ERR {
						subscriptionSkippedCount += 1
						return nil
					} else {
						return err
					}
				}
			}

			subscriptionDeletedCount += 1
			return nil
		})
	}

	if err := egSubscriptions.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("Subscriptions deleted: %d, skipped: %d\n", subscriptionDeletedCount, subscriptionSkippedCount)
}

func deleteTopics(client *pubsub.Client, ctx context.Context, conf config.Configuration, opts Options) {
	egTopics := errgroup.Group{}
	topicSkippedCount := 0
	topicDeletedCount := 0

	fmt.Println("Start deleting topics...")

	for _, tp := range topic.FromConfig(conf, opts.Variables, client) {
		tp := tp

		egTopics.Go(func() error {
			if !opts.IsDryRun {
				if err := tp.Delete(client, ctx); err != nil {
					if errors.Cause(err) == topic.TOPIC_NOT_FOUND_ERR {
						topicSkippedCount += 1
						return nil
					} else {
						return err
					}
				}
			}

			topicDeletedCount += 1
			return nil
		})
	}

	if err := egTopics.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("Topics deleted: %d, skipped: %d\n", topicDeletedCount, topicSkippedCount)
}
