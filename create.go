package main

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	config "pubsubroller/config"
	"pubsubroller/subscription"
	"pubsubroller/topic"
)

func createSubscriptions(client pubsubClient, callbacks SubscriptionCallbacks, ctx context.Context, conf config.Configuration, opts Options) {
	egSubscriptions := errgroup.Group{}
	subscriptionSkippedCount := 0
	subscriptionCreatedCount := 0

	callbacks.Initialized()

	for _, sub := range subscription.FromConfig(conf, opts.Variables) {
		sub := sub

		egSubscriptions.Go(func() error {
			if !opts.IsDryRun {
				if err := client.CreateSubscription(sub); err != nil {
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

func createTopics(client pubsubClient, callbacks TopicCallbacks, ctx context.Context, conf config.Configuration, opts Options) {
	egTopics := errgroup.Group{}
	topicSkippedCount := 0
	topicCreatedCount := 0

	callbacks.Initialized()

	for _, tp := range topic.FromConfig(conf, opts.Variables) {
		tp := tp

		egTopics.Go(func() error {
			if !opts.IsDryRun {
				if err := client.CreateTopic(tp); err != nil {
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
