package main

import (
	"os"

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

	var jiraRepo *jira.Repository
	if c.JIRAIntegration != (config.JIRAConfig{}) {
		jiraRepo = jira.New(jira.RepositoryConfig{
			AuthToken:   c.JIRAIntegration.APIToken,
			Email:       c.JIRAIntegration.UserEmail,
			JIRAAddress: c.JIRAIntegration.Host,
		})
	}

	filesystem := fs.New(c)
	timetrace := core.New(c, filesystem, []core.Provider{jiraRepo})

	if err := cli.RootCommand(timetrace, version).Execute(); err != nil {
		out.Err("%s", err.Error())
		os.Exit(1)
	}
}
