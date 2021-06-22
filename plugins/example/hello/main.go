package main

import (
	"fmt"
	"github.com/aligator/goplug/goplug"
	core0 "github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/plugin"
	"strings"
	"time"
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

	p.RegisterCobraCommand(plugin.RegisterCobraCommand{
		Use:     "testhour",
		Short:   "saves a hour record to project 'test'",
		Example: "save super",
		Action: func(args []string) error {
			project, err := p.LoadProject("test")
			if err != nil {
				return err
			}

			end := time.Now().Add(1 * time.Hour)
			return p.SaveRecord(core0.Record{
				Start:      time.Now(),
				End:        &end,
				Project:    project,
				IsBillable: true,
			}, true)
		},
	})

	p.Run()
}
