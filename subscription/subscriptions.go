package subscription

import (
	"pubsubroller/config"
	"strings"

	"cloud.google.com/go/pubsub"
)

func FromConfig(conf config.Configuration, variables map[string]string, client *pubsub.Client) []Subscription {
	var subscriptions []Subscription

	for topicName, topic := range conf.Topics() {
		topicName := topicName
		topic := topic

		for _, sub := range topic.Subscriptions() {
			endpoint := sub.Endpoint
			for key, value := range variables {
				endpoint = strings.Replace(endpoint, "${"+key+"}", value, -1)
			}

			subscriptions =
				append(
					subscriptions,
					New(sub.Name, endpoint, sub.Pull, client.Topic(topicName)),
				)
		}
	}

	return subscriptions
}
