package main

import (
	"os"

	"github.com/dominikbraun/timetrace/cli"
	"github.com/dominikbraun/timetrace/out"
)

var version = "UNDEFINED"

func main() {
	if err := cli.RootCommand(version).Execute(); err != nil {
		out.Err("%s", err.Error())
		os.Exit(1)
	}
}
