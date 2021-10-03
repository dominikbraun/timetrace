package cli

import (
	"fmt"
	"strings"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/out"

	"github.com/spf13/cobra"
)

const TagsPrefix = "+"

type startOptions struct {
	isBillable    bool
	isNonBillable bool // Used for overwriting `billable: true` in the project config.
}

func startCommand(t *core.Timetrace) *cobra.Command {
	var options startOptions

	start := &cobra.Command{
		Use:   "start <PROJECT KEY> [+TAG1, +TAG2, ...]",
		Short: "Start tracking time",
		Run: func(cmd *cobra.Command, args []string) {
			projectKey := args[0]
			tags := args[1:]

			// Limit number of tags to 3
			if len(tags) > 3 {
				out.Err("Failed to start tracking: At most 3 tags are allowed, got %v tags", len(tags))
				return
			}

			isBillable := options.isBillable

			// If there is a default configuration for the project key, use that configuration.
			if projectConfig, ok := t.Config().Projects[projectKey]; ok {
				isBillable = projectConfig.Billable
			}

			if options.isNonBillable {
				isBillable = false
			}

			tagNames, err := extractTagNames(tags)
			if err != nil {
				out.Err("failed to start tracking: %s", err.Error())
				return
			}

			if err := t.Start(projectKey, isBillable, tagNames); err != nil {
				out.Err("failed to start tracking: %s", err.Error())
				return
			}

			out.Success("Started tracking time")
		},
	}

	start.Flags().BoolVarP(&options.isBillable, "billable", "b",
		false, `mark tracked time as billable`)

	start.Flags().BoolVar(&options.isNonBillable, "non-billable",
		false, `mark tracked time as non-billable if the project is configured as billable`)

	return start
}

func extractTagNames(tagsWithPrefix []string) ([]string, error) {
	tagNames := make([]string, 0)

	for _, tagWithPrefix := range tagsWithPrefix {
		if !strings.HasPrefix(tagWithPrefix, TagsPrefix) {
			return nil, fmt.Errorf("'%s' is not a valid tag. Tags must start with %s", tagWithPrefix, TagsPrefix)
		}
		tagNames = append(tagNames, tagWithPrefix[1:])
	}

	return tagNames, nil
}
