package plugin

import (
	"github.com/dominikbraun/timetrace/core"
)

// Actions contains the host-implementations of actions.
type Actions struct {
	t *core.Timetrace
}

func (a *Actions) LoadLatestRecord(args interface{}, reply *core.Record) error {
	latestRecord, err := a.t.LoadLatestRecord()
	if err != nil {
		return err
	}

	*reply = *latestRecord
	return nil
}

func (t *Timetrace) LoadLatestRecord() *core.Record {
	response := &core.Record{}
	err := t.plugin.Call("LoadLatestRecord", nil, &response)
	if err != nil {
		panic(err)
	}

	return response
}

func (t *Timetrace) Print(text string) {
	// Just pass through to make it available for plugins.
	t.plugin.Print(text)
}
