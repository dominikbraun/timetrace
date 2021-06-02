package cli

import (
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"
	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

type reportOptions struct {
	isBillable   bool
	projectKey   string
	outputFormat string
	outputPath   string
	fromTime     string
	toTime       string
}

func generateReportCommand(t *core.Timetrace) *cobra.Command {
	var options reportOptions

	report := &cobra.Command{
		Use:   "report",
		Short: "Report allows to view or output tracked records as defined report",
		Run: func(cmd *cobra.Command, args []string) {
			var fromDate, toDate time.Time
			var formatErr error

			// TODO: find better way to check if flag is default
			// matching the default string will be buggy
			if options.fromTime != "from oldest" {
				fromDate, formatErr = t.Formatter().ParseDate(options.fromTime)
				if formatErr != nil {
					out.Err("failed to parse date: %s", formatErr.Error())
					return
				}
			}

			// TODO: find better way to check if flag is default
			// matching the default string will be buggy
			if options.toTime != "to newest" {
				toDate, formatErr = t.Formatter().ParseDate(options.toTime)
				if formatErr != nil {
					out.Err("failed to parse date: %s", formatErr.Error())
					return
				}
			}

			// set-up filter options based on cmd flags
			var filter = []func(*core.Record) bool{
				// this will ignore records which end time to not set
				// so current tracked times for example
				core.FilterNoneNilEndTime,
				core.FilterByTimeRange(fromDate, toDate),
			}
			// TODO: find better way to check if flag is default
			// matching the default string will be buggy
			if options.projectKey != "include all projects" {
				filter = append(filter, core.FilterByProject(options.projectKey))
			}
			if options.isBillable {
				filter = append(filter, core.FilterBillable)
			}

			report, err := t.Report(filter...)
			if err != nil {
				out.Err(err.Error())
			}

			// check what to do with the report
			// if options.outputFormat is default only table will be
			// printed to os.Stdout
			switch options.outputFormat {
			case "json":
				data, err := report.Json()
				if err != nil {
					out.Err(err.Error())
				}
				t.WriteReport(options.outputPath, data)
			default:
				projects, total := report.Table()
				out.Table(
					[]string{"Project", "Date", "Start", "End", "Billable", "Total"},
					projects,
					[]string{"", "", "", "", "TOTAL", total},
					out.TableWithCellMerge(0), // merge cells over "Project" (index:0) column
					out.TableFooterColor(
						tablewriter.Colors{}, tablewriter.Colors{},
						tablewriter.Colors{}, tablewriter.Colors{},
						tablewriter.Colors{tablewriter.Bold},          // text "TOTAL"
						tablewriter.Colors{tablewriter.FgGreenColor}), // digit of "TOTAL"
				)
			}
		},
	}

	report.Flags().BoolVarP(&options.isBillable, "billable", "b",
		false, "filter for only billable records")

	report.Flags().StringVarP(&options.fromTime, "start", "s",
		"from oldest", "filter records from a given start date <YYYY-MM-DD>")

	report.Flags().StringVarP(&options.toTime, "end", "e",
		"to newest", "filter records to a given end date (end is inclusive) <YYYY-MM-DD>")

	report.Flags().StringVarP(&options.projectKey, "project", "p",
		"include all projects", "filter records by a specific project")

	report.Flags().StringVarP(&options.outputFormat, "format", "f",
		"print table", "output format for report (json/csv)")

	report.Flags().StringVarP(&options.outputPath, "out", "o",
		"", "choose output path for report")

	return report
}
