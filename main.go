package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
)

func main() {
	projectIdPtr := flag.String("projectId", "", "target GCP project ID")
	flag.Parse()
	projectId := *projectIdPtr

	if projectId == "" {
		fmt.Println("Error: GCP project ID required with `-projectId` option.")
		return
	}

	fmt.Println("Target project ID:", projectId)

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		fmt.Println("Error on initializing pubsub client:", err.Error())
		return
	}

	topicId := "hogeTopic"

	exists, err := client.Topic(topicId).Exists(ctx)
	if err != nil {
		fmt.Println("Error on checking existence of pub/sub topic:", err.Error())
		return
	}

	if exists {
		fmt.Println("Skip creating pub/sub topic")
		return
	}

	_, err = client.CreateTopic(ctx, topicId)
	if err != nil {
		fmt.Println("Error on creating new pub/sub topic:", err.Error())
	}

	return
}
