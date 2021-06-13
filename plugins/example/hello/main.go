package main

import (
	"fmt"
	"github.com/dominikbraun/timetrace/plugin"
	"strings"
)

type HelloPlugin struct {
	plugin.Timetrace
}

func New() HelloPlugin {
	return HelloPlugin{
		Timetrace: plugin.New("hello_plugin"),
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
			r := p.LoadLatestRecord()

			p.Print(fmt.Sprintf("Latest record: %v, %v, %v, %v", r.Project.Key, r.Start, r.End, r.IsBillable) + "\n")
			return nil
		},
	})

	p.Run()
}
