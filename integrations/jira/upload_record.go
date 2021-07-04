package jira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dominikbraun/timetrace/core"
)

// newWorklogBody
// https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issue-worklogs/#api-rest-api-3-issue-issueidorkey-worklog-get
type newWorklogBody struct {
	Started          string `json:"started"`
	TimeSpentSeconds uint   `json:"timeSpentSeconds"`
}

type worklogResponse struct {
	ID string `json:"id"`
}

// UploadRecord creates a worklog in JIRA and adds metadata to show it was from
// timetrace
func (r *Repository) UploadRecord(record *core.Record) error {
	worklogID, err := r.createWorklog(record)
	if err != nil {
		return err
	}

	return r.addTimetraceMetadata(record, worklogID)
}

func (r *Repository) addTimetraceMetadata(record *core.Record, worklogID string) error {
	issueID := extractIssueID(record.Project)
	url := fmt.Sprintf("https://%s/%s/issue/%s/worklog/%s/properties/timetrace_id",
		r.jiraAddress,
		jiraV3BaseURL,
		issueID,
		worklogID,
	)

	reqBody, err := json.Marshal(timetraceProperty{
		TimetraceRecordID: r.formatter.RecordKey(record),
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequest("PUT", url, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	request.SetBasicAuth(r.email, r.authToken)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err := r.client.Do(request)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("err adding worklog to issue: %s: status: %d: %s", issueID, response.StatusCode, string(body))
	}

	return nil
}

// createWorklog takes a record and attempts to create a worklog in JIRA
func (r *Repository) createWorklog(record *core.Record) (string, error) {
	issueID := extractIssueID(record.Project)
	url := fmt.Sprintf("https://%s/%s/issue/%s/worklog",
		r.jiraAddress,
		jiraV3BaseURL,
		issueID,
	)

	requestBody, err := createWorklogBody(record)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(requestBody))
	if err != nil {
		return "", nil
	}
	request.SetBasicAuth(r.email, r.authToken)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")

	response, err := r.client.Do(request)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("err pushing to issue: %s: status: %d: %s", issueID, response.StatusCode, string(body))
	}

	var wlr worklogResponse
	err = json.Unmarshal(body, &wlr)
	return wlr.ID, err
}

func createWorklogBody(record *core.Record) ([]byte, error) {
	timeSpentSeconds := uint(record.End.Sub(record.Start).Seconds())
	if timeSpentSeconds < 60 {
		return nil, errors.New("unfortunately JIRA only accepts log durations > 60 seconds")
	}
	b := newWorklogBody{
		// TODO: sort out the timezone
		Started:          record.Start.Format("2006-01-02T15:04:05.999+0000"),
		TimeSpentSeconds: timeSpentSeconds,
	}

	return json.Marshal(b)
}
