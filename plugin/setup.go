package plugin

import (
	"fmt"
	"github.com/aligator/goplug"
	"github.com/dominikbraun/timetrace/config"
	"github.com/spf13/cobra"
)

type RegisterCobraCommand struct {
	// ToDo: PluginID will become private or removed completely
	//       when commands get only sent to one specific plugin.
	PluginID string `json:"pluginId"`
	Use      string `json:"use"`
	Short    string `json:"short"`
	Long     string `json:"long"`
	Example  string `json:"example"`
}

type OnCommand struct {
	Cmd  RegisterCobraCommand `json:"cmd"`
	Args []string             `json:"args"`
}

type Plugins struct {
	GoPlug   *goplug.GoPlug
	commands []RegisterCobraCommand
}

func (p *Plugins) Init(c *config.Config) error {
	if c.PluginFolder == "" {
		c.PluginFolder = "plugins"
	}

	p.GoPlug = &goplug.GoPlug{
		PluginFolder: c.PluginFolder,
	}

	// Just a simple print command.
	p.GoPlug.RegisterOnCommand("print", func() interface{} {
		var s string
		return &s
	}, func(info goplug.PluginInfo, message interface{}) error {
		text := message.(*string)
		fmt.Print(*text)
		return nil
	})

	p.GoPlug.RegisterOnCommand("registerCobraCommand", func() interface{} {
		return &RegisterCobraCommand{}
	}, func(info goplug.PluginInfo, message interface{}) error {
		cobraData := message.(*RegisterCobraCommand)
		cobraData.PluginID = info.ID

		p.commands = append(p.commands, *cobraData)
		return nil
	})

	err := p.GoPlug.Init()
	if err != nil {
		return err
	}

	return nil
}

func (p *Plugins) AddToCobra(root *cobra.Command) error {
	for _, pluginCmd := range p.commands {
		root.AddCommand(&cobra.Command{
			Use:     pluginCmd.Use,
			Short:   pluginCmd.Short,
			Long:    pluginCmd.Long,
			Example: pluginCmd.Example,
			Run: func(cmd *cobra.Command, args []string) {
				// ToDo: for now it is sent to all plugins.
				//       This has to be changed, so only the plugin which
				//       registered the command receives it.
				p.GoPlug.Send("command", OnCommand{
					Cmd:  pluginCmd,
					Args: args,
				})
			},
		})
	}

	return nil
}

func (p *Plugins) Close() error {
	return p.GoPlug.Close()
}
