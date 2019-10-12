package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	config "pubsubroller/config"
	topic "pubsubroller/topic"
)

func createTopics(client *pubsub.Client, ctx context.Context, conf config.Configuration, opts Options) {
	egTopics := errgroup.Group{}
	topicSkippedCount := 0
	topicCreatedCount := 0

	fmt.Printf("\nStart creating topics...\n\n")
	for topicName, _ := range conf.Topics() {
		topicName := topicName

		egTopics.Go(func() error {
			err :=
				topic.
					New(topicName).
					Create(client, ctx)

			if err != nil {
				if errors.Cause(err) == topic.TOPIC_EXISTS_ERR {
					topicSkippedCount += 1
					return nil
				} else {
					return err
				}
			}

			topicCreatedCount += 1
			return nil
		})
	}

	if err := egTopics.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("\nTopics created: %d, skipped: %d\n", topicCreatedCount, topicSkippedCount)
}
