package core

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	defaultTotalSymbol = "∑"
)

func FilterNonNilEndTime(r *Record) bool {
	return r.End != nil
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
// If "to" is nil the upper boundary is ignored and vice versa with "from". If both are nil returns true
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

// Reporter holds map of projects with slice of project records
type Reporter struct {
	t *Timetrace
	// report stores project-key:tracked-records
	report map[string][]*Record
	// total stores the overall time spend on a project
	totals map[string]time.Duration
}

// sortAndMerge assigns each record in the given slice to the correct project key in the
// Reporter.report map and computes each projects total time.Duration
func (r *Reporter) sortAndMerge(records []*Record) {
	for _, record := range records {
		cached, ok := r.report[record.Project.Key]
		if !ok {
			r.report[record.Project.Key] = []*Record{record}
			r.totals[record.Project.Key] = record.End.Sub(record.Start)
			continue
		}
		r.report[record.Project.Key] = append(cached, record)
		// compute updated total
		tmp := r.totals[record.Project.Key] + record.End.Sub(record.Start)
		r.totals[record.Project.Key] = tmp
	}
}

// Table prepares the r.report and r.totals data in a way that it can be consumed by the out.Table
// It returns a [][]string where each []string represents one record of a project and
// the total sum of time for all projects
func (r Reporter) Table() ([][]string, string) {
	var rows = make([][]string, 0)
	var totalSum time.Duration

	for key, records := range r.report {
		for _, record := range records {
			project := key
			billable := "no"
			if record.IsBillable {
				billable = "yes"
			}
			date := r.t.Formatter().PrettyDateString(record.Start)
			start := r.t.Formatter().TimeString(record.Start)
			end := r.t.Formatter().TimeString(*record.End)

			rows = append(rows, []string{project, date, start, end, billable, ""})
		}
		// append with last row for total of tracked time for project
		rows = append(rows, []string{key, "", "", "", defaultTotalSymbol, formatDuration(r.totals[key])})
		totalSum += r.totals[key]
	}
	return rows, formatDuration(totalSum)
}

// Json prepares the r.report and r.totals data so that it can be written to a json file
func (r Reporter) Json() ([]byte, error) {
	var result = make(map[string]interface{})

	for key, records := range r.report {
		var total time.Duration
		if t, ok := r.totals[key]; ok {
			total = t
		}
		result[key] = map[string]interface{}{
			"records": records,
			"total":   total,
		}
	}
	b, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("could not marshal report to json")
	}
	return b, nil
}
