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
					out.Err("failed to revert project: %s", err.Error())
				} else {
					out.Info("project backup restored successfuly")
				}
				return
			}

			if err := t.BackupProject(key); err != nil {
				out.Err("failed to backup project before edit: %s", err.Error())
				return
			}
			out.Info("opening %s in default editor", key)

			if err := t.EditProject(key); err != nil {
				out.Err("failed to edit project: %s", err.Error())
				return
			}

			out.Success("successfully edited %s", key)
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
				out.Err("plus and minus flag can not be combined: %s", errors.New("edit not possible"))
				return
			}

			recordTime, err := getRecordTimeFromArg(t, args[0])

			if err != nil {
				out.Err(err.Error())
				return
			}

			if options.Revert {
				if err := t.RevertRecord(recordTime); err != nil {
					out.Err("failed to revert record: %s", err.Error())
				} else {
					out.Info("Record backup restored successfully")
				}
				return
			}

			if err := t.BackupRecord(recordTime); err != nil {
				out.Err("failed to backup record before edit: %s", err.Error())
				return
			}

			if options.Minus == "" && options.Plus == "" {
				out.Info("Opening %s in default editor", recordTime)
				if err := t.EditRecordManual(recordTime); err != nil {
					out.Err("failed to edit record: %s", err.Error())
					return
				}
			} else {
				if err := t.EditRecord(recordTime, options.Plus, options.Minus); err != nil {
					out.Err("failed to edit record: %s", err.Error())
					return
				}
			}

			out.Success("successfully edited %s", recordTime)
		},
	}

	editRecord.PersistentFlags().StringVarP(&options.Plus, "plus", "p", "", "Adds the given duration to the end time of the record")
	editRecord.PersistentFlags().StringVarP(&options.Minus, "minus", "m", "", "Substracts the given duration to the end time of the record")
	editRecord.PersistentFlags().BoolVarP(&options.Revert, "revert", "r", false, "Restores the record to it's state prior to the last 'edit' command.")

	return editRecord
}

func getRecordTimeFromArg(t *core.Timetrace, arg string) (time.Time, error) {
	var recordTime time.Time
	var err error
	// if more aliases are needed, this should be expanded to a switch
	if strings.ToLower(arg) == "latest" {
		rec, err := t.LoadLatestRecord()
		if err != nil {
			err = errors.New("error on loading last record: " + err.Error())
			return recordTime, err
		}
		recordTime = rec.Start
	} else if strings.Contains(arg, "@") {
		id, err := strconv.Atoi(arg[1:])
		if err != nil {
			err = errors.New("error on parsing ID: " + err.Error())
			return recordTime, err
		}
		rec, err := t.LoadRecordByID(id)
		if err != nil {
			err = errors.New("error on loading last record: " + err.Error())
			return recordTime, err
		}
		if rec == nil {
			err = errors.New("no record of given ID started today")
			return recordTime, err
		}
		recordTime = rec.Start
	} else {
		recordTime, err = t.Formatter().ParseRecordKey(arg)
		if err != nil {
			err = errors.New("failed to parse date argument: " + err.Error())
			return recordTime, err
		}
	}

	return recordTime, err
}
