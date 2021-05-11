package cli

import "github.com/spf13/cobra"

const (
	defaultString = "---"
	defaultBool   = "no"
)

func RootCommand(version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "timetrace",
		Short:         "timetrace is a simple CLI for tracking your working time.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	root.AddCommand(createCommand())
	root.AddCommand(getCommand())
	root.AddCommand(editCommand())
	root.AddCommand(deleteCommand())
	root.AddCommand(startCommand())
	root.AddCommand(statusCommand())
	root.AddCommand(stopCommand())
	root.AddCommand(versionCommand(version))

	return root
}
