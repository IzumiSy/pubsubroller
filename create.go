package main

import (
	"context"
	"pubsubroller/client"
	config "pubsubroller/config"
	subscription "pubsubroller/subscription"
	topic "pubsubroller/topic"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func createSubscriptions(c client.PubsubClient, callbacks SubscriptionCallbacks, ctx context.Context, conf config.Configuration, opts Options) {
	egSubscriptions := errgroup.Group{}
	subscriptionSkippedCount := 0
	subscriptionCreatedCount := 0

	callbacks.Initialized()

	for _, sub := range subscription.FromConfig(conf, opts.Variables, c) {
		sub := sub

		egSubscriptions.Go(func() error {
			if !opts.IsDryRun {
				if err := sub.Create(c, ctx); err != nil {
					if errors.Cause(err) == subscription.SUBSCRIPTION_EXISTS_ERR {
						subscriptionSkippedCount += 1
						return nil
					} else {
						return err
					}
				}
			}

			subscriptionCreatedCount += 1
			callbacks.Each(sub)
			return nil
		})
	}

	if err := egSubscriptions.Wait(); err != nil {
		panic(err)
	}

	callbacks.Finalized(subscriptionCreatedCount, subscriptionSkippedCount)
}

func createTopics(c client.PubsubClient, callbacks TopicCallbacks, ctx context.Context, conf config.Configuration, opts Options) {
	egTopics := errgroup.Group{}
	topicSkippedCount := 0
	topicCreatedCount := 0

	callbacks.Initialized()

	for _, tp := range topic.FromConfig(conf, opts.Variables, c) {
		tp := tp

		egTopics.Go(func() error {
			if !opts.IsDryRun {
				if err := tp.Create(c, ctx); err != nil {
					if errors.Cause(err) == topic.TOPIC_EXISTS_ERR {
						topicSkippedCount += 1
						return nil
					} else {
						return err
					}
				}
			}

			topicCreatedCount += 1
			callbacks.Each(tp)
			return nil
		})
	}

	if err := egTopics.Wait(); err != nil {
		panic(err)
	}

	callbacks.Finalized(topicCreatedCount, topicSkippedCount)
}
