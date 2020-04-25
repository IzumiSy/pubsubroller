package main

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	client "pubsubroller/adapters/google"
	subscription "pubsubroller/adapters/google/subscription"
	topic "pubsubroller/adapters/google/topic"
	config "pubsubroller/config"
)

func deleteSubscriptions(c client.PubsubClient, callbacks SubscriptionCallbacks, ctx context.Context, conf config.Configuration, opts Options) {
	egSubscriptions := errgroup.Group{}
	subscriptionSkippedCount := 0
	subscriptionDeletedCount := 0

	callbacks.Initialized()

	for _, sub := range subscription.FromConfig(conf, opts.Variables, c) {
		sub := sub

		egSubscriptions.Go(func() error {
			if !opts.IsDryRun {
				if err := sub.Delete(c, ctx); err != nil {
					if errors.Cause(err) == subscription.SUBSCRIPTION_NOT_FOUND_ERR {
						subscriptionSkippedCount += 1
						return nil
					} else {
						return err
					}
				}
			}

			subscriptionDeletedCount += 1
			callbacks.Each(sub)
			return nil
		})
	}

	if err := egSubscriptions.Wait(); err != nil {
		panic(err)
	}

	callbacks.Finalized(subscriptionDeletedCount, subscriptionSkippedCount)
}

func deleteTopics(c client.PubsubClient, callbacks TopicCallbacks, ctx context.Context, conf config.Configuration, opts Options) {
	egTopics := errgroup.Group{}
	topicSkippedCount := 0
	topicDeletedCount := 0

	callbacks.Initialized()

	for _, tp := range topic.FromConfig(conf, opts.Variables, c) {
		tp := tp

		egTopics.Go(func() error {
			if !opts.IsDryRun {
				if err := tp.Delete(c, ctx); err != nil {
					if errors.Cause(err) == topic.TOPIC_NOT_FOUND_ERR {
						topicSkippedCount += 1
						return nil
					} else {
						return err
					}
				}
			}

			topicDeletedCount += 1
			callbacks.Each(tp)
			return nil
		})
	}

	if err := egTopics.Wait(); err != nil {
		panic(err)
	}

	callbacks.Finalized(topicDeletedCount, topicSkippedCount)
}
