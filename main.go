// Note: go install needs min Go 1.16
//go:generate go install github.com/aligator/goplug@v0.0.8
//go:generate goplug -o plugin/actions -allow-structs -allow-slices .
//go:generate go build -o ./plugins ./plugins/example/hello

package main

import (
	"github.com/dominikbraun/timetrace/plugin"
	"os"

	"github.com/dominikbraun/timetrace/cli"
	"github.com/dominikbraun/timetrace/config"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/fs"
	"github.com/dominikbraun/timetrace/out"
)

var version = "UNDEFINED"

func main() {
	c, err := config.FromFile()
	if err != nil {
		out.Warn("%s", err.Error())
	}

	filesystem := fs.New(c)
	timetrace := core.New(c, filesystem)

	pluginHost := &plugin.Host{
		T: timetrace,
	}

	err = pluginHost.Init(c)
	if err != nil {
		out.Err("%s", err.Error())
		os.Exit(1)
	}

	if err := cli.RootCommand(timetrace, version, pluginHost).Execute(); err != nil {
		out.Err("%s", err.Error())
		os.Exit(1)
	}
}
