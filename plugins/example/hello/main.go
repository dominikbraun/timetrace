package main

import (
	"fmt"
	"github.com/aligator/goplug/goplug"
	"github.com/dominikbraun/timetrace/plugin"
	"strings"
)

type HelloPlugin struct {
	plugin.Plugin
}

func New() HelloPlugin {
	return HelloPlugin{
		Plugin: plugin.New(goplug.PluginInfo{
			ID:         "superplugin",
			PluginType: goplug.OneShot,
		}),
	}
}

func main() {
	p := New()
	p.RegisterCobraCommand(plugin.RegisterCobraCommand{
		Use:     "hello",
		Short:   "prints Hello and all arguments",
		Example: "hello",
		Action: func(args []string) error {
			p.Print(strings.Join(args, ", ") + "\n")
			return nil
		},
	})

	p.RegisterCobraCommand(plugin.RegisterCobraCommand{
		Use:     "record",
		Short:   "prints the current record",
		Example: "record",
		Action: func(args []string) error {
			r, err := p.LoadLatestRecord()
			if err != nil {
				return err
			}

			p.Print(fmt.Sprintf("Latest record: %v, %v, %v, %v", r.Project.Key, r.Start, r.End, r.IsBillable) + "\n")
			return nil
		},
	})

	p.Run()
}
