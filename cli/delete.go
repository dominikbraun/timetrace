package cli

import (
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func deleteCommand() *cobra.Command {
	delete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	delete.AddCommand(deleteProjectCommand())

	return delete
}

func deleteProjectCommand() *cobra.Command {
	deleteProject := &cobra.Command{
		Use:  "project <KEY>",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]

			project := core.Project{
				Key: key,
			}

			if err := core.DeleteProject(project); err != nil {
				out.Err("Failed to delete %s", err.Error())
				return
			}

			out.Success("Deleted project %s", key)
		},
	}

	return deleteProject
}
