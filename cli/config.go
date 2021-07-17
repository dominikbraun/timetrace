package cli

import (
	"strconv"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"
	"github.com/spf13/cobra"
)

func configCommand(t *core.Timetrace) *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configs",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	configCmd.AddCommand(setConfigCommand(t))
	configCmd.AddCommand(getConfigCommand(t))
	return configCmd
}

func setConfigCommand(t *core.Timetrace) *cobra.Command {
	config := &cobra.Command{
		Use:   "set <KEY> <VALUE>",
		Short: "Set config value",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	return config
}

func getConfigCommand(t *core.Timetrace) *cobra.Command {
	config := &cobra.Command{
		Use:   "get <KEY>",
		Short: "Get config values",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			configKeys := make([][]string, 4)
			configKeys[0] = createLine("store", t.Config().Store)
			configKeys[1] = createLine("editor", t.Config().Editor)
			configKeys[2] = createLine("reportPath", t.Config().ReportPath)
			configKeys[3] = createLine("use12Hours", strconv.FormatBool(t.Config().Use12Hours))
			if len(args) == 0 {

				out.Table([]string{"Key", "Value"}, configKeys, nil)
				return
			}
			for _, configKey := range configKeys {
				if configKey[0] == args[0] {
					slice := configKeys[0:1]
					out.Table([]string{"Key", "Value"}, slice, nil)
				}
			}
		},
	}
	return config
}

func createLine(key string, value string) []string {
	line := make([]string, 2)
	line[0] = key
	line[1] = value
	return line
}
