package cli

import (
	"errors"
	"strconv"
	"strings"
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
	var options editOptions
	editProject := &cobra.Command{
		Use:   "project <KEY>",
		Short: "Edit a project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			if options.Revert {
				if err := t.RevertProject(key); err != nil {
					out.Err("Failed to revert project: %s", err.Error())
				} else {
					out.Info("Project backup restored successfuly")
				}
				return
			}

			if err := t.BackupProject(key); err != nil {
				out.Err("Failed to backup project before edit: %s", err.Error())
				return
			}
			out.Info("Opening %s in default editor", key)

			if err := t.EditProject(key); err != nil {
				out.Err("Failed to edit project: %s", err.Error())
				return
			}

			out.Success("Successfully edited %s", key)
		},
	}

	editProject.PersistentFlags().BoolVarP(&options.Revert, "revert", "r", false, "Restores the project to it's state prior to the last 'edit' command.")

	return editProject
}

type editOptions struct {
	Plus   string
	Minus  string
	Revert bool
}

func editRecordCommand(t *core.Timetrace) *cobra.Command {
	var options editOptions

	editRecord := &cobra.Command{
		Use:   "record {<KEY>|latest|@ID}",
		Short: "Edit a record",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if options.Plus != "" && options.Minus != "" {
				out.Err("Plus and minus flag can not be combined: %s", errors.New("edit not possible"))
				return
			}

			var recordTime time.Time
			var err error
			// if more aliases are needed, this should be expanded to a switch
			if strings.ToLower(args[0]) == "latest" {
				rec, err := t.LoadLatestRecord()
				if err != nil {
					out.Err("Error on loading last record: %s", err.Error())
					return
				}
				recordTime = rec.Start
			} else if strings.Contains(args[0], "@") {
				id, err := strconv.Atoi(args[0][1:])
				if err != nil {
					out.Err("Error on parsing ID: %s", err.Error())
					return
				}
				rec, err := t.LoadRecordByID(id)
				if err != nil {
					out.Err("Error on loading last record: %s", err.Error())
					return
				}
				if rec == nil {
					out.Err("No record of given ID started today.")
					return
				}
				recordTime = rec.Start
			} else {
				recordTime, err = t.Formatter().ParseRecordKey(args[0])
				if err != nil {
					out.Err("Failed to parse date argument: %s", err.Error())
					return
				}
			}

			if options.Revert {
				if err := t.RevertRecord(recordTime); err != nil {
					out.Err("Failed to revert record: %s", err.Error())
				} else {
					out.Info("Record backup restored successfully")
				}
				return
			}

			if err := t.BackupRecord(recordTime); err != nil {
				out.Err("Failed to backup record before edit: %s", err.Error())
				return
			}

			if options.Minus == "" && options.Plus == "" {
				out.Info("Opening %s in default editor", recordTime)
				if err := t.EditRecordManual(recordTime); err != nil {
					out.Err("Failed to edit record: %s", err.Error())
					return
				}
			} else {
				if err := t.EditRecord(recordTime, options.Plus, options.Minus); err != nil {
					out.Err("Failed to edit record: %s", err.Error())
					return
				}
			}

			out.Success("Successfully edited %s", recordTime)
		},
	}

	editRecord.PersistentFlags().StringVarP(&options.Plus, "plus", "p", "", "Adds the given duration to the end time of the record")
	editRecord.PersistentFlags().StringVarP(&options.Minus, "minus", "m", "", "Substracts the given duration to the end time of the record")
	editRecord.PersistentFlags().BoolVarP(&options.Revert, "revert", "r", false, "Restores the record to it's state prior to the last 'edit' command.")

	return editRecord
}
