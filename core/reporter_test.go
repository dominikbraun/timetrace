package core

import (
	"testing"
	"time"
)

func TestReportFilterTimeRange(t *testing.T) {
	f := &Formatter{
		use12Hours: true,
	}
	tt := []struct {
		From     time.Time
		To       time.Time
		Record   Record
		Expected bool
	}{}

	d1, _ := f.ParseDate("2021-04-01")
	d2, _ := f.ParseDate("2021-06-01")
	r1 := d1.AddDate(0, 1, 0) // record.Start will be 2021-05-01
	// both times given and in range -> must be true
	tt = append(tt, struct {
		From     time.Time
		To       time.Time
		Record   Record
		Expected bool
	}{
		From: d1,
		To:   d2,
		Record: Record{
			Start: r1,
		},
		Expected: true,
	})
	// to is nil -> must be true
	tt = append(tt, struct {
		From     time.Time
		To       time.Time
		Record   Record
		Expected bool
	}{
		From: d1,
		To:   time.Time{},
		Record: Record{
			Start: r1,
		},
		Expected: true,
	})
	// both from and to are nil -> must be true
	tt = append(tt, struct {
		From     time.Time
		To       time.Time
		Record   Record
		Expected bool
	}{
		From: time.Time{},
		To:   time.Time{},
		Record: Record{
			Start: r1,
		},
		Expected: true,
	})

	d3, _ := f.ParseDate("2021-04-01")
	d4, _ := f.ParseDate("2021-06-01")
	r2 := d1.AddDate(0, 3, 0) // record.Start will be 2021-07-01
	// both times given but r.Start is not in range -> must be false
	tt = append(tt, struct {
		From     time.Time
		To       time.Time
		Record   Record
		Expected bool
	}{
		From: d3,
		To:   d4,
		Record: Record{
			Start: r2,
		},
		Expected: false,
	})
	for _, tc := range tt {
		ok := FilterByTimeRange(tc.From, tc.To)(&tc.Record)
		if ok != tc.Expected {
			t.Fatalf("filter time range: filter returned wrong statement")
		}
	}
}

func TestProjectFilter(t *testing.T) {
	tt := []struct {
		Key      string
		R        Record
		Expected bool
	}{
		{
			Key: "test",
			R: Record{
				Project: &Project{
					Key: "test",
				},
			},
			Expected: true,
		},
		{
			Key: "mod@test",
			R: Record{
				Project: &Project{
					Key: "mod@test",
				},
			},
			Expected: true,
		},
		{
			Key: "test",
			R: Record{
				Project: &Project{
					Key: "mod@test",
				},
			},
			Expected: true,
		},
		{
			Key: "mod@test",
			R: Record{
				Project: &Project{
					Key: "test",
				},
			},
			Expected: false,
		},
		{
			Key: "mod@test",
			R: Record{
				Project: &Project{
					Key: "not-test",
				},
			},
			Expected: false,
		},
		{
			Key: "test",
			R: Record{
				Project: &Project{
					Key: "not-test",
				},
			},
			Expected: false,
		},
	}

	for _, tc := range tt {
		filter := FilterByProject(tc.Key)
		filtered := filter(&tc.R)
		if filtered != tc.Expected {
			t.Fatalf("project-filter failed: want %v, have: %v for key: %s and project.key: %s", tc.Expected, filtered, tc.Key, tc.R.Project.Key)
		}
	}
}
