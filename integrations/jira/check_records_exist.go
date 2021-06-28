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
	var recordsToKeep []*core.Record

	// find common projects to save on calls to jira. The user may be trying to
	// upload many records for a single project
	recordsByProjectID := buildProjectRecordMap(records)

	// for every record
	for _, records := range recordsByProjectID {
		// all records here must have the same project ID, so just take the first
		issueID := extractIssueID(records[0].Project)

		fmt.Printf("checking issueID: %s\n", issueID)

		// get the worklogs attached to the jira ticket for that project
		worklogIDs, err := r.getIssueWorklogs(issueID)
		// if the issue isn't found, ignore all the records for that project
		if errors.Is(err, errIssueNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}
		fmt.Printf("found worklog IDs: %#v\n", worklogIDs)

		// and for every worklog, check if it has a timetraceID in its "metadata"
		// (called entity properties in JIRA). It won't if the worklog hasn't come
		// from here
		for _, worklogID := range worklogIDs {
			timetraceRecordID, err := r.getWorklogTimeTraceRecord(issueID, worklogID)
			if err != nil {
				return nil, err
			}
			if timetraceRecordID == nil {
				continue
			}

			// if we do have a timetrace recordID, then there is a chance this is one
			// of the records we are trying to upload!
			for _, record := range records {
				if *timetraceRecordID != r.formatter.RecordKey(record) {
					// TODO: ignore the record if it's already been pushed
				}
			}
		}
		recordsToKeep = append(recordsToKeep, records...)
	}

	return recordsToKeep, nil
}

func buildProjectRecordMap(records []*core.Record) map[string][]*core.Record {
	m := make(map[string][]*core.Record)
	for _, record := range records {
		m[record.Project.Key] = append(m[record.Project.Key], record)
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

// getIssueWorklogs calls this endpoint
// https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-worklogs/#api-rest-api-3-issue-issueidorkey-worklog-get
// and returns a list of worklog IDs, with their timetrace records not nil if
// there are any present on the worklogs
func (r *Repository) getIssueWorklogs(issueName string) ([]string, error) {
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
