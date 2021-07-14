package jira_test

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/dominikbraun/timetrace/core"
	"github.com/dominikbraun/timetrace/integrations/jira"
)

func TestUploadRecord(t *testing.T) {
	t.Run("it creates a worklog and adds the required metadata", func(t *testing.T) {
		var (
			// started 8 hours ago, finished 5 minutes ago.
			endTime    = time.Now().Add(-time.Minute * 5)
			startTime  = endTime.Add(-time.Hour * 8)
			projectKey = "project.key"
			worklogID  = "12345"
			record     = &core.Record{
				Start: startTime,
				End:   &endTime,
				Project: &core.Project{
					Key: projectKey,
				},
			}
		)

		// use the roundTripper as a poor mans mock.
		expectedRequests := []*http.Request{
			{
				Method: "POST",
				URL: &url.URL{
					Path: fmt.Sprintf("/rest/api/3/issue/%s/worklog", projectKey),
				},
			},
			{
				Method: "PUT",
				URL: &url.URL{
					Path: fmt.Sprintf("/rest/api/3/issue/%s/worklog/%s/properties/timetrace_id", projectKey, worklogID),
				},
			},
		}

		responses := []*http.Response{
			{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(strings.NewReader(fmt.Sprintf(`{"id": "%s"}`, worklogID))),
			},
			{
				StatusCode: http.StatusCreated,
			},
		}
		var i = 0
		roundTripper := roundTripFunc(func(r *http.Request) (*http.Response, error) {
			defer func() { i++ }()
			checkRequestsEqual(r, expectedRequests[i], t)
			// TODO: check that the body in the request is as expected
			return responses[i], nil
		})

		err := jira.New(jira.RepositoryConfig{
			HTTPClient: &http.Client{Transport: roundTripper},
		}).UploadRecord(record)

		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("it returns an error on non 201 from JIRA when creating worklog", func(t *testing.T) {
		var (
			// started 8 hours ago, finished 5 minutes ago.
			endTime    = time.Now().Add(-time.Minute * 5)
			startTime  = endTime.Add(-time.Hour * 8)
			projectKey = "project.key"
			record     = &core.Record{
				Start: startTime,
				End:   &endTime,
				Project: &core.Project{
					Key: projectKey,
				},
			}
		)

		roundTripper := roundTripFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
			}, nil
		})

		err := jira.New(jira.RepositoryConfig{
			HTTPClient: &http.Client{Transport: roundTripper},
		}).UploadRecord(record)

		if err == nil {
			t.Fatal("wanted an error but didn't get one")
		}
	})
}

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (s roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return s(r)
}

func checkRequestsEqual(got, want *http.Request, t *testing.T) {
	if got.Method != want.Method {
		t.Fatalf("unexpected request method. want: %s, got: %s", want.Method, got.Method)
	}

	if got.URL.Path != want.URL.Path {
		t.Fatalf("unexpected request method. want: %s, got: %s", want.URL.Path, got.URL.Path)
	}
}
