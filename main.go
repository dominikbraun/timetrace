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

	plugins := &plugin.Plugins{}
	plugins.Init(c)
	defer plugins.Close()

	filesystem := fs.New(c)
	timetrace := core.New(c, filesystem)

	if err := cli.RootCommand(timetrace, version, plugins).Execute(); err != nil {
		out.Err("%s", err.Error())
		os.Exit(1)
	}
}
