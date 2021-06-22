package actions

import (
	"github.com/aligator/goplug/goplug"
	core0 "github.com/dominikbraun/timetrace/core"
)

// HostActions contains the host-implementations of actions.
type HostActions struct {
	Core0TimetraceRef *core0.Timetrace
}

type ClientActions struct {
	client *goplug.Client
}

func NewClientActions(plugin *goplug.Client) ClientActions {
	return ClientActions{
		client: plugin,
	}
}

// Make some plugin-methods available to the plugins.

func (c *ClientActions) Print(text string) error {
	return c.client.Print(text)
}

// Action implementations for host and client.

type LoadLatestRecordRequest struct {
}

type LoadLatestRecordResponse struct {
	Res0 *core0.Record `json:"res0"`
}

// LoadLatestRecord loads the youngest record. This may also be a record from
// another day. If there is no latest record, nil and no error will be returned.
func (h *HostActions) LoadLatestRecord(args LoadLatestRecordRequest, reply *LoadLatestRecordResponse) error {
	// Host implementation.
	res0, err := h.Core0TimetraceRef.LoadLatestRecord()

	if err != nil {
		return err
	}

	*reply = LoadLatestRecordResponse{
		Res0: res0,
	}

	return nil
}

// LoadLatestRecord loads the youngest record. This may also be a record from
// another day. If there is no latest record, nil and no error will be returned.
func (c *ClientActions) LoadLatestRecord() (res0 *core0.Record, err error) {
	// Calling from the plugin.
	response := LoadLatestRecordResponse{}
	err = c.client.Call("LoadLatestRecord", LoadLatestRecordRequest{}, &response)
	return response.Res0, err
}
