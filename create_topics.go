package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	config "pubsubroller/config"
)

func createTopics(client *pubsub.Client, ctx context.Context, conf config.Configuration, opts Options) {
	egTopics := errgroup.Group{}
	topicSkippedCount := 0
	topicCreatedCount := 0

	fmt.Printf("\nStart creating topics...\n\n")
	for topicName, _ := range conf.Topics() {
		topicName := topicName

		egTopics.Go(func() error {
			isCreated, err := createTopic(client, ctx, topicName, opts.IsDryRun)
			if isCreated {
				topicCreatedCount += 1
			} else {
				topicSkippedCount += 1
			}
			return err
		})
	}

	if err := egTopics.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("\nTopics created: %d, skipped: %d\n", topicCreatedCount, topicSkippedCount)
}

func createTopic(client *pubsub.Client, ctx context.Context, topicId string, isDryRun bool) (bool, error) {
	exists, err := client.Topic(topicId).Exists(ctx)
	if err != nil {
		return false, err
	}

	if exists {
		fmt.Println("Skip:", topicId)
		return false, nil
	}

	// dryrunであれば必ず成功扱いとするため作成系の実行はスキップする
	if !isDryRun {
		_, err = client.CreateTopic(ctx, topicId)
		if err != nil {
			return false, err
		}
	}

	fmt.Println("Created:", topicId)
	return true, nil
}
