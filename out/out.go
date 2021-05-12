// Package out provides functions for printing messages to the standard output.
package out

import (
	"fmt"
	"os"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// Success prints a colored, formatted success message prefixed with an emoji.
func Success(format string, a ...interface{}) {
	p(color.FgGreen, emoji.CheckMark, format, a...)
}

// Info prints a colored, formatted info message prefixed with an emoji.
func Info(format string, a ...interface{}) {
	p(color.FgCyan, emoji.LightBulb, format, a...)
}

// Err prints a colored, formatted error message prefixed with an emoji.
func Err(format string, a ...interface{}) {
	p(color.FgHiRed, emoji.ExclamationMark, format, a...)
}

// Table renders a table with the given rows to the standard output.
func Table(header []string, rows [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(rows)
	table.Render()
}

func p(attribute color.Attribute, emoji emoji.Emoji, format string, a ...interface{}) {
	formatWithEmoji := fmt.Sprintf("%v %s\n", emoji, format)
	_, _ = color.New(attribute).Printf(formatWithEmoji, a...)
}
