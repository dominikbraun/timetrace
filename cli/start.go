package cli

import (
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

type startOptions struct {
	isBillable bool
}

func startCommand(t *core.Timetrace) *cobra.Command {
	var options startOptions

	start := &cobra.Command{
		Use:   "start <PROJECT KEY>",
		Short: "Start tracking time",
		Run: func(cmd *cobra.Command, args []string) {
			projectKey := args[0]

			if err := t.Start(projectKey, options.isBillable); err != nil {
				out.Err("Failed to start tracking: %s", err.Error())
				return
			}

			out.Success("Started tracking time")
		},
	}

	start.Flags().BoolVarP(&options.isBillable, "billable", "b",
		false, `mark tracked time as billable`)

	return start
}
