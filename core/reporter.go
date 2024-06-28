package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	defaultTotalSymbol = "∑"
)

func FilterNoneNilEndTime(r *Record) bool {
	return r.End != nil
}

// FilterBillable returns a record filter for the billable flag
// that filters based on the passed display value
func FilterBillable(display bool) func(*Record) bool {
	return func(r *Record) bool {
		return r.IsBillable == display
	}
}

// FilterByProject returns true if the record matches the given project keys
// if a module is given "mod@project" filter checks if module and project match
func FilterByProject(key string) func(*Record) bool {
	module := strings.Split(key, "@") // if a module is given mod@project
	return func(r *Record) bool {
		recordModule := strings.Split(r.Project.Key, "@")
		// mod@key - mod@key -> search for project with module X
		if r.Project.IsModule() && len(module) > 1 {
			return recordModule[0] == module[0] && recordModule[1] == module[1]
		}
		// key - mod@key || key - key -> search for project where key
		if len(module) == 1 {
			if r.Project.IsModule() {
				return recordModule[1] == key
			}
			return r.Project.Key == key
		}
		return false
	}
}

// FilterByTimeRange allows to determine whether a given records is in-between a time-range.
// If "to" is nil the upper boundary is ignored and vice versa with "from". If both are nil returns true
// start and end time are both inclusive.
// Explanition for the `to.AddDate(0,0,1)`:
// the "to" input will be YYYY-MM-DD 00:00:00, hence the actual tracked records of that
// date will be ignored as they are all bigger since their hh:mm:ss will be grather then 00:00:00
// of the "to" time. Adding one day to the "to" time will include records tracked on that date thus
// will make the "to" time inclusive
func FilterByTimeRange(start, end time.Time) func(*Record) bool {

	return func(r *Record) bool {
		if start.IsZero() && end.IsZero() {
			return true
		}
		if end.IsZero() {
			return r.Start.Unix() >= start.Unix()
		}
		if start.IsZero() {
			// adding one day end the "end" date is required in or for
			// the end-time to be inclusive
			return r.Start.Unix() <= end.AddDate(0, 0, 1).Unix()
		}

		// adding one day end the "end" date is required in or for
		// the end-time to be inclusive
		return r.Start.Unix() >= start.Unix() && r.Start.Unix() <= end.AddDate(0, 0, 1).Unix()
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
// projects with module will be grouped by
func (r *Reporter) sortAndMerge(records []*Record) {
	for _, record := range records {
		key := record.Project.Key
		keyParts := strings.Split(key, "@")
		if len(keyParts) > 1 {
			key = keyParts[1]
		}
		cached, ok := r.report[key]
		if !ok {
			r.report[key] = []*Record{record}
			r.totals[key] = record.Duration()
			continue
		}
		r.report[key] = append(cached, record)
		// compute updated total
		tmp := r.totals[key] + record.Duration()
		r.totals[key] = tmp
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
			keyParts := strings.Split(record.Project.Key, "@")
			module, key := "", keyParts[0]
			if len(keyParts) > 1 {
				module = keyParts[0]
				key = keyParts[1]
			}
			billable := "no"
			if record.IsBillable {
				billable = "yes"
			}
			date := r.t.Formatter().PrettyDateString(record.Start)
			start := r.t.Formatter().TimeString(record.Start)
			end := r.t.Formatter().TimeString(*record.End)

			rows = append(rows, []string{key, module, date, start, end, billable, strings.Join(record.Notes, "\n"), ""})
		}
		// append with last row for total of tracked time for project
		rows = append(rows, []string{"", "", "", "", "", "", defaultTotalSymbol, r.t.Formatter().FormatDuration(r.totals[key])})
		totalSum += r.totals[key]
	}
	return rows, r.t.Formatter().FormatDuration(totalSum)
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
	b, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("could not marshal report to json")
	}
	return b, nil
}
