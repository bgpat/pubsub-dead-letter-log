package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gocloud.dev/blob"

	"github.com/wantedly/pubsub-dead-letter-log"
)

func newPublishCommand() *cobra.Command {
	c := cobra.Command{
		Use:          "publish",
		Short:        "publish messages to Google Cloud Pub/Sub",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			ctx := context.Background()
			var logs []psdll.DeadLetterLog
			for _, arg := range args {
				bucket, files, err := retrieveFiles(ctx, arg)
				if err != nil {
					return errors.Wrapf(err, "failed to retrive %q", arg)
				}
				defer func() {
					if err := bucket.Close(); err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
				}()
				for _, file := range files {
					b, err := bucket.ReadAll(ctx, file)
					if err != nil {
						return errors.Wrapf(err, "failed to read file %q", file)
					}
					var log psdll.DeadLetterLog
					if err := json.Unmarshal(b, &log); err != nil {
						return errors.Wrapf(err, "%q is not a JSON file", file)
					}
					logs = append(logs, log)
				}
			}
			sort.Slice(logs, func(i, j int) bool {
				return logs[i].Timestamp.Before(logs[j].Timestamp)
			})
			clients := make(map[string]*pubsub.Client)
			topics := make(map[string]*pubsub.Topic)
			for _, log := range logs {
				var err error
				client, cached := clients[log.Project]
				if !cached {
					client, err = pubsub.NewClient(ctx, log.Project)
					if err != nil {
						return errors.Wrapf(err, "failed to initialize pubsub client for %q project", log.Project)
					}
					clients[log.Project] = client
				}
				topic, cached := topics[log.Project+"/"+log.Topic]
				if !cached {
					topic = client.Topic(log.Topic)
					topics[log.Project+"/"+log.Topic] = topic
				}
				result := topic.Publish(ctx, &pubsub.Message{Data: log.Message.Data, Attributes: log.Message.Attributes})
				id, err := result.Get(ctx)
				if err != nil {
					return errors.Wrapf(err, "failed to publish %v", log.Message.Attributes)
				}
				fmt.Printf("published: id=%s, attributes=%v\n", id, log.Attributes)
			}
			return nil
		},
	}
	return &c
}

func retrieveFiles(ctx context.Context, dir string) (*blob.Bucket, []string, error) {
	bucket, err := blob.OpenBucket(ctx, dir)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	u, err := url.Parse(dir)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	key := strings.TrimLeft(u.Path, "/")
	if len(key) > 0 {
		if exists, err := bucket.Exists(ctx, key); err != nil {
			return nil, nil, errors.WithStack(err)
		} else if exists {
			return bucket, []string{key}, nil
		}
	}
	var files []string
	i := bucket.List(&blob.ListOptions{
		Prefix: key,
	})
	for {
		obj, err := i.Next(ctx)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		if isMatchFileOrDir(key, obj.Key) {
			files = append(files, obj.Key)
		}
	}
	return bucket, files, nil
}

func isMatchFileOrDir(query, actual string) bool {
	a := strings.Split(actual, "/")
	for i, s := range strings.Split(query, "/") {
		if s != "" && s != a[i] {
			return false
		}
	}
	return true
}
