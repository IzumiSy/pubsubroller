package topic

import (
	"pubsubroller/config"

	"cloud.google.com/go/pubsub"
)

func FromConfig(conf config.Configuration, variables map[string]string, client *pubsub.Client) []Topic {
	var topics []Topic

	for topicName, _ := range conf.Topics() {
		topicName := topicName
		topics = append(topics, New(topicName))
	}

	return topics
}
