package cli

import (
	"errors"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func statusCommand(t *core.Timetrace) *cobra.Command {
	status := &cobra.Command{
		Use:   "status",
		Short: "Display the current tracking status",
		Run: func(cmd *cobra.Command, args []string) {
			report, err := t.Status()
			if errors.Is(err, core.ErrTrackingNotStarted) {
				out.Info("You haven't started tracking time today")
				return
			} else if err != nil {
				out.Err("Failed to obtain status: %s", err.Error())
				return
			}

			if report == nil {
				out.Info("You're not tracking time at the moment")
				return
			}

			project := defaultString

			if report.Current != nil {
				project = report.Current.Project.Key
			}

			trackedTimeCurrent := defaultString

			if report.TrackedTimeCurrent != nil {
				trackedTimeCurrent = report.FormatCurrentTime()
			}

			rows := [][]string{
				{
					project,
					trackedTimeCurrent,
					report.FormatTodayTime(),
				},
			}
			out.Table([]string{"Current project", "Worked since start", "Worked today"}, rows, nil)
		},
	}

	return status
}
