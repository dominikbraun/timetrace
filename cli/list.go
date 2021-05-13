package cli

import (
	"strconv"

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
