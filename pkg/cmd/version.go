package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	appName = "file2pubsub"
	version = "v0.0.0"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:           "version",
		Short:         "Print the version information",
		SilenceErrors: true,
		SilenceUsage:  true,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf(
				"%s %s (%s %s/%s)\n",
				appName,
				version,
				runtime.Version(),
				runtime.GOOS,
				runtime.GOARCH,
			)
		},
	}
}
