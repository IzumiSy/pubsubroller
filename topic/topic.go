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

func (topic Topic) Name() string {
	return topic.name
}

var (
	INTERNAL_ERR        error = errors.New("Internal error")
	TOPIC_EXISTS_ERR    error = errors.New("Topic already exists")
	TOPIC_NOT_FOUND_ERR error = errors.New("Topic not found")
)

type pubsubClient interface {
	Topic(id string) *pubsub.Topic
	Subscription(id string) *pubsub.Subscription
	CreateTopic(ctx context.Context, id string) (*pubsub.Topic, error)
}

func (topic Topic) Create(client pubsubClient, ctx context.Context) error {
	exists, err := client.Topic(topic.name).Exists(ctx)
	if err != nil {
		return errors.Wrap(err, INTERNAL_ERR.Error())
	}

	if exists {
		return TOPIC_EXISTS_ERR
	}

	_, err = client.CreateTopic(ctx, topic.name)
	return errors.Wrap(err, INTERNAL_ERR.Error())
}

func (topic Topic) Delete(client pubsubClient, ctx context.Context) error {
	tp := client.Topic(topic.name)

	exists, err := tp.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, INTERNAL_ERR.Error())
	}

	if !exists {
		return TOPIC_NOT_FOUND_ERR
	}

	return errors.Wrap(tp.Delete(ctx), INTERNAL_ERR.Error())
}

func FromConfig(conf config.Configuration, variables map[string]string, client pubsubClient) []Topic {
	var topics []Topic

	for topicName, _ := range conf.Topics() {
		topicName := topicName
		topics = append(topics, New(topicName))
	}

	return topics
}
