package topic

import (
	"github.com/pkg/errors"
	"pubsubroller/config"
)

type Topic struct {
	Name string
}

var (
	INTERNAL_ERR        error = errors.New("Internal error")
	TOPIC_EXISTS_ERR    error = errors.New("Topic already exists")
	TOPIC_NOT_FOUND_ERR error = errors.New("Topic not found")
)

func FromConfig(conf config.Configuration, variables map[string]string) []Topic {
	var topics []Topic

	for topicName := range conf.Topics() {
		topicName := topicName
		topics = append(topics, Topic{topicName})
	}

	return topics
}
