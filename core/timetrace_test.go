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

func newTestRecTracked(s int) Record {
	start := time.Now().Add(time.Duration(s) * time.Minute)
	return Record{Start: start}
}

func checkConsistent(t *testing.T, expect, result []*Record) {
	sameLen := len(result) == len(expect)
	sameContent := true

	if sameLen {
		for i := range result {
			if expect[i] != result[i] {
				sameContent = false
			}
		}
	}

	if !(sameLen && sameContent) {
		t.Errorf("should collide with :\n")
		for _, r := range expect {
			t.Errorf("%v\n", r)
		}
		t.Errorf("while collides return :\n")
		for _, r := range result {
			t.Errorf("%v\n", r)
		}
	}

}

func TestCollides(t *testing.T) {
	savedRec := newTestRecord(-60, -1)
	allRecs := []*Record{&savedRec}
	savedRecTracked := newTestRecTracked(-60)
	allRecsTracked := []*Record{&savedRecTracked}

	// rec1 starts and end after savedRec
	rec1 := newTestRecord(-1, 0)

	if collide, collidingRecs := collides(rec1, allRecs); collide && len(collidingRecs) == 0 {
		t.Error("records should not collide")
	}

	// rec2 starts in savedRec, ends after
	rec2 := newTestRecord(-30, 1)

	if collide, collidingRecs := collides(rec2, allRecs); !collide {
		t.Error("records should collide")
	} else {
		checkConsistent(t, allRecs, collidingRecs)
	}

	// rec3 start before savedRec, ends inside
	rec3 := newTestRecord(-75, -30)

	if collide, collidingRecs := collides(rec3, allRecs); !collide {
		t.Error("records should collide")
	} else {
		checkConsistent(t, allRecs, collidingRecs)
	}

	// rec4 starts and ends before savedRec
	rec4 := newTestRecord(-75, -70)

	if collide, collidingRecs := collides(rec4, allRecs); collide && len(collidingRecs) == 0 {
		t.Error("records should not collide")
	}

	// rec5 starts and ends inside savedRec
	rec5 := newTestRecord(-40, -20)

	if collide, collidingRecs := collides(rec5, allRecs); !collide {
		t.Error("records should collide")
	} else {
		checkConsistent(t, allRecs, collidingRecs)
	}

	// rec6 starts before and ends after savedRec
	rec6 := newTestRecord(-70, 10)

	if collide, collidingRecs := collides(rec6, allRecs); !collide {
		t.Error("records should collide")
	} else {
		checkConsistent(t, allRecs, collidingRecs)
	}

	// rec7 starts and ends at the same time as savedRec
	rec7 := newTestRecord(-60, -1)

	if collide, collidingRecs := collides(rec7, allRecs); !collide {
		t.Error("records should collide")
	} else {
		checkConsistent(t, allRecs, collidingRecs)
	}

	// rec7 starts at the same time as savedRecTracked
	rec8 := newTestRecord(-60, -1)

	if collide, collidingRecs := collides(rec8, allRecsTracked); !collide {
		t.Error("records should collide")
	} else {
		checkConsistent(t, allRecsTracked, collidingRecs)
	}

	// rec9 ends at the time savedRecTracked ends
	rec9 := newTestRecord(-80, -60)

	if collide, collidingRecs := collides(rec9, allRecsTracked); !collide {
		t.Error("records should collide")
	} else {
		checkConsistent(t, allRecsTracked, collidingRecs)
	}

	// rec10 ends after savedRecTracked starts
	rec10 := newTestRecord(-80, -50)

	if collide, collidingRecs := collides(rec10, allRecsTracked); !collide {
		t.Error("records should collide")
	} else {
		checkConsistent(t, allRecsTracked, collidingRecs)
	}

	// rec11 ends before savedRecTracked starts
	rec11 := newTestRecord(-80, -70)

	if collide, collidingRecs := collides(rec11, allRecsTracked); collide && len(collidingRecs) == 0 {
		t.Error("records should not collide")
	}
}

func TestFormatDuration(t *testing.T) {

	tt := []struct {
		Duration     time.Duration
		Expected     string
		ExpectedDec  string
		ExpectedBoth string
	}{
		{
			Duration:     time.Duration(12 * time.Second),
			Expected:     "0h 0min",
			ExpectedDec:  "0.0h",
			ExpectedBoth: "0h 0min 0.0h",
		},
		{
			Duration:     time.Duration(60 * time.Minute),
			Expected:     "1h 0min",
			ExpectedDec:  "1.0h",
			ExpectedBoth: "1h 0min 1.0h",
		},
		{
			Duration:     time.Duration(24 * time.Minute),
			Expected:     "0h 24min",
			ExpectedDec:  "0.4h",
			ExpectedBoth: "0h 24min 0.4h",
		},
		{
			Duration:     time.Duration((60*8 + 24) * time.Minute),
			Expected:     "8h 24min",
			ExpectedDec:  "8.4h",
			ExpectedBoth: "8h 24min 8.4h",
		},
		{
			Duration:     time.Duration((60*8+24)*time.Minute + 12*time.Second),
			Expected:     "8h 24min",
			ExpectedDec:  "8.4h",
			ExpectedBoth: "8h 24min 8.4h",
		},
		{
			Duration:     time.Duration(0 * time.Second),
			Expected:     "0h 0min",
			ExpectedDec:  "0.0h",
			ExpectedBoth: "0h 0min 0.0h",
		},
	}

	formatter := Formatter{}

	//Default Case
	for _, test := range tt {
		strFormat := formatter.FormatDuration(test.Duration)
		if strFormat != test.Expected {
			t.Fatalf("format error: %s != %s", strFormat, test.Expected)
		}
	}
	//Decimal Hours true
	formatter.useDecimalHours = "On"
	for _, test := range tt {
		strFormat := formatter.FormatDuration(test.Duration)
		if strFormat != test.ExpectedDec {
			t.Fatalf("format error: %s != %s", strFormat, test.ExpectedDec)
		}
	}
	//Decimal Hours both
	formatter.useDecimalHours = "Both"
	for _, test := range tt {
		strFormat := formatter.FormatDuration(test.Duration)
		if strFormat != test.ExpectedBoth {
			t.Fatalf("format error: %s != %s", strFormat, test.ExpectedBoth)
		}
	}
}
