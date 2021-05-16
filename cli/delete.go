package cli

import (
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func deleteCommand(t *core.Timetrace) *cobra.Command {
	delete := &cobra.Command{
		Use:   "delete",
		Short: "Delete a resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	delete.AddCommand(deleteProjectCommand(t))
	delete.AddCommand(deleteRecordCommand(t))

	return delete
}

func deleteProjectCommand(t *core.Timetrace) *cobra.Command {
	deleteProject := &cobra.Command{
		Use:   "project <KEY>",
		Short: "Delete a project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]

			project := core.Project{
				Key: key,
			}

			if err := t.DeleteProject(project); err != nil {
				out.Err("Failed to delete %s", err.Error())
				return
			}

			out.Success("Deleted project %s", key)
		},
	}

	return deleteProject
}

func deleteRecordCommand(t *core.Timetrace) *cobra.Command {
	deleteRecord := &cobra.Command{
		Use:   "record YYYY-MM-DD-HH-MM",
		Short: "Delete a record",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			layout := defaultRecordArgLayout

			if t.Config().Use12Hours {
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

			if err := t.DeleteRecord(*record); err != nil {
				out.Err("Failed to delete %s", err.Error())
				return
			}

			out.Success("Deleted record %s", args[0])
		},
	}

	return deleteRecord
}
