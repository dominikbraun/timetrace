package cli

import (
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func editCommand() *cobra.Command {
	edit := &cobra.Command{
		Use:   "edit",
		Short: "Edit a resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	edit.AddCommand(editProjectCommand())

	return edit
}

func editProjectCommand() *cobra.Command {
	editProject := &cobra.Command{
		Use:   "project <KEY>",
		Short: "Edit a project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			out.Info("Opening %s in default editor", key)

			if err := core.EditProject(key); err != nil {
				out.Err("Failed to edit project: %s", err.Error())
				return
			}

			out.Success("Successfully edited %s", key)
		},
	}

	return editProject
}
