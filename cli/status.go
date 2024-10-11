package cli

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

type statusReport struct {
	Project            string `json:"project"`
	TrackedTimeCurrent string `json:"trackedTimeCurrent"`
	TrackedTimeToday   string `json:"trackedTimeToday"`
	BreakTimeToday     string `json:"breakTimeToday"`
}

type statusOptions struct {
	format string
	output string
}

func statusCommand(t *core.Timetrace) *cobra.Command {
	var options statusOptions

	status := &cobra.Command{
		Use:   "status",
		Short: "Display the current tracking status",
		Run: func(cmd *cobra.Command, args []string) {
			report, err := t.Status()
			if errors.Is(err, core.ErrTrackingNotStarted) {
				out.Info("You haven't started tracking time today")
				return
			} else if err != nil {
				out.Err("failed to obtain status: %s", err.Error())
				return
			}

			if report == nil {
				out.Info("you're not tracking time at the moment")
				return
			}

			statusReport := statusReport{
				Project:            defaultString,
				TrackedTimeCurrent: defaultString,
				TrackedTimeToday:   t.Formatter().FormatDuration(report.TrackedTimeToday),
				BreakTimeToday:     t.Formatter().FormatDuration(report.BreakTimeToday),
			}

			if report.Current != nil {
				statusReport.Project = report.Current.Project.Key
			}

			if report.TrackedTimeCurrent != nil {
				statusReport.TrackedTimeCurrent = t.Formatter().FormatDuration(*report.TrackedTimeCurrent)
			}

			if options.format != "" {
				format := options.format
				format = strings.ReplaceAll(format, "{project}", statusReport.Project)
				format = strings.ReplaceAll(format, "{trackedTimeCurrent}", statusReport.TrackedTimeCurrent)
				format = strings.ReplaceAll(format, "{trackedTimeToday}", statusReport.TrackedTimeToday)
				format = strings.ReplaceAll(format, "{breakTimeToday}", statusReport.BreakTimeToday)
				format = strings.ReplaceAll(format, `\n`, "\n")
				fmt.Printf(format)
				return
			}

			if options.output != "" {
				switch options.output {
				case "json":
					bytes, err := json.MarshalIndent(statusReport, "", "\t")
					if err != nil {
						out.Err("error printing JSON: %s", err.Error())
						return
					}
					fmt.Println(string(bytes))
					return
				case "csv":
					var buf bytes.Buffer
					writer := csv.NewWriter(&buf)

					headers := []string{"Project", "Tracked Time Current", "Tracked Time Today", "Break Time Today"}
					if err := writer.Write(headers); err != nil {
						out.Err("error writing CSV header: %s", err.Error())
						return
					}

					row := []string{
						statusReport.Project,
						statusReport.TrackedTimeCurrent,
						statusReport.TrackedTimeToday,
						statusReport.BreakTimeToday,
					}
					if err := writer.Write(row); err != nil {
						out.Err("error writing CSV row: %s", err.Error())
						return
					}

					writer.Flush()

					if err := writer.Error(); err != nil {
						out.Err("error flushing CSV: %s", err.Error())
						return
					}

					csvBytes := buf.Bytes()

					fmt.Println(string(csvBytes))
					return
				default:
					out.Err("unknown output format: %s", options.format)
					return
				}
			}

			rows := [][]string{
				{
					statusReport.Project,
					statusReport.TrackedTimeCurrent,
					statusReport.TrackedTimeToday,
					statusReport.BreakTimeToday,
				},
			}

			out.Table([]string{"Current project", "Worked since start", "Worked today", "Breaks"}, rows, nil)
		},
	}

	status.Flags().StringVarP(&options.format, "format", "f", "", "Format string, availiable:\n{project}, {trackedTimeCurrent}, {trackedTimeToday}, {breakTimeToday}")
	status.Flags().StringVarP(&options.output, "output", "o", "", "The output format. Available: json")

	return status
}
