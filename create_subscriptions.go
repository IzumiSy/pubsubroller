package main

import (
	"context"
	"fmt"
	config "pubsubroller/config"
	subscription "pubsubroller/subscription"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func createSubscriptions(client *pubsub.Client, ctx context.Context, conf config.Configuration, opts Options) {
	egSubscriptions := errgroup.Group{}
	subscriptionSkippedCount := 0
	subscriptionCreatedCount := 0

	fmt.Printf("\nStart creating subscriptions...\n\n")

	for _, sub := range subscription.FromConfig(conf, opts.Variables, client) {
		sub := sub

		egSubscriptions.Go(func() error {
			if err := sub.Create(client, ctx); err != nil {
				if errors.Cause(err) == subscription.SUBSCRIPTION_EXISTS_ERR {
					subscriptionSkippedCount += 1
					return nil
				} else {
					return err
				}
			}

			subscriptionCreatedCount += 1
			return nil
		})
	}

	if err := egSubscriptions.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("\nSubscriptions created: %d, skipped: %d\n", subscriptionCreatedCount, subscriptionSkippedCount)
}
