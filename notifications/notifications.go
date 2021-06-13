package notifications

import (
	"fmt"
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/gen2brain/beeep"
)

// Notifier provides methods to send notifications to the user
type Notifier struct {
	T     *core.Timetrace
	Timer time.Duration
}

// Run starts the Notifier, who sends notifactions to the user about the current record
// if it is counting time for freq minutes
func (n *Notifier) Run() {
	for {
		time.Sleep(10 * time.Second)

		rec, err := n.T.LoadLatestRecord()
		if err != nil || rec == nil {
			fmt.Printf(err.Error())
			continue
		}

		if rec.End == nil && rec.Start.Before(time.Now().Add(-n.Timer)) {
			err := beeep.Notify("timetrace", "Hey, are you still working on "+rec.Project.Key+"?", "")
			if err != nil {
				fmt.Printf(err.Error())
			}
			break
		}

		if rec.End != nil {
			break
		}
	}
}
