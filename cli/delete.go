package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

var confirmed bool

func deleteCommand(t *core.Timetrace) *cobra.Command {
	delete := &cobra.Command{
		Use:   "delete",
		Short: "Delete a resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	delete.AddCommand(deleteProjectCommand(t))
	delete.AddCommand(deleteRecordCommand(t))
	delete.PersistentFlags().BoolVar(&confirmed, "yes", false, "Do not ask for confirmation")

	return delete
}

func deleteProjectCommand(t *core.Timetrace) *cobra.Command {
	deleteProject := &cobra.Command{
		Use:   "project <KEY>",
		Short: "Delete a project",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]

			project := core.Project{
				Key: key,
			}

			if err := t.DeleteProject(project); err != nil {
				out.Err("Failed to delete %s", err.Error())
				return
			}

			out.Success("Deleted project %s", key)
		},
	}

	return deleteProject
}

func deleteRecordCommand(t *core.Timetrace) *cobra.Command {

	// Depending on the use12hours setting, the command syntax either is
	// `record YYYY-MM-DD-HH-MM` or `record YYYY-MM-DD-HH-MMPM`.
	use := fmt.Sprintf("record %s", t.Formatter().RecordKeyLayout())

	deleteRecord := &cobra.Command{
		Use:   use,
		Short: "Delete a record",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			start, err := t.Formatter().ParseRecordKeyString(args[0])
			if err != nil {
				out.Err("Failed to parse date argument: %s", err.Error())
				return
			}

			record, err := t.LoadRecord(start)
			if err != nil {
				out.Err("Failed to read record: %s", err.Error())
				return
			}

			showRecord(record, t.Formatter())
			if !confirmed {
				if !askForConfirmation() {
					out.Info("Record NOT deleted.")
					return
				}
			}

			if err := t.DeleteRecord(*record); err != nil {
				out.Err("Failed to delete %s", err.Error())
				return
			}

			out.Success("Deleted record %s", args[0])
		},
	}

	return deleteRecord
}

func askForConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, "Please confirm (Y/N): ")
		s, _ := reader.ReadString('\n')
		s = strings.TrimSuffix(s, "\n")
		s = strings.ToLower(s)
		if len(s) > 1 {
			continue
		}
		if strings.Compare(s, "n") == 0 {
			return false
		} else if strings.Compare(s, "y") == 0 {
			break
		}
	}
	return true
}
