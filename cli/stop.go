package cli

import (
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func stopCommand(t *core.Timetrace) *cobra.Command {
	stop := &cobra.Command{
		Use:   "stop",
		Short: "Stop tracking your time",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if err := t.Stop(); err != nil {
				out.Err("Failed to stop tracking: %s", err.Error())
				return
			}

			out.Success("Stopped tracking time")
		},
	}

	return stop
}
