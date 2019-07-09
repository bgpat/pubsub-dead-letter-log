package cmd

import (
	"context"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/wantedly/pubsub-dead-letter-log"
)

func newListCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "list",
		Short: "list log files from URL",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
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

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"timestamp", "project", "topic", "attributes"})
			table.SetRowLine(true)
			for _, l := range logs {
				attrs := []string{}
				for k, v := range l.Attributes {
					attrs = append(attrs, k+"="+v)
				}
				table.Append([]string{l.Timestamp.String(), l.Project, l.Topic, strings.Join(attrs, "\n")})
			}
			table.Render()
			return nil
		},
	}
	return &c
}
