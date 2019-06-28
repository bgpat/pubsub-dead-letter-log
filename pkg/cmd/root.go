package cmd

import "github.com/spf13/cobra"

// New returns cobra.Command object.
func New() *cobra.Command {
	c := cobra.Command{
		Use: "psdll",
	}
	c.AddCommand(
		newVersionCommand(),
		newPublishCommand(),
	)
	return &c
}
