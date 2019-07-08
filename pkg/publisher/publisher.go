package publisher

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

var (
	clients    = make(map[string]*pubsub.Client)
	publishers = make(map[string]*Publisher)
)

// Publisher has topic information to publish the pubsub message.
type Publisher struct {
	*pubsub.Topic
}

// New returns a publisher.
func New(ctx context.Context, project, topic string) (*Publisher, error) {
	key := project + "/" + topic
	if p, cached := publishers[key]; cached {
		return p, nil
	}
	var err error
	c, cached := clients[project]
	if !cached {
		c, err = pubsub.NewClient(ctx, project)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to initialize pubsub client for %q project", project)
		}
		clients[project] = c
	}
	p := &Publisher{
		Topic: c.Topic(topic),
	}
	publishers[key] = p
	return p, nil
}

// Close waits to close all clients.
func Close() error {
	var result *multierror.Error
	for p, c := range clients {
		err := c.Close()
		result = multierror.Append(
			result,
			errors.Wrapf(err, "failed to close pubsub client for %q project", p),
		)
	}
	return result.ErrorOrNil()
}
