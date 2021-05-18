package cli

import (
	"strconv"
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

const (
	defaultTimeLayout = "15:04"
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

			out.Table([]string{"#", "Key"}, rows)
		},
	}

	return listProjects
}

func listRecordsCommand(t *core.Timetrace) *cobra.Command {
	listRecords := &cobra.Command{
		Use:   "records <YYYY-MM-DD>",
		Short: "List all records from a date",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			date, err := time.Parse("2006-01-02", args[0])
			if err != nil {
				out.Err("failed to parse date: %s", err.Error())
				return
			}

			records, err := t.ListRecords(date)
			if err != nil {
				out.Err("failed to list records: %s", err.Error())
				return
			}

			timeLayout := defaultTimeLayout
			keyLayout := defaultRecordArgLayout

			if t.Config().Use12Hours {
				timeLayout = "03:04PM"
				keyLayout = "2006-01-02-03-04PM"
			}

			rows := make([][]string, len(records))

			for i, record := range records {
				end := defaultString

				if record.End != nil {
					end = record.End.Format(timeLayout)
				}

				billable := defaultBool

				if record.IsBillable {
					billable = "yes"
				}

				rows[i] = make([]string, 6)
				rows[i][0] = strconv.Itoa(i + 1)
				rows[i][1] = record.Start.Format(keyLayout)
				rows[i][2] = record.Project.Key
				rows[i][3] = record.Start.Format(timeLayout)
				rows[i][4] = end
				rows[i][5] = billable
			}

			out.Table([]string{"#", "Key", "Project", "Start", "End", "Billable"}, rows)
		},
	}

	return listRecords
}
