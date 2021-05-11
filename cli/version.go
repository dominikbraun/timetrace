package cli

import (
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func versionCommand(value string) *cobra.Command {
	version := &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			out.Err("timetrace version %s", value)
		},
	}

	return version
}
