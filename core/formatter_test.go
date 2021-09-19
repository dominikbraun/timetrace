package core

import "testing"

func TestFormatter_FormatTags(t *testing.T) {
	tests := map[string]struct {
		tags   []string
		output string
	}{
		"one tag": {
			tags:   []string{"coffee"},
			output: "coffee",
		},
		"two tags": {
			tags:   []string{"coffee", "espresso"},
			output: "coffee, espresso",
		},
		"four tags": {
			tags:   []string{"coffee", "espresso", "morning"},
			output: "coffee, espresso, morning",
		},
	}

	formatter := Formatter{}

	for name, tc := range tests {
		formattedTags := formatter.FormatTags(tc.tags)
		if formattedTags != tc.output {
			t.Errorf("%s: expected %s, got %s", name, tc.output, formattedTags)
		}
	}
}
