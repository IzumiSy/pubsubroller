package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
)

func deleteTopics(client *pubsub.Client, ctx context.Context, config Configuration, opts Options) {
	egTopics := errgroup.Group{}
	topicSkippedCount := 0
	topicDeletedCount := 0

	fmt.Printf("\nStart deleting topics...\n\n")
	for topicName, _ := range config.Topics {
		topicName := topicName

		egTopics.Go(func() error {
			isDeleted, err := deleteTopic(client, ctx, topicName, opts.IsDryRun)
			if isDeleted {
				topicDeletedCount += 1
			} else {
				topicSkippedCount += 1
			}
			return err
		})
	}

	if err := egTopics.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("\nTopics deleted: %d, skipped: %d\n", topicDeletedCount, topicSkippedCount)
}

func deleteTopic(client *pubsub.Client, ctx context.Context, topicId string, isDryRun bool) (bool, error) {
	topic := client.Topic(topicId)

	exists, err := topic.Exists(ctx)
	if err != nil {
		return false, err
	}

	if !exists {
		fmt.Println("Skip:", topicId)
		return false, nil
	}

	// dryrunであれば必ず成功扱いとするため削除系の実行はスキップする
	if !isDryRun {
		if err = topic.Delete(ctx); err != nil {
			return false, err
		}
	}

	fmt.Println("Deleted:", topicId)
	return true, nil
}
