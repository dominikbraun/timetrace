// Package out provides functions for printing messages to the standard output.
package out

import (
	"fmt"
	"os"

	"github.com/enescakir/emoji"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

var (
	// used for the headers on the tables
	backgroundColor = []int{
		tablewriter.BgCyanColor,
		tablewriter.BgMagentaColor,
		tablewriter.BgGreenColor,
		tablewriter.BgRedColor,
		tablewriter.BgYellowColor,
	}
)

// Success prints a colored, formatted success message prefixed with an emoji.
func Success(format string, a ...interface{}) {
	p(color.FgGreen, emoji.CheckMark, format, a...)
}

// Info prints a colored, formatted info message prefixed with an emoji.
func Info(format string, a ...interface{}) {
	p(color.FgCyan, emoji.LightBulb, format, a...)
}

// Warn prints a colored, formatted warning message prefixed with an emoji.
func Warn(format string, a ...interface{}) {
	p(color.FgHiYellow, emoji.Warning, format, a...)
}

// Err prints a colored, formatted error message prefixed with an emoji.
func Err(format string, a ...interface{}) {
	p(color.FgHiRed, emoji.ExclamationMark, format, a...)
}

// Table renders a table with the given rows to the standard output.
func Table(header []string, rows [][]string, footer []string, opts ...TableOption) {
	paddedHeaders := headersWithPadding(header)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(paddedHeaders)
	setHeaderColor(table, paddedHeaders)
	// If footer array is not empty, then render footer in table.
	if len(footer) > 0 {
		paddedFooters := headersWithPadding(footer)
		table.SetFooter(paddedFooters)
		table.SetFooterAlignment(tablewriter.ALIGN_LEFT)
	}

	// if provided apply table options to table
	// var table must be a pointer else options wont be apply
	for _, opt := range opts {
		opt(table)
	}
	table.AppendBulk(rows)
	table.Render()
}

// setHeaderColor set colors for the headers on the table
func setHeaderColor(table *tablewriter.Table, header []string) {
	colors := []tablewriter.Colors{}
	for i := range header {
		color := tablewriter.Colors{tablewriter.Bold, backgroundColor[i%len(backgroundColor)]}
		colors = append(colors, color)
	}
	table.SetHeaderColor(colors...)
}

// headersWithPadding prepends and appends a space to each header
func headersWithPadding(headers []string) []string {
	newHeaders := make([]string, len(headers))
	for idx, val := range headers {
		// Do not pad empty headers.
		if val == "" {
			continue
		}
		newHeaders[idx] = " " + val + " "
	}
	return newHeaders
}

func p(attribute color.Attribute, emoji emoji.Emoji, format string, a ...interface{}) {
	formatWithEmoji := fmt.Sprintf("%v %s\n", emoji, format)
	_, _ = color.New(attribute).Printf(formatWithEmoji, a...)
}

// TableOptions allows to modify the table instance
// with different functionalities
type TableOption func(*tablewriter.Table)

// TableWithCellMerge apply tablewriter.SetAuthMergeCellsByColumnIndex to the
// table instance and enables tablewriter.SetRowLine.
// Allows to group rows by a column index
func TableWithCellMerge(mergeByIndex int) func(*tablewriter.Table) {
	return func(t *tablewriter.Table) {
		var index = mergeByIndex
		if mergeByIndex > t.NumLines() {
			index = 0
		}
		t.SetAutoMergeCellsByColumnIndex([]int{index})
		t.SetRowLine(true)
	}
}

// TableFooterColor adds colors to the tablewriter.Footer
func TableFooterColor(colors ...tablewriter.Colors) func(*tablewriter.Table) {
	return func(t *tablewriter.Table) {
		t.SetFooterColor(colors...)
	}
}
