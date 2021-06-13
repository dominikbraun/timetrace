package plugin

import (
	"encoding/json"
	"github.com/aligator/goplug/goplug"
	"github.com/aligator/goplug/plugin"
	"os"
)

type RegisterCobraCommand struct {
	Use     string                    `json:"use"`
	Short   string                    `json:"short"`
	Long    string                    `json:"long"`
	Example string                    `json:"example"`
	Action  func(args []string) error `json:"-"`
}

type Metadata struct {
	Commands []RegisterCobraCommand `json:"commands"`
}

// Timetrace provides the methods, which can be used by plugins.
type Timetrace struct {
	plugin   plugin.Plugin
	commands []RegisterCobraCommand
}

func New(id string) Timetrace {
	return Timetrace{
		plugin: plugin.Plugin{
			PluginInfo: goplug.PluginInfo{
				ID:         id,
				PluginType: goplug.OneShot,
			},
		},
	}
}

func (t *Timetrace) RegisterCobraCommand(cmd RegisterCobraCommand) {
	t.commands = append(t.commands, cmd)
}

func (t *Timetrace) Run() error {
	metaJson, err := json.Marshal(Metadata{
		Commands: t.commands,
	})
	if err != nil {
		return err
	}
	t.plugin.Metadata = metaJson

	t.plugin.Init()

	for _, cmd := range t.commands {
		if cmd.Use != os.Args[1] {
			continue
		}

		err := cmd.Action(os.Args[1:])
		if err != nil {
			panic(err)
		}
	}

	return nil
}
