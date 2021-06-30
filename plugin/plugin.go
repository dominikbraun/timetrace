package plugin

import (
	"encoding/json"
	"github.com/aligator/goplug/goplug"
	"github.com/dominikbraun/timetrace/plugin/actions"
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

// Plugin provides the methods, which can be used by plugins.
type Plugin struct {
	actions.ClientActions
	client   *goplug.Client
	commands []RegisterCobraCommand
}

func New(info goplug.PluginInfo) Plugin {

	client := &goplug.Client{
		PluginInfo: info,
	}
	return Plugin{
		ClientActions: actions.NewClientActions(client),
		client:        client,
	}
}

func (t *Plugin) RegisterCobraCommand(cmd RegisterCobraCommand) {
	t.commands = append(t.commands, cmd)
}

func (t *Plugin) Run() error {
	metaJson, err := json.Marshal(Metadata{
		Commands: t.commands,
	})
	if err != nil {
		return err
	}
	t.client.Metadata = string(metaJson)

	t.client.Init()

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
