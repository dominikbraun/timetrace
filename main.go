package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/dominikbraun/timetrace/cli"
	"github.com/dominikbraun/timetrace/config"
	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/fs"
	"github.com/dominikbraun/timetrace/integrations/jira"
	"github.com/dominikbraun/timetrace/out"
)

var version = "UNDEFINED"

func main() {
	c, err := config.FromFile()
	if err != nil {
		out.Warn("%s", err.Error())
	}

	// TODO: only load jira repo if specified
	rand.Seed(time.Now().UTC().UnixNano())
	jiraRepo := jira.New(jira.RepositoryConfig{})

	filesystem := fs.New(c)
	timetrace := core.New(c, filesystem, []core.Provider{jiraRepo})

	if err := cli.RootCommand(timetrace, version).Execute(); err != nil {
		out.Err("%s", err.Error())
		os.Exit(1)
	}
}
