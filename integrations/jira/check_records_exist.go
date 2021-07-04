package jira

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dominikbraun/timetrace/core"
)

// CheckRecordsExist queries JIRA for all current record states returning
// records that do not currently exist in JIRA. This places JIRA as the source
// of truth on whether records have been uploaded to JIRA already. It would be
// possible to save them to disk, but this runs the risk of reuploading a lot
// of work records to JIRA. Best to ask JIRA what it thinks the current state
// of its records is.
func (r *Repository) CheckRecordsExist(records []*core.Record) ([]*core.Record, error) {
	recordsByIssueID := buildIssueRecordMap(records)
	recordsByIssueIDToKeep := buildIssueRecordMap(records)

	for issueID, records := range recordsByIssueID {
		worklogs, err := r.getIssueWorklogs(issueID)
		if errors.Is(err, errIssueNotFound) {
			delete(recordsByIssueIDToKeep, issueID)
			continue
		}
		if err != nil {
			return nil, err
		}

		for _, wl := range worklogs {
			if wl.timetraceID == nil {
				continue
			}

			var reducedRecords []*core.Record
			for _, record := range records {
				if r.formatter.RecordKey(record) != *wl.timetraceID {
					reducedRecords = append(reducedRecords, record)
				}
			}
			recordsByIssueIDToKeep[issueID] = reducedRecords
		}
	}

	var recordsToKeep []*core.Record
	for _, rcs := range recordsByIssueIDToKeep {
		recordsToKeep = append(recordsToKeep, rcs...)
	}
	return recordsToKeep, nil
}

func removeRecord(a []*core.Record, i int) []*core.Record {
	// Remove the element at index i from a.
	a[i] = a[len(a)-1] // Copy last element to index i.
	a[len(a)-1] = nil  // Erase last element (write zero value).
	a = a[:len(a)-1]   // Truncate slice.
	return a
}

func buildIssueRecordMap(records []*core.Record) map[string][]*core.Record {
	m := make(map[string][]*core.Record)
	for _, record := range records {
		issueID := extractIssueID(record.Project)
		m[issueID] = append(m[issueID], record)
	}

	return m
}

func extractIssueID(project *core.Project) string {
	if !project.IsModule() {
		return project.Key
	}

	tokens := strings.Split(project.Key, "@")
	return tokens[0]
}

// timetraceProperty is the custom metadata we set on the records
// https://developer.atlassian.com/cloud/jira/platform/jira-entity-properties/#example-2--retrieving-data
type jiraEntityProperty struct {
	Value timetraceProperty `json:"value"`
}

type timetraceProperty struct {
	TimetraceRecordID string `json:"timetrace_record_id"`
}

func (r *Repository) getWorklogTimeTraceRecord(issueName string, worklogID string) (*string, error) {
	url := fmt.Sprintf("https://%s/%s/issue/%s/worklog/%s/properties/timetrace_id",
		r.jiraAddress,
		jiraV3BaseURL,
		issueName,
		worklogID,
	)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(r.email, r.authToken)

	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}

	// if 404, it's not an error, just that the worklog has no timetrace records
	if response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var worklogProperties jiraEntityProperty
	if err := json.Unmarshal(body, &worklogProperties); err != nil {
		return nil, err
	}

	return &worklogProperties.Value.TimetraceRecordID, nil
}

// issueNotFound is returned to signify all records associated with this
// project are not in jira. it can be quietly ignored
var errIssueNotFound = errors.New("issue not found")

type jiraWorklog struct {
	worklogID   string
	timetraceID *string
}

// getIssueWorklogs calls this endpoint
// https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-worklogs/#api-rest-api-3-issue-issueidorkey-worklog-get
// and returns a list of worklog IDs, with their timetrace records not nil if
// there are any present on the worklogs
func (r *Repository) getIssueWorklogs(issueName string) ([]jiraWorklog, error) {
	worklogIDs, err := r.getIssueWorklogIDs(issueName)
	if err != nil {
		return nil, err
	}

	worklogs := make([]jiraWorklog, len(worklogIDs))
	for i, worklogID := range worklogIDs {
		ttRecord, err := r.getWorklogTimeTraceRecord(issueName, worklogID)
		if err != nil {
			return nil, err
		}

		worklogs[i] = jiraWorklog{
			worklogID:   worklogID,
			timetraceID: ttRecord,
		}
	}

	return worklogs, nil
}

func (r *Repository) getIssueWorklogIDs(issueName string) ([]string, error) {
	url := fmt.Sprintf("https://%s/%s/issue/%s/worklog", r.jiraAddress, jiraV3BaseURL, issueName)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(r.email, r.authToken)

	response, err := r.client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, errIssueNotFound
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var i issue
	if err := json.Unmarshal(body, &i); err != nil {
		return nil, err
	}

	worklogIDs := make([]string, len(i.Worklogs))
	for j := range i.Worklogs {
		worklogIDs[j] = i.Worklogs[j].ID
	}

	return worklogIDs, nil
}
