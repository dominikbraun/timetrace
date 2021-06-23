package cli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"
	"github.com/enescakir/emoji"
	"github.com/spf13/cobra"
)

func pushCommand(t *core.Timetrace) *cobra.Command {
	push := &cobra.Command{
		Use:   "push <INTEGRATION>",
		Short: "Pushes all unlogged records to the provided 3rd party provider (e.g. JIRA)",
		Args:  cobra.ExactArgs(1),
		// Use a prerun to clarify that the user is satisfied with the proposed
		// state change in JIRA
		PreRunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			integrationName := args[0]
			if _, exists := t.ListIntegrations()[integrationName]; !exists {
				err := fmt.Errorf("integration %s is not available", integrationName)
				out.Err(err.Error())
			}

			recordsToPush, err := t.VerifyPush(integrationName)
			if err != nil {
				out.Err("Aborting push. err: %s", err)
				return
			}

			// print what is about to be pushed as a table
			printPushTable(integrationName, recordsToPush)

			if !awaitConfirmation() {
				out.Info("not pushing records. Goodbye!")
				return
			}

			tw := out.NewTableWriter(len(recordsToPush))
			defer tw.Finish()
			if err := t.Push(integrationName, recordsToPush, tw); err != nil {
				out.Err(err.Error())
			}
		},
	}

	return push
}

func awaitConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		out.Info("[y/n]: ")

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func printPushTable(integrationName string, records []*core.Record) {
	headings := []string{
		"STATUS",
		"OPERATION",
		"DURATION",
		fmt.Sprintf("%s RECORD", integrationName),
	}

	out.Info("About to push to %s", integrationName)
	tableRowCols := make([][]string, len(records))
	for i, record := range records {
		duration := record.End.Sub(record.Start).Round(time.Second)
		tableRowCols[i] = []string{emoji.OutboxTray.String(), "PUSH", duration.String(), record.Project.Key}
	}
	out.Table(headings, tableRowCols, nil)
}
