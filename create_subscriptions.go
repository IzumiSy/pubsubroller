package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"strings"
)

func createSubscriptions(client *pubsub.Client, ctx context.Context, config Configuration, opts Options) {
	egSubscriptions := errgroup.Group{}
	subscriptionSkippedCount := 0
	subscriptionCreatedCount := 0

	fmt.Printf("\nStart creating subscriptions...\n\n")
	for topicName, topic := range config.Topics {
		topicName := topicName
		topic := topic

		for _, subscription := range topic.Subsciptions {
			subscription := subscription
			endpoint := subscription.Endpoint
			for key, value := range opts.Variables {
				endpoint = strings.Replace(endpoint, "${"+key+"}", value, -1)
			}

			subscription.Endpoint = endpoint
			egSubscriptions.Go(func() error {
				isCreated, err := createSubscription(client, ctx, subscription, topicName, opts.IsDryRun)
				if isCreated {
					subscriptionCreatedCount += 1
				} else {
					subscriptionSkippedCount += 1
				}
				return err
			})
		}
	}

	if err := egSubscriptions.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("\nSubscriptions Created: %d, Skipped: %d\n", subscriptionCreatedCount, subscriptionSkippedCount)
}

func createSubscription(client *pubsub.Client, ctx context.Context, subscription Subscription, topicName string, isDryRun bool) (bool, error) {
	name := subscription.Name

	endpoint := subscription.Endpoint

	s := client.Subscription(name)
	exists, err := s.Exists(ctx)
	if err != nil {
		return false, err
	}

	if exists {
		fmt.Println("Skip:", name)
		return false, nil
	}

	// dryrunであれば必ず成功扱いとするため作成系の実行はスキップする
	if !isDryRun {
		pushConfig := pubsub.PushConfig{
			Endpoint: endpoint,
		}

		// 空のpubsub.PushConfigを指定してpullなsubscriptionにする
		if subscription.Pull {
			pushConfig = pubsub.PushConfig{}
		} else {
			if subscription.Endpoint == "" {
				return false, fmt.Errorf("Failed because no endpoint specified to subscription: %s", name)
			}
		}

		topic := client.Topic(topicName)
		_, err = client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
			Topic:      topic,
			PushConfig: pushConfig,
		})
		if err != nil {
			return false, err
		}
	}

	fmt.Println("Created:", name)
	return true, nil
}
