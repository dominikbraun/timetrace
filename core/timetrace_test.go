package core

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {

	tt := []struct {
		Duration time.Duration
		Expected string
	}{
		{
			Duration: time.Duration(12 * time.Second),
			Expected: "0h 0min 12sec",
		},
		{
			Duration: time.Duration(60 * time.Minute),
			Expected: "1h 0min",
		},
		{
			Duration: time.Duration(24 * time.Minute),
			Expected: "0h 24min",
		},
		{
			Duration: time.Duration((60*8 + 24) * time.Minute),
			Expected: "8h 24min",
		},
		{
			Duration: time.Duration((60*8+24)*time.Minute + 12*time.Second),
			Expected: "8h 24min",
		},
		{
			Duration: time.Duration(0 * time.Second),
			Expected: "0h 0min 0sec",
		},
	}

	for _, test := range tt {
		strFormat := formatDuration(test.Duration)
		if strFormat != test.Expected {
			t.Fatalf("format error: %s != %s", strFormat, test.Expected)
		}
	}
}
