package cli

import (
	"time"

	"github.com/dominikbraun/timetrace/config"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

const (
	defaultRecordArgLayout = "2006-01-02-15-04"
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

			out.Table([]string{"Key"}, [][]string{{project.Key}})
		},
	}

	return getProject
}

func getRecordCommand(t *core.Timetrace) *cobra.Command {
	getRecord := &cobra.Command{
		Use:   "record YYYY-MM-DD-HH-MM",
		Short: "Display a record",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			layout := defaultRecordArgLayout

			if config.Get().Use12Hours {
				layout = "2006-01-02-03-04PM"
			}

			start, err := time.Parse(layout, args[0])
			if err != nil {
				out.Err("Failed to parse date argument: %s", err.Error())
				return
			}

			record, err := t.LoadRecord(start)
			if err != nil {
				out.Err("Failed to read record: %s", err.Error())
				return
			}

			isBillable := defaultBool

			if record.IsBillable {
				isBillable = "yes"
			}

			end := defaultString

			if record.End != nil {
				end = record.End.Format("15:04")
			}

			project := defaultString

			if record.Project != nil {
				project = record.Project.Key
			}

			rows := [][]string{
				{
					record.Start.Format("15:04"),
					end,
					project,
					isBillable,
				},
			}

			out.Table([]string{"Start", "End", "Project", "Billable"}, rows)
		},
	}

	return getRecord
}
