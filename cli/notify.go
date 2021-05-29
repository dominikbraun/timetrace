package cli

import (
	"fmt"
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/notifications"
	"github.com/spf13/cobra"
)

func notifyCommand(t *core.Timetrace) *cobra.Command {
	notify := &cobra.Command{
		Use:   "notify",
		Short: "Send notification about running record",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			timer, err := time.ParseDuration(args[0])
			if err != nil {
				fmt.Printf("could not parse duration: %s", err)
			}
			notifier := notifications.Notifier{T: t, Timer: timer}
			notifier.Run()
		},
	}

	return notify
}
