package cli

import (
	"reflect"
	"testing"
	"time"

	"github.com/dominikbraun/timetrace/core"
)

func TestFilterBillableRecords(t *testing.T) {

	tt := []struct {
		title    string
		records  []*core.Record
		expected []*core.Record
	}{
		{
			title: "all records are billable",
			records: []*core.Record{
				{Project: &core.Project{Key: "a"}, IsBillable: true},
				{Project: &core.Project{Key: "b"}, IsBillable: true},
			},
			expected: []*core.Record{
				{Project: &core.Project{Key: "a"}, IsBillable: true},
				{Project: &core.Project{Key: "b"}, IsBillable: true},
			},
		},
		{
			title: "no records are billable",
			records: []*core.Record{
				{Project: &core.Project{Key: "a"}, IsBillable: false},
				{Project: &core.Project{Key: "b"}, IsBillable: false},
			},
			expected: []*core.Record{},
		},
		{
			title: "half of records are billable",
			records: []*core.Record{
				{Project: &core.Project{Key: "a"}, IsBillable: true},
				{Project: &core.Project{Key: "b"}, IsBillable: true},
				{Project: &core.Project{Key: "c"}, IsBillable: false},
				{Project: &core.Project{Key: "d"}, IsBillable: false},
			},
			expected: []*core.Record{
				{Project: &core.Project{Key: "a"}, IsBillable: true},
				{Project: &core.Project{Key: "b"}, IsBillable: true},
			},
		},
	}

	for _, test := range tt {
		billableRecords := filterBillableRecords(test.records)
		if !reflect.DeepEqual(billableRecords, test.expected) {
			t.Fatalf("error when %s: %v != %v", test.title, billableRecords, test.expected)
		}
	}
}

func TestFilterProjectRecords(t *testing.T) {

	tt := []struct {
		title    string
		filter   string
		records  []*core.Record
		expected []*core.Record
	}{
		{
			title:  "filter by project a",
			filter: "a",
			records: []*core.Record{
				{Project: &core.Project{Key: "a"}, IsBillable: false},
				{Project: &core.Project{Key: "b"}, IsBillable: true},
			},
			expected: []*core.Record{
				{Project: &core.Project{Key: "a"}, IsBillable: false},
			},
		},
		{
			title:  "filter by project b",
			filter: "b",
			records: []*core.Record{
				{Project: &core.Project{Key: "a"}, IsBillable: false},
				{Project: &core.Project{Key: "b"}, IsBillable: true},
			},
			expected: []*core.Record{
				{Project: &core.Project{Key: "b"}, IsBillable: true},
			},
		},
		{
			title:  "filter by project module@b",
			filter: "b",
			records: []*core.Record{
				{Project: &core.Project{Key: "module@b"}, IsBillable: false},
			},
			expected: []*core.Record{
				{Project: &core.Project{Key: "module@b"}, IsBillable: false},
			},
		},
		{
			title:  "no records found",
			filter: "a",
			records: []*core.Record{
				{Project: &core.Project{Key: "c"}, IsBillable: false},
				{Project: &core.Project{Key: "d"}, IsBillable: false},
			},
			expected: []*core.Record{},
		},
	}

	for _, test := range tt {
		projectRecords := filterProjectRecords(test.records, test.filter)
		if !reflect.DeepEqual(projectRecords, test.expected) {
			t.Fatalf("error when %s: %v != %v", test.title, projectRecords, test.expected)
		}
	}
}

func TestTotalTrackedTime(t *testing.T) {
	tt := []struct {
		records  []*core.Record
		expected time.Duration
	}{
		{records: []*core.Record{
			{
				Start: time.Date(2021, 06, 07, 16, 00, 00, 00, time.Local),          // 4:00PM
				End:   timePtr(time.Date(2021, 06, 07, 16, 25, 00, 00, time.Local)), // 4:25PM
			},
			{
				Start: time.Date(2021, 06, 07, 16, 30, 00, 00, time.Local),          // 4:30PM
				End:   timePtr(time.Date(2021, 06, 07, 16, 50, 00, 00, time.Local)), // 4:50PM
			},
			{
				Start: time.Date(2021, 06, 07, 16, 55, 00, 00, time.Local),          // 4:55PM
				End:   timePtr(time.Date(2021, 06, 07, 17, 10, 00, 00, time.Local)), // 5:10PM
			},
		},
			expected: time.Duration(time.Hour),
		},
	}
	for _, test := range tt {
		totalTime := getTotalTrackedTime(test.records)
		if totalTime != test.expected {
			t.Fatalf("error when %v != %v", totalTime, test.expected)
		}
	}
}

// timePtr gives a pointer of `time.Time` (an alias of int64).
func timePtr(t time.Time) *time.Time {
	return &t
}
