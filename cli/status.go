package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func statusCommand(t *core.Timetrace) *cobra.Command {
	var format string

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
				trackedTimeCurrent = t.Formatter().FormatCurrentTime(report)
			}

			rows := [][]string{
				{
					project,
					trackedTimeCurrent,
					t.Formatter().FormatTodayTime(report),
					t.Formatter().FormatBreakTime(report),
				},
			}
			if format != "" {
				format = strings.ReplaceAll(format, "{project}", project)
				format = strings.ReplaceAll(format, "{trackedTimeCurrent}", trackedTimeCurrent)
				format = strings.ReplaceAll(format, "{trackedTimeToday}", t.Formatter().FormatTodayTime(report))
				format = strings.ReplaceAll(format, "{breakTimeToday}", t.Formatter().FormatBreakTime(report))
				format = strings.ReplaceAll(format, `\n`, "\n")
				fmt.Printf(format)
				return
			}

			out.Table([]string{"Current project", "Worked since start", "Worked today", "Breaks"}, rows, nil)
		},
	}

	status.Flags().StringVarP(&format, "format", "f", "", "Format string, availiable:\n{project}, {trackedTimeCurrent}, {todayTime}, {breakTime}")

	return status
}
