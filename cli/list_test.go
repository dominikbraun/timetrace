package cli

import (
	"reflect"
	"testing"

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
