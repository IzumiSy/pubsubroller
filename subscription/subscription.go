package subscription

import (
	"github.com/pkg/errors"
	"pubsubroller/config"
	"pubsubroller/topic"
	"strings"
)

type Subscription struct {
	Topic    topic.Topic
	Name     string
	Endpoint string
	Pull     bool
}

var (
	INTERNAL_ERR               error = errors.New("Internal error")
	SUBSCRIPTION_EXISTS_ERR    error = errors.New("Subscription already exists")
	SUBSCRIPTION_NOT_FOUND_ERR error = errors.New("Subscription not found")
	NO_ENDPOINT_SPECIFIED_ERR  error = errors.New("No endpoint specified")
)

func FromConfig(conf config.Configuration, variables map[string]string) []Subscription {
	var subscriptions []Subscription

	for topicName, tp := range conf.Topics() {
		tp := tp

		for _, sub := range tp.Subscriptions() {
			endpoint := sub.Endpoint
			name := sub.Name
			for key, value := range variables {
				endpoint = strings.Replace(endpoint, "${"+key+"}", value, -1)
				name = strings.Replace(name, "${"+key+"}", value, -1)
			}

			subscriptions =
				append(subscriptions, Subscription{
					Topic:    topic.Topic{topicName},
					Name:     name,
					Endpoint: endpoint,
					Pull:     sub.Pull,
				})
		}
	}

	return subscriptions
}
