package main

import (
	"github.com/dominikbraun/timetrace/plugin"
	"strings"
)

func main() {
	p := plugin.New("HelloWorld")
	err := p.Register()
	if err != nil {
		panic(err)
	}

	logger := p.Logger()
	logger.Println("Initializing")

	p.RegisterCobraCommand(plugin.RegisterCobraCommand{
		Use:     "hello",
		Short:   "prints Hello World and all arguments",
		Example: "hello",
	})

	p.OnCommand(func(cmd plugin.OnCommand) error {
		if cmd.Cmd.Use == "hello" {
			return p.Print("Hello World" + strings.Join(cmd.Args, ", "))
		}

		return nil
	})

	p.OnAllInitialized(func() error {
		return nil
	})

	err = p.Run()
	if err != nil {
		logger.Println(err)
	}
}
