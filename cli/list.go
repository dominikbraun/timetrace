package cli

import (
	"strconv"
	"strings"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func listCommand(t *core.Timetrace) *cobra.Command {
	list := &cobra.Command{
		Use:   "list",
		Short: "List all resources",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	list.AddCommand(listProjectsCommand(t))
	list.AddCommand(listRecordsCommand(t))

	return list
}

func listProjectsCommand(t *core.Timetrace) *cobra.Command {
	listProjects := &cobra.Command{
		Use:   "projects",
		Short: "List all projects",
		Run: func(cmd *cobra.Command, args []string) {
			projects, err := t.ListProjects()
			if err != nil {
				out.Err("Failed to list projects: %s", err.Error())
				return
			}

			rows := make([][]string, len(projects))

			for i, project := range projects {
				rows[i] = make([]string, 2)
				rows[i][0] = strconv.Itoa(i + 1)
				rows[i][1] = project.Key
			}

			out.Table([]string{"#", "Key"}, rows, nil)
		},
	}

	return listProjects
}

type listRecordsOptions struct {
	isOnlyDisplayingBillable bool
	projectKeyFilter         string
}

func listRecordsCommand(t *core.Timetrace) *cobra.Command {
	var options listRecordsOptions

	listRecords := &cobra.Command{
		Use:   "records {<YYYY-MM-DD>|today|yesterday}",
		Short: "List all records from a date",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			date, err := t.Formatter().ParseDate(args[0])
			if err != nil {
				out.Err("failed to parse date: %s", err.Error())
				return
			}

			records, err := t.ListRecords(date)
			if err != nil {
				out.Err("failed to list records: %s", err.Error())
				return
			}

			if len(options.projectKeyFilter) > 0 {
				records = filterProjectRecords(records, options.projectKeyFilter)
			}

			if options.isOnlyDisplayingBillable {
				records = filterBillableRecords(records)
			}

			rows := make([][]string, len(records))

			for i, record := range records {
				end := defaultString
				if record.End != nil {
					end = t.Formatter().TimeString(*record.End)
				}

				billable := defaultBool

				if record.IsBillable {
					billable = "yes"
				}

				rows[i] = make([]string, 6)
				rows[i][0] = strconv.Itoa(i + 1)
				rows[i][1] = t.Formatter().RecordKey(record)
				rows[i][2] = record.Project.Key
				rows[i][3] = t.Formatter().TimeString(record.Start)
				rows[i][4] = end
				rows[i][5] = billable
			}

			out.Table([]string{"#", "Key", "Project", "Start", "End", "Billable"}, rows, nil)
		},
	}

	listRecords.Flags().BoolVarP(&options.isOnlyDisplayingBillable, "billable", "b",
		false, `only display billable records`)

	listRecords.Flags().StringVarP(&options.projectKeyFilter, "project", "p",
		"", "filter by project key")

	return listRecords
}

func filterBillableRecords(records []*core.Record) []*core.Record {
	billableRecords := []*core.Record{}
	for _, record := range records {
		if record.IsBillable == true {
			billableRecords = append(billableRecords, record)
		}
	}
	return billableRecords
}

func filterProjectRecords(records []*core.Record, projectKey string) []*core.Record {
	projectRecords := []*core.Record{}
	for _, record := range records {
		if strings.Compare(record.Project.Key, projectKey) == 0 {
			projectRecords = append(projectRecords, record)
		}
	}
	return projectRecords
}
