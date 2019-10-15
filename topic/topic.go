package topic

import (
	"context"
	"pubsubroller/client"
	"pubsubroller/config"

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

func (topic Topic) Create(c client.TopicClient, ctx context.Context) error {
	exists, err := c.Topic(topic.name).Exists(ctx)
	if err != nil {
		return errors.Wrap(err, INTERNAL_ERR.Error())
	}

	if exists {
		return TOPIC_EXISTS_ERR
	}

	_, err = c.CreateTopic(ctx, topic.name)
	return errors.Wrap(err, INTERNAL_ERR.Error())
}

func (topic Topic) Delete(c client.TopicClient, ctx context.Context) error {
	tp := c.Topic(topic.name)

	exists, err := tp.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, INTERNAL_ERR.Error())
	}

	if !exists {
		return TOPIC_NOT_FOUND_ERR
	}

	return errors.Wrap(tp.Delete(ctx), INTERNAL_ERR.Error())
}

func FromConfig(conf config.Configuration, variables map[string]string, c client.TopicClient) []Topic {
	var topics []Topic

	for topicName, _ := range conf.Topics() {
		topicName := topicName
		topics = append(topics, New(topicName))
	}

	return topics
}
