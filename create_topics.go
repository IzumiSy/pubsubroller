package main

import (
	"context"
	"fmt"
	config "pubsubroller/config"
	topic "pubsubroller/topic"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func createTopics(client *pubsub.Client, ctx context.Context, conf config.Configuration, opts Options) {
	egTopics := errgroup.Group{}
	topicSkippedCount := 0
	topicCreatedCount := 0

	fmt.Printf("\nStart creating topics...\n\n")

	for _, tp := range topic.FromConfig(conf, opts.Variables, client) {
		tp := tp

		egTopics.Go(func() error {
			if !opts.IsDryRun {
				if err := tp.Create(client, ctx); err != nil {
					if errors.Cause(err) == topic.TOPIC_EXISTS_ERR {
						topicSkippedCount += 1
						return nil
					} else {
						return err
					}
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
