package main

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	config "pubsubroller/config"
	"pubsubroller/subscription"
	"pubsubroller/topic"
)

func deleteSubscriptions(client pubsubClient, callbacks SubscriptionCallbacks, ctx context.Context, conf config.Configuration, opts appOptions) {
	subscriptions := subscription.FromConfig(conf, opts.Variables)
	counter := counter{Total: len(subscriptions)}
	egSubscriptions := errgroup.Group{}

	callbacks.Initialized()

	for _, sub := range subscriptions {
		sub := sub

		egSubscriptions.Go(func() error {
			if !opts.IsDryRun {
				if err := client.DeleteSubscription(sub); err != nil {
					if errors.Cause(err) == subscription.SUBSCRIPTION_NOT_FOUND_ERR {
						counter.Skipped()
						return nil
					} else {
						return err
					}
				}
			}

			counter.Done()
			callbacks.Each(sub)
			return nil
		})
	}

	if err := egSubscriptions.Wait(); err != nil {
		panic(err)
	}

	callbacks.Finalized(&counter)
}

func deleteTopics(client pubsubClient, callbacks TopicCallbacks, ctx context.Context, conf config.Configuration, opts appOptions) {
	topics := topic.FromConfig(conf, opts.Variables)
	counter := counter{Total: len(topics)}
	egTopics := errgroup.Group{}

	callbacks.Initialized()

	for _, tp := range topics {
		tp := tp

		egTopics.Go(func() error {
			if !opts.IsDryRun {
				if err := client.DeleteTopic(tp); err != nil {
					if errors.Cause(err) == topic.TOPIC_NOT_FOUND_ERR {
						counter.Skipped()
						return nil
					} else {
						return err
					}
				}
			}

			counter.Done()
			callbacks.Each(tp)
			return nil
		})
	}

	if err := egTopics.Wait(); err != nil {
		panic(err)
	}

	callbacks.Finalized(&counter)
}
