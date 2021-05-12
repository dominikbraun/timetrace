package cli

import (
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func statusCommand() *cobra.Command {
	status := &cobra.Command{
		Use:   "status",
		Short: "Display the current tracking status",
		Run: func(cmd *cobra.Command, args []string) {
			report, err := core.Status()
			if err != nil {
				out.Err("Failed to obtain status: %s", err.Error())
				return
			}

			if report == nil {
				out.Info("You're not tracking time at the moment")
				return
			}

			project := defaultString

			if report.Current.Project != nil {
				project = report.Current.Project.Key
			}

			rows := [][]string{
				{
					project,
					report.TrackedTimeCurrent.String(),
					report.TrackedTimeToday.String(),
				},
			}

			out.Table([]string{"Current project", "Worked since start", "Worked today"}, rows)
		},
	}

	return status
}
