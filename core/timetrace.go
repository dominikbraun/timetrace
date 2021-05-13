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
	TrackedTimeCurrent time.Duration
	TrackedTimeToday   time.Duration
}

// Filesystem represents a filesystem used for storing and loading resources.
type Filesystem interface {
	ProjectFilepath(key string) string
	ProjectFilepaths() ([]string, error)
	RecordFilepath(start time.Time) string
	RecordFilepaths(dir string, less func(a, b string) bool) ([]string, error)
	RecordDirs(less func(a, b string) bool) ([]string, error)
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

// Status creates and returns a report including the time worked since the last
// start and the overall time worked today.
func (t *Timetrace) Status() (*Report, error) {
	latestRecord, err := t.loadLatestRecord()
	if err != nil {
		return nil, err
	}

	if latestRecord == nil {
		return nil, ErrTrackingNotStarted
	}

	now := time.Now()

	firstRecord, err := t.loadOldestRecord(now)
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
