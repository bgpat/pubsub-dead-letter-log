package cmd

import (
	"context"
	"os"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/wantedly/psdll"
)

func newListCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "list <url> ...",
		Short: "list log files from URL",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			ctx := context.Background()
			type pair struct {
				log  psdll.DeadLetterLog
				path string
			}
			var logs []pair
			for _, arg := range args {
				l, err := psdll.ReadFromURL(ctx, arg)
				if err != nil {
					return errors.Wrapf(err, "invalid URL: %q", arg)
				}
				for path, log := range l {
					logs = append(logs, pair{
						log:  log,
						path: path,
					})
				}
			}
			sort.Slice(logs, func(i, j int) bool {
				return logs[i].log.Timestamp.Before(logs[j].log.Timestamp)
			})

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"path", "timestamp", "project", "topic", "attributes"})
			table.SetRowLine(true)
			for _, l := range logs {
				attrs := []string{}
				for k, v := range l.log.Attributes {
					attrs = append(attrs, k+"="+v)
				}
				table.Append([]string{
					l.path,
					l.log.Timestamp.String(),
					l.log.Project,
					l.log.Topic,
					strings.Join(attrs, "\n")},
				)
			}
			table.Render()
			return nil
		},
	}
	return &c
}
