package cli

import (
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func stopCommand() *cobra.Command {
	stop := &cobra.Command{
		Use:  "stop",
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := core.Stop(); err != nil {
				out.Err("Failed to stop tracking: %s", err.Error())
				return
			}

			out.Success("Stopped tracking time")
		},
	}

	return stop
}
