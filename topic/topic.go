package topic

import (
	"context"
	"pubsubroller/config"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

type Topic struct {
	name string
}

func New(name string) Topic {
	return Topic{name: name}
}

var (
	INTERNAL_ERR     error = errors.New("Internal error")
	TOPIC_EXISTS_ERR error = errors.New("Topic already exists")
)

func (topic Topic) Create(client *pubsub.Client, ctx context.Context) error {
	exists, err := client.Topic(topic.name).Exists(ctx)
	if err != nil {
		return errors.Wrap(err, INTERNAL_ERR.Error())
	}

	if exists {
		return TOPIC_EXISTS_ERR
	}

	_, err = client.CreateTopic(ctx, topic.name)
	if err != nil {
		return errors.Wrap(err, INTERNAL_ERR.Error())
	}

	return nil
}

func FromConfig(conf config.Configuration, variables map[string]string, client *pubsub.Client) []Topic {
	var topics []Topic

	for topicName, _ := range conf.Topics() {
		topicName := topicName
		topics = append(topics, New(topicName))
	}

	return topics
}
