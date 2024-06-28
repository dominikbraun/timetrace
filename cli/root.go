package cli

import (
	"github.com/dominikbraun/timetrace/core"

	"github.com/spf13/cobra"
)

const (
	defaultString = "---"
	defaultBool   = "no"
)

func RootCommand(t *core.Timetrace, version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "timetrace",
		Short:         "timetrace is a simple CLI for tracking your working time.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return t.EnsureDirectories()
		},
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	root.AddCommand(createCommand(t))
	root.AddCommand(attachCommand(t))
	root.AddCommand(getCommand(t))
	root.AddCommand(listCommand(t))
	root.AddCommand(editCommand(t))
	root.AddCommand(deleteCommand(t))
	root.AddCommand(startCommand(t))
	root.AddCommand(statusCommand(t))
	root.AddCommand(stopCommand(t))
	root.AddCommand(generateReportCommand(t))
	root.AddCommand(versionCommand(version))

	return root
}
