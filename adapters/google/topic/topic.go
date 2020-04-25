package topic

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/pkg/errors"
	"pubsubroller/topic"
)

type Topic struct {
	name string
}

func New(topic topic.Topic) Topic {
	return Topic{
		name: topic.Name,
	}
}

func (tp Topic) Create(c *pubsub.Client, ctx context.Context) error {
	exists, err := c.Topic(tp.name).Exists(ctx)
	if err != nil {
		return errors.Wrap(err, topic.INTERNAL_ERR.Error())
	}

	if exists {
		return topic.TOPIC_EXISTS_ERR
	}

	_, err = c.CreateTopic(ctx, tp.name)
	return errors.Wrap(err, topic.INTERNAL_ERR.Error())
}

func (tp Topic) Delete(c *pubsub.Client, ctx context.Context) error {
	deletingTopic := c.Topic(tp.name)

	exists, err := deletingTopic.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, topic.INTERNAL_ERR.Error())
	}

	if !exists {
		return topic.TOPIC_NOT_FOUND_ERR
	}

	return errors.Wrap(deletingTopic.Delete(ctx), topic.INTERNAL_ERR.Error())
}
