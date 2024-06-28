package cli

import (
	"log"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

func attachCommand(t *core.Timetrace) *cobra.Command {
	attach := &cobra.Command{
		Use:   "attach",
		Short: "Attach a resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	attach.AddCommand(attachNoteCommand(t))

	return attach
}

func attachNoteCommand(t *core.Timetrace) *cobra.Command {
	var options attachOptions
	attachNote := &cobra.Command{
		Use:   "note",
		Short: "Attach a note",
		//		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			latestRecord, err := t.LoadLatestRecord()
			if err != nil {
				log.Fatal(err)
			}

			if options.InLine != "" {
				latestRecord.Notes = append(latestRecord.Notes, options.InLine)

				err = t.SaveRecord(*latestRecord, true)
				if err != nil {
					log.Fatal(err)
				}
				out.Success("successfully attached note %q to record", options.InLine)

			} else {

				out.Err("this action is not available yet, please use -m to add a note")

				// I think what will be good here is when the command is run, it allows
				// the user to view the record that the user wants to attach a note to
				// just like how `git commit -v` works.
			}
		},
	}

	attachNote.PersistentFlags().StringVarP(&options.InLine, "message", "m", "", "Directly attach a note to record")

	return attachNote
}

type attachOptions struct {
	InLine string
}
