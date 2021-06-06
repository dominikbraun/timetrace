package plugin

import "github.com/aligator/goplug"

// TimetracePlugin provides the methods, which can be used by plugins.
type TimetracePlugin struct {
	goplug.Plugin
}

func New(ID string) TimetracePlugin {
	return TimetracePlugin{
		Plugin: goplug.Plugin{
			ID: ID,
		},
	}
}

func (p *TimetracePlugin) OnCommand(listener func(cmd OnCommand) error) {
	p.RegisterCommand("command", func() interface{} {
		return &OnCommand{}
	}, func(message interface{}) error {
		data := message.(*OnCommand)
		if p.ID == data.Cmd.PluginID {
			return listener(*data)
		}

		return nil
	})
}

func (p TimetracePlugin) RegisterCobraCommand(cmd RegisterCobraCommand) error {
	return p.Send("registerCobraCommand", cmd)
}

func (p TimetracePlugin) Print(message string) error {
	return p.Send("print", message)
}
