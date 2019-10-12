package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/pkg/errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"strings"
	subscription "pubsubroller/subscription"
	config "pubsubroller/config"
)

func createSubscriptions(client *pubsub.Client, ctx context.Context, conf config.Configuration, opts Options) {
	egSubscriptions := errgroup.Group{}
	subscriptionSkippedCount := 0
	subscriptionCreatedCount := 0

	fmt.Printf("\nStart creating subscriptions...\n\n")
	for topicName, topic := range conf.Topics() {
		topicName := topicName
		topic := topic

		for _, sub := range topic.Subscriptions() {
			egSubscriptions.Go(func() error {
				endpoint := sub.Endpoint
				for key, value := range opts.Variables {
					endpoint = strings.Replace(endpoint, "${"+key+"}", value, -1)
				}

				err :=
					subscription.
						New(sub.Name, endpoint, sub.Pull, client.Topic(topicName)).
						Create(client, ctx)

				if err != nil {
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
	}

	if err := egSubscriptions.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("\nSubscriptions created: %d, skipped: %d\n", subscriptionCreatedCount, subscriptionSkippedCount)
}
