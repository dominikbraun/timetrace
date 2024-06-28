package cli

import (
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"
	"github.com/olekukonko/tablewriter"

	"github.com/spf13/cobra"
)

type reportOptions struct {
	isBillable    bool
	isNonBillable bool
	projectKey    string
	outputFormat  string
	filePath      string
	startTime     string
	endTime       string
}

func generateReportCommand(t *core.Timetrace) *cobra.Command {
	var options reportOptions

	report := &cobra.Command{
		Use:    "report",
		Short:  "Report allows to view or output tracked records as defined report",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			var startDate, endDate time.Time
			var formatErr error

			if options.startTime != "" {
				startDate, formatErr = t.Formatter().ParseDate(options.startTime)
				if formatErr != nil {
					out.Err("failed to parse date: %s", formatErr.Error())
					return
				}
			}

			if options.endTime != "" {
				endDate, formatErr = t.Formatter().ParseDate(options.endTime)
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
				core.FilterByTimeRange(startDate, endDate),
			}

			if options.projectKey != "" {
				filter = append(filter, core.FilterByProject(options.projectKey))
			}
			// wont hurt table will just be empty but makes sense to let the user know
			if options.isBillable && options.isNonBillable {
				out.Err("cannot filter for billable and none billable records")
				return
			}
			if options.isBillable {
				filter = append(filter, core.FilterBillable(true))
			}
			if options.isNonBillable {
				filter = append(filter, core.FilterBillable(false))
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
				t.WriteReport(options.filePath, data)
			default:
				projects, total := report.Table()
				out.Table(
					[]string{"Project", "Module", "Date", "Start", "End", "Billable", "Notes", "Total"},
					projects,
					[]string{"", "", "", "", "", "", "TOTAL", total},
					out.TableWithCellMerge(0), // merge cells over "Project" (index:0) column
					out.TableFooterColor(
						tablewriter.Colors{}, tablewriter.Colors{},
						tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{}, tablewriter.Colors{},
						tablewriter.Colors{tablewriter.Bold},          // text "TOTAL"
						tablewriter.Colors{tablewriter.FgGreenColor}), // digit of "TOTAL"
				)
			}
		},
	}

	report.Flags().BoolVarP(&options.isBillable, "billable", "b",
		false, "filter for only billable records")

	report.Flags().BoolVarP(&options.isNonBillable, "non-billable", "B",
		false, "filter for only none billable records")

	report.Flags().StringVarP(&options.startTime, "start", "s",
		"", "filter records from a given start date <YYYY-MM-DD>")

	report.Flags().StringVarP(&options.endTime, "end", "e",
		"", "filter records to a given end date (end is inclusive) <YYYY-MM-DD>")

	report.Flags().StringVarP(&options.projectKey, "project", "p",
		"", "filter records by a specific project")

	report.Flags().StringVarP(&options.outputFormat, "output", "o",
		"print table", "output format for report file (json)")

	report.Flags().StringVarP(&options.filePath, "file", "f",
		"", "file to write report to")

	return report
}
