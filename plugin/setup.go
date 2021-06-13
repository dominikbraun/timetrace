package plugin

import (
	"encoding/json"
	"github.com/aligator/goplug/goplug"
	"github.com/dominikbraun/timetrace/config"
	"github.com/dominikbraun/timetrace/core"
	"github.com/spf13/cobra"
)

type Host struct {
	commands []RegisterCobraCommand
	action   goplug.OnOneShot
	T        *core.Timetrace
}

func (h *Host) RegisterOneShot(info goplug.PluginInfo, action goplug.OnOneShot) error {
	var meta Metadata
	err := json.Unmarshal(info.Metadata, &meta)
	if err != nil {
		return err
	}

	h.commands = meta.Commands
	h.action = action
	return nil
}

func (h *Host) Init(c *config.Config) error {
	if c.PluginFolder == "" {
		c.PluginFolder = "plugins"
	}

	g := goplug.GoPlug{
		PluginFolder: c.PluginFolder,
		Host:         h,
		Actions:      &Actions{h.T},
	}
	return g.Init()
}

func (h *Host) AddToCobra(root *cobra.Command) error {
	for _, pluginCmd := range h.commands {
		root.AddCommand(&cobra.Command{
			Use:     pluginCmd.Use,
			Short:   pluginCmd.Short,
			Long:    pluginCmd.Long,
			Example: pluginCmd.Example,
			Run: func(cmd *cobra.Command, args []string) {
				h.action(append([]string{cmd.Name()}, args...))
			},
		})
	}

	return nil
}
