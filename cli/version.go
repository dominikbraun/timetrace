package cli

import (
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func versionCommand(value string) *cobra.Command {
	version := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		Run: func(cmd *cobra.Command, args []string) {
			out.Info("timetrace version %s", value)
		},
	}

	return version
}
