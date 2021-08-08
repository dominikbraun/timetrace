package cli

import (
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

type startOptions struct {
	isBillable    bool
	isNonBillable bool // Used for overwriting `billable: true` in the project config.
}

func startCommand(t *core.Timetrace) *cobra.Command {
	var options startOptions

	start := &cobra.Command{
		Use:   "start <PROJECT KEY>",
		Short: "Start tracking time",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			projectKey := args[0]

			isBillable := options.isBillable

			if projectConfig, ok := t.Config().Projects[projectKey]; ok {
				isBillable = projectConfig.Billable
			}

			if options.isNonBillable {
				isBillable = false
			}

			if err := t.Start(projectKey, isBillable); err != nil {
				out.Err("Failed to start tracking: %s", err.Error())
				return
			}

			out.Success("Started tracking time")
		},
	}

	start.Flags().BoolVarP(&options.isBillable, "billable", "b",
		false, `mark tracked time as billable`)

	start.Flags().BoolVar(&options.isNonBillable, "non-billable",
		false, `mark tracked time as non-billable if the project is configured as billable`)

	return start
}
