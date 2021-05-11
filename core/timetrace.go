package core

import (
	"errors"
	"time"
)

var (
	ErrNoEndTime          = errors.New("no end time for last record")
	ErrTrackingNotStarted = errors.New("start tracking first")
)

type Report struct {
	Current            *Record
	TrackedTimeCurrent time.Duration
	TrackedTimeToday   time.Duration
}

// Start starts tracking time for the given project key. This will create a new
// record with the current time as start time.
//
// Since parallel work isn't supported, the previous work must be stopped first.
func Start(projectKey string, isBillable bool) error {
	latestRecord, err := loadLatestRecord()
	if err != nil {
		return err
	}

	// If there is no end time of the latest record, the user has to stop first.
	if latestRecord != nil && latestRecord.End == nil {
		return ErrNoEndTime
	}

	var project *Project

	if projectKey != "" {
		if project, err = LoadProject(projectKey); err != nil {
			return err
		}
	}

	record := Record{
		Start:      time.Now(),
		Project:    project,
		IsBillable: isBillable,
	}

	return SaveRecord(record, false)
}

// Status creates and returns a report including the time worked since the last
// start and the overall time worked today.
func Status() (*Report, error) {
	latestRecord, err := loadLatestRecord()
	if err != nil {
		return nil, err
	}

	if latestRecord == nil {
		return nil, ErrTrackingNotStarted
	}

	now := time.Now()

	firstRecord, err := loadOldestRecord(now)
	if err != nil {
		return nil, err
	}

	trackedTimeCurrent := now.Sub(latestRecord.Start)
	trackedTimeToday := now.Sub(firstRecord.Start)

	return &Report{
		Current:            latestRecord,
		TrackedTimeCurrent: trackedTimeCurrent,
		TrackedTimeToday:   trackedTimeToday,
	}, nil
}

// Stop stops the time tracking and marks the current record as ended.
func Stop() error {
	latestRecord, err := loadLatestRecord()
	if err != nil {
		return err
	}

	if latestRecord == nil || latestRecord.End != nil {
		return ErrTrackingNotStarted
	}

	end := time.Now()
	latestRecord.End = &end

	return SaveRecord(*latestRecord, false)
}
