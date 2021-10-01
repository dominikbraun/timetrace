package cli

import (
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func createCommand(t *core.Timetrace) *cobra.Command {
	create := &cobra.Command{
		Use:   "create",
		Short: "Create a new resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	create.AddCommand(createProjectCommand(t))
	create.AddCommand(createRecordCommand(t))

	return create
}

func createProjectCommand(t *core.Timetrace) *cobra.Command {
	createProject := &cobra.Command{
		Use:   "project <KEY>",
		Short: "Create a new project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]

			project := core.Project{
				Key: key,
			}

			if err := t.SaveProject(project, false); err != nil {
				out.Err("failed to create project: %s", err.Error())
				return
			}

			out.Success("Created project %s", key)
		},
	}

	return createProject
}

func createRecordCommand(t *core.Timetrace) *cobra.Command {
	var options startOptions
	createRecord := &cobra.Command{
		Use:   "record <PROJECT KEY> {<YYYY-MM-DD>|today|yesterday} <HH:MM> <HH:MM>",
		Short: "Create a new record",
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			project, err := t.LoadProject(key)
			if err != nil {
				out.Err("failed to get project: %s", key)
				return
			}

			date, err := t.Formatter().ParseDate(args[1])
			if err != nil {
				out.Err("failed to parse date: %s", err.Error())
				return
			}

			start, err := t.Formatter().ParseTime(args[2])
			if err != nil {
				out.Err("failed to parse start time: %s", err.Error())
				return
			}
			start = t.Formatter().CombineDateAndTime(date, start)

			end, err := t.Formatter().ParseTime(args[3])
			if err != nil {
				out.Err("failed to parse end time: %s", err.Error())
				return
			}
			end = t.Formatter().CombineDateAndTime(date, end)

			if end.Before(start) {
				out.Err("end time is before start time")
				return
			}

			now := time.Now()
			if now.Before(start) || now.Before(end) {
				out.Err("provided record happens in the future")
				return
			}

			record := core.Record{
				Project:    project,
				Start:      start,
				End:        &end,
				IsBillable: options.isBillable,
			}

			collides, err := t.RecordCollides(record)
			if err != nil {
				out.Err("error on check if record collides: %s", err.Error())
				return
			}
			if collides {
				return
			}

			if err := t.SaveRecord(record, false); err != nil {
				out.Err("failed to create record: %s", err.Error())
				return
			}

			out.Success("created record %s in project %s", t.Formatter().TimeString(record.Start), key)
		},
	}

	createRecord.Flags().BoolVarP(&options.isBillable, "billable", "b",
		false, `mark tracked time as billable`)

	return createRecord
}
