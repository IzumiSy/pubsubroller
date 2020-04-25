package subscription

import (
	"github.com/pkg/errors"
	"pubsubroller/config"
	"strings"
)

type Subscription struct {
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

	for _, topic := range conf.Topics() {
		topic := topic

		for _, sub := range topic.Subscriptions() {
			endpoint := sub.Endpoint
			for key, value := range variables {
				endpoint = strings.Replace(endpoint, "${"+key+"}", value, -1)
			}

			subscriptions =
				append(subscriptions, Subscription{sub.Name, endpoint, sub.Pull})
		}
	}

	return subscriptions
}
