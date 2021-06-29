package core

import (
	"testing"
	"time"
)

func newTestRecord(s int, e int) Record {
	start := time.Now().Add(time.Duration(s) * time.Minute)
	end := time.Now().Add(time.Duration(e) * time.Minute)
	return Record{Start: start, End: &end}
}

func TestCollides(t *testing.T) {
	savedRec := newTestRecord(-60, -1)
	allRecs := []*Record{&savedRec}

	// rec1 starts and end after savedRec
	rec1 := newTestRecord(-1, 0)

	if collides(rec1, allRecs) {
		t.Error("records should not collide")
	}

	// rec2 starts in savedRec, ends after
	rec2 := newTestRecord(-30, 1)

	if !collides(rec2, allRecs) {
		t.Error("records should collide")
	}

	// rec3 start before savedRec, ends inside
	rec3 := newTestRecord(-75, -30)

	if !collides(rec3, allRecs) {
		t.Error("records should collide")
	}

	// rec4 starts and ends before savedRec
	rec4 := newTestRecord(-75, -70)

	if collides(rec4, allRecs) {
		t.Error("records should not collide")
	}

	// rec5 starts and ends inside savedRec
	rec5 := newTestRecord(-40, -20)

	if !collides(rec5, allRecs) {
		t.Error("records should collide")
	}

	// rec6 starts before and ends after savedRec
	rec6 := newTestRecord(-70, 10)

	if !collides(rec6, allRecs) {
		t.Error("records should collide")
	}
}

func TestFormatDuration(t *testing.T) {

	tt := []struct {
		Duration time.Duration
		Expected string
	}{
		{
			Duration: time.Duration(12 * time.Second),
			Expected: "0h 0min",
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
			Expected: "0h 0min",
		},
	}

	formatter := Formatter{}

	for _, test := range tt {
		strFormat := formatter.FormatDuration(test.Duration)
		if strFormat != test.Expected {
			t.Fatalf("format error: %s != %s", strFormat, test.Expected)
		}
	}
}
