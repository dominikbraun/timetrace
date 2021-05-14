package core

import (
	"errors"
	"time"

	"github.com/dominikbraun/timetrace/config"
)

var (
	ErrNoEndTime          = errors.New("no end time for last record")
	ErrTrackingNotStarted = errors.New("start tracking first")
)

type Report struct {
	Current            *Record
	TrackedTimeCurrent *time.Duration
	TrackedTimeToday   time.Duration
}

// Filesystem represents a filesystem used for storing and loading resources.
type Filesystem interface {
	ProjectFilepath(key string) string
	ProjectFilepaths() ([]string, error)
	RecordFilepath(start time.Time) string
	RecordFilepaths(dir string, less func(a, b string) bool) ([]string, error)
	RecordDirs() ([]string, error)
	RecordDirFromDate(date time.Time) string
	EnsureDirectories() error
	EnsureRecordDir(date time.Time) error
}

type Timetrace struct {
	config *config.Config
	fs     Filesystem
}

func New(config *config.Config, fs Filesystem) *Timetrace {
	return &Timetrace{
		config: config,
		fs:     fs,
	}
}

// Start starts tracking time for the given project key. This will create a new
// record with the current time as start time.
//
// Since parallel work isn't supported, the previous work must be stopped first.
func (t *Timetrace) Start(projectKey string, isBillable bool) error {
	latestRecord, err := t.loadLatestRecord()
	if err != nil {
		return err
	}

	// If there is no end time of the latest record, the user has to stop first.
	if latestRecord != nil && latestRecord.End == nil {
		return ErrNoEndTime
	}

	var project *Project

	if projectKey != "" {
		if project, err = t.LoadProject(projectKey); err != nil {
			return err
		}
	}

	record := Record{
		Start:      time.Now(),
		Project:    project,
		IsBillable: isBillable,
	}

	return t.SaveRecord(record, false)
}

// Status calculates and returns a status report.
//
// If the user isn't tracking time at the moment of calling this function, the
// Report.Current and Report.TrackedTimeCurrent fields will be nil. If the user
// hasn't tracked time today, ErrTrackingNotStarted will be returned.
func (t *Timetrace) Status() (*Report, error) {
	now := time.Now()

	firstRecord, err := t.loadOldestRecord(time.Now())
	if err != nil {
		return nil, err
	}

	if firstRecord == nil {
		return nil, ErrTrackingNotStarted
	}

	latestRecord, err := t.loadLatestRecord()
	if err != nil {
		return nil, err
	}

	trackedTimeToday, err := t.trackedTime(now)
	if err != nil {
		return nil, err
	}

	report := &Report{
		TrackedTimeToday: trackedTimeToday,
	}

	// If the latest record has been stopped, there is no active time tracking.
	// Therefore, just calculate the tracked time of today and return.
	if latestRecord.End != nil {
		return report, nil
	}

	report.Current = latestRecord

	// If the latest record has not been stopped yet, time tracking is active.
	// Calculate the time tracked for the current record and for today.
	trackedTimeCurrent := now.Sub(latestRecord.Start)
	report.TrackedTimeCurrent = &trackedTimeCurrent

	return report, nil
}

// Stop stops the time tracking and marks the current record as ended.
func (t *Timetrace) Stop() error {
	latestRecord, err := t.loadLatestRecord()
	if err != nil {
		return err
	}

	if latestRecord == nil || latestRecord.End != nil {
		return ErrTrackingNotStarted
	}

	end := time.Now()
	latestRecord.End = &end

	return t.SaveRecord(*latestRecord, false)
}

func (t *Timetrace) EnsureDirectories() error {
	return t.fs.EnsureDirectories()
}

func (t *Timetrace) Config() *config.Config {
	return t.config
}

func (t *Timetrace) trackedTime(date time.Time) (time.Duration, error) {
	records, err := t.loadAllRecords(date)
	if err != nil {
		return 0, err
	}

	var trackedTime time.Duration

	for _, record := range records {
		// If the record doesn't have an end time, it is expected that this is
		// the current record and time is still being tracked.
		if record.End == nil {
			trackedTime += time.Now().Sub(record.Start)
			continue
		}

		trackedTime += record.End.Sub(record.Start)
	}

	return trackedTime, nil
}
