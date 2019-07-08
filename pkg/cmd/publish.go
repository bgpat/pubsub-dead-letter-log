package cmd

import (
	"context"
	"fmt"
	"sort"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/wantedly/pubsub-dead-letter-log"
	"github.com/wantedly/pubsub-dead-letter-log/pkg/publisher"
)

var publishOption = struct {
	project string
	topic   string
}{}

func newPublishCommand() *cobra.Command {
	c := cobra.Command{
		Use:          "publish",
		Short:        "publish messages to Google Cloud Pub/Sub",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			ctx := context.Background()
			var logs []psdll.DeadLetterLog
			for _, arg := range args {
				l, err := psdll.ReadFromURL(ctx, arg)
				if err != nil {
					return errors.Wrapf(err, "invalid URL: %q", arg)
				}
				logs = append(logs, l...)
			}
			sort.Slice(logs, func(i, j int) bool {
				return logs[i].Timestamp.Before(logs[j].Timestamp)
			})

			for _, log := range logs {
				project := publishOption.project
				if project == "" {
					project = log.Project
				}

				topic := publishOption.topic
				if topic == "" {
					topic = log.Topic
				}

				p, err := publisher.New(ctx, project, topic)
				if err != nil {
					return errors.WithStack(err)
				}
				result := p.Publish(ctx, &pubsub.Message{
					Data:       log.Message.Data,
					Attributes: log.Message.Attributes,
				})
				id, err := result.Get(ctx)
				if err != nil {
					return errors.Wrapf(err, "failed to publish %v", log.Message.Attributes)
				}
				fmt.Printf("published: id=%s, attributes=%v\n", id, log.Attributes)
			}
			return errors.WithStack(publisher.Close())
		},
	}
	c.Flags().StringVar(&publishOption.project, "project", "", "pubsub project name. If not set, publish to the original project described in the dead-letter log.")
	c.Flags().StringVar(&publishOption.topic, "topic", "", "pubsub topic name. If not set, publish to the original topic described in the dead-letter log.")
	return &c
}
