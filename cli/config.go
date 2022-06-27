package cli

import (
	"reflect"
	"strconv"
	"strings"

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
			header := []string{"Key", "Value"}
			valueOfConfig := reflect.ValueOf(t.Config())
			typeOfConfig := reflect.Indirect(valueOfConfig).Type()
			values := make([][]string, reflect.Indirect(valueOfConfig).NumField())
			for i := 0; i < reflect.Indirect(valueOfConfig).NumField(); i++ {
				switch reflect.Indirect(valueOfConfig).Field(i).Interface().(type) {
				case bool:
					values[i] = []string{typeOfConfig.Field(i).Name, strconv.FormatBool(reflect.Indirect(valueOfConfig).Field(i).Bool())}
				default:
					values[i] = []string{typeOfConfig.Field(i).Name, reflect.Indirect(valueOfConfig).Field(i).String()}
				}
			}
			if len(args) == 0 {
				out.Table(header, values, nil)
				return
			}
			for _, configSetting := range values {

				if strings.ToUpper(configSetting[0]) == strings.ToUpper(args[0]) {
					out.Table(header, [][]string{configSetting}, nil)
					return
				}
			}

		},
	}
	return config
}
