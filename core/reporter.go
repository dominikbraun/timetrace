package core

import (
	"fmt"
	"time"
)

// Reporter holds map of projects with slice of project records
type Reporter struct {
	// report stores project-key:tracked-records
	report map[string][]*Record
	// total stores the overall time spend on a project
	totals map[string]time.Duration
}

// sortAndMerge assigns each record in the given slice to the correct project key in the
// Reporter.report map
func (r *Reporter) sortAndMerge(reocrds []*Record) {
	for _, record := range reocrds {
		cached, ok := r.report[record.Project.Key]
		if !ok {
			r.report[record.Project.Key] = []*Record{record}
			r.totals[record.Project.Key] = record.End.Sub(record.Start)
		}
		r.report[record.Project.Key] = append(cached, record)
		// compute updated total
		tmp := r.totals[record.Project.Key] + record.End.Sub(record.Start)
		r.totals[record.Project.Key] = tmp
	}
	fmt.Println(r.totals)
}

// FilterBillable returns true if a records is listed as billable
func FilterBillable(r *Record) bool {
	return r.IsBillable
}

// FilterByProject returns true if the record matches the given project keys
func FilterByProject(key string) func(*Record) bool {
	return func(r *Record) bool {
		return r.Project.Key == key
	}
}

// FilterByTimeRange allows to determine whether a given records is in-between a time-range.
// If "to" is nil upper boundary is ignored and vice versa with "from". If both are nil returns true
func FilterByTimeRange(from, to *time.Time) func(*Record) bool {
	return func(r *Record) bool {
		if from == nil && to == nil {
			return true
		}
		if to == nil {
			return r.Start.Unix() >= from.Unix()
		}
		if from == nil {
			return r.Start.Unix() <= to.Unix()
		}
		return r.Start.Unix() >= from.Unix() && r.Start.Unix() <= to.Unix()
	}
}
