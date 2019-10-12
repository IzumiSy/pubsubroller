package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	config "pubsubroller/config"
)

func deleteSubscriptions(client *pubsub.Client, ctx context.Context, conf config.Configuration, opts Options) {
	egSubscriptions := errgroup.Group{}
	subscriptionSkippedCount := 0
	subscriptionDeletedCount := 0

	fmt.Printf("\nStart deleting subscriptions...\n\n")

	for _, topic := range conf.Topics() {
		for _, subscription := range topic.Subscriptions() {
			subscription := subscription

			egSubscriptions.Go(func() error {
				isDeleted, err := deleteSubscription(client, ctx, subscription.Name, opts.IsDryRun)
				if isDeleted {
					subscriptionDeletedCount += 1
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

	fmt.Printf("\nSubscriptions deleted: %d, skipped: %d\n", subscriptionDeletedCount, subscriptionSkippedCount)
}

func deleteSubscription(client *pubsub.Client, ctx context.Context, subscriptionId string, isDryRun bool) (bool, error) {
	s := client.Subscription(subscriptionId)

	exists, err := s.Exists(ctx)
	if err != nil {
		return false, err
	}

	if !exists {
		fmt.Println("Skip:", subscriptionId)
		return false, nil
	}

	// dryrunであれば必ず成功扱いとするため削除系の実行はスキップする
	if !isDryRun {
		if err := s.Delete(ctx); err != nil {
			return false, err
		}
	}

	fmt.Println("Deleted:", subscriptionId)
	return true, nil
}
