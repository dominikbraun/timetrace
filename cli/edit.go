package cli

import (
	"errors"
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func editCommand(t *core.Timetrace) *cobra.Command {
	edit := &cobra.Command{
		Use:   "edit",
		Short: "Edit a resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	edit.AddCommand(editProjectCommand(t))
	edit.AddCommand(editRecordCommand(t))

	return edit
}

func editProjectCommand(t *core.Timetrace) *cobra.Command {
	editProject := &cobra.Command{
		Use:   "project <KEY>",
		Short: "Edit a project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			out.Info("Opening %s in default editor", key)

			if err := t.EditProject(key); err != nil {
				out.Err("Failed to edit project: %s", err.Error())
				return
			}

			out.Success("Successfully edited %s", key)
		},
	}

	return editProject
}

type editOptions struct {
	Plus  string
	Minus string
}

func editRecordCommand(t *core.Timetrace) *cobra.Command {
	var options editOptions

	editRecord := &cobra.Command{
		Use:   "record <KEY>",
		Short: "Edit a record",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if options.Plus != "" && options.Minus != "" {
				out.Err("Plus and minus flag can not be combined: %s", errors.New("edit not possible"))
				return
			}

			layout := defaultRecordArgLayout

			if t.Config().Use12Hours {
				layout = "2006-01-02-03-04PM"
			}

			recordTime, err := time.Parse(layout, args[0])
			if err != nil {
				out.Err("Failed to parse date argument: %s", err.Error())
				return
			}

			if options.Minus == "" && options.Plus == "" {
				out.Info("Opening %s in default editor", recordTime)
				if err := t.EditRecordManual(recordTime); err != nil {
					out.Err("Failed to edit project: %s", err.Error())
					return
				}
			} else {
				if err := t.EditRecord(recordTime, options.Plus, options.Minus); err != nil {
					out.Err("Failed to edit project: %s", err.Error())
					return
				}
			}

			out.Success("Successfully edited %s", recordTime)
		},
	}

	editRecord.PersistentFlags().StringVarP(&options.Plus, "plus", "p", "", "Adds the given duration to the end time of the record")
	editRecord.PersistentFlags().StringVarP(&options.Minus, "minus", "m", "", "Substracts the given duration to the end time of the record")

	return editRecord
}
