package cli

import (
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

type reportOptions struct {
	isBillable   bool
	projectKey   string
	outputFormat string
	fromTime     string
	toTime       string
}

func generateReportCommand(t *core.Timetrace) *cobra.Command {
	var options reportOptions

	report := &cobra.Command{
		Use:   "report",
		Short: "Report allows to view or output tracked records",
		// Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var fromDate, toDate *time.Time
			var formatErr error

			// TODO: find better way to check if flag is not provided
			// matching the default string will be buggy
			if options.fromTime != "from oldest" {
				*fromDate, formatErr = t.Formatter().ParseDate(options.fromTime)
				if formatErr != nil {
					out.Err("failed to parse date: %s", formatErr.Error())
					return
				}
			}

			// TODO: find better way to check if flag is not provided
			// matching the default string will be buggy
			if options.toTime != "to newest" {
				*toDate, formatErr = t.Formatter().ParseDate(options.toTime)
				if formatErr != nil {
					out.Err("failed to parse date: %s", formatErr.Error())
					return
				}
			}

			// set-up filter options based on cmd flags
			var filter = []func(*core.Record) bool{
				core.FilterByTimeRange(fromDate, toDate),
			}
			// TODO: find better way to check for string default, this is to unstable
			if options.projectKey != "include all projects" {
				filter = append(filter, core.FilterByProject(options.projectKey))
			}
			if options.isBillable {
				filter = append(filter, core.FilterBillable)
			}

			_, err := t.Report(filter...)
			if err != nil {
				out.Err(err.Error())
			}
		},
	}

	report.Flags().BoolVarP(&options.isBillable, "billable", "b",
		false, "report only billable records")

	report.Flags().StringVarP(&options.fromTime, "start", "s",
		"from oldest", "filter records from a given start date <YYYY-MM-DD>")

	report.Flags().StringVarP(&options.toTime, "end", "e",
		"to newest", "filter records to a given end date (end is inclusive) <YYYY-MM-DD>")

	report.Flags().StringVarP(&options.projectKey, "project", "p",
		"include all projects", "filter records by a specific project")

	report.Flags().StringVarP(&options.outputFormat, "format", "f",
		"json", "choose output format for report")

	return report
}
