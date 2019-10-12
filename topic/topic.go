package topic

import (
	"github.com/pkg/errors"
	"cloud.google.com/go/pubsub"
	"context"
)

type Topic struct {
	name string
}

func New(name string) Topic {
	return Topic{name: name}
}

var (
	INTERNAL_ERR error = errors.New("Internal error")
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