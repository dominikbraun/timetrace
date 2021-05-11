package cli

import (
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

type startOptions struct {
	isBillable bool
}

func startCommand() *cobra.Command {
	var options startOptions

	start := &cobra.Command{
		Use:  "start [KEY]",
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var projectKey string

			if len(args) > 0 {
				projectKey = args[0]
			}

			if err := core.Start(projectKey, options.isBillable); err != nil {
				out.Err("Failed to start tracking: %s", err.Error())
				return
			}

			out.Success("Started tracking time")
		},
	}

	start.Flags().BoolVarP(&options.isBillable, "billable", "b",
		false, `Mark work as billable`)

	return start
}
