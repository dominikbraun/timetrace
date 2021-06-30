package cli

import (
	"fmt"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func getCommand(t *core.Timetrace) *cobra.Command {
	get := &cobra.Command{
		Use:   "get",
		Short: "Display a resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	get.AddCommand(getProjectCommand(t))
	get.AddCommand(getRecordCommand(t))

	return get
}

func getProjectCommand(t *core.Timetrace) *cobra.Command {
	getProject := &cobra.Command{
		Use:   "project <KEY>",
		Short: "Display a project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]

			project, err := t.LoadProject(key)
			if err != nil {
				out.Err("Failed to get project: %s", key)
				return
			}

			out.Table([]string{"Key"}, [][]string{{project.Key}}, nil)
		},
	}

	return getProject
}

func getRecordCommand(t *core.Timetrace) *cobra.Command {

	// Depending on the use12hours setting, the command syntax either is
	// `record YYYY-MM-DD-HH-MM` or `record YYYY-MM-DD-HH-MMPM`.
	use := fmt.Sprintf("record %s", t.Formatter().RecordKeyLayout())

	getRecord := &cobra.Command{
		Use:   use,
		Short: "Display a record",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			start, err := t.Formatter().ParseRecordKey(args[0])
			if err != nil {
				out.Err("Failed to parse date argument: %s", err.Error())
				return
			}

			record, err := t.LoadRecord(start)
			if err != nil {
				out.Err("Failed to read record: %s", err.Error())
				return
			}

			showRecord(&record, t.Formatter())
		},
	}

	return getRecord
}

func showRecord(record *core.Record, formatter *core.Formatter) {
	isBillable := defaultBool

	if record.IsBillable {
		isBillable = "yes"
	}

	end := defaultString
	if record.End != nil {
		end = formatter.TimeString(*record.End)
	}

	project := defaultString

	if record.Project != nil {
		project = record.Project.Key
	}

	rows := [][]string{
		{
			formatter.TimeString(record.Start),
			end,
			project,
			isBillable,
		},
	}

	out.Table([]string{"Start", "End", "Project", "Billable"}, rows, nil)
}
