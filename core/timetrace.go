package core

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/dominikbraun/timetrace/config"
)

var (
	ErrNoEndTime           = errors.New("no end time for last record")
	ErrTrackingNotStarted  = errors.New("start tracking first")
	ErrAllDirectoriesEmpty = errors.New("all directories empty")
)

type Report struct {
	Current            *Record
	TrackedTimeCurrent *time.Duration
	TrackedTimeToday   time.Duration
	BreakTimeToday     time.Duration
}

// Filesystem represents a filesystem used for storing and loading resources.
type Filesystem interface {
	ProjectFilepath(key string) string
	ProjectBackupFilepath(key string) string
	ProjectFilepaths() ([]string, error)
	RecordFilepath(start time.Time) string
	RecordBackupFilepath(start time.Time) string
	RecordFilepaths(dir string, less func(a, b string) bool) ([]string, error)
	RecordDirs() ([]string, error)
	RecordDirFromDate(date time.Time) string
	EnsureDirectories() error
	EnsureRecordDir(date time.Time) error
}

type Timetrace struct {
	config    *config.Config
	fs        Filesystem
	formatter *Formatter
}

func New(config *config.Config, fs Filesystem) *Timetrace {
	return &Timetrace{
		config: config,
		fs:     fs,
		formatter: &Formatter{
			use12Hours: config.Use12Hours,
		},
	}
}

// Start starts tracking time for the given project key. This will create a new
// record with the current time as start time.
//
// Since parallel work isn't supported, the previous work must be stopped first.
func (t *Timetrace) Start(projectKey string, isBillable bool) error {
	latestRecord, err := t.LoadLatestRecord()
	if err != nil && !errors.Is(err, ErrAllDirectoriesEmpty) {
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

	latestRecord, err := t.LoadLatestRecord()
	if err != nil {
		return nil, err
	}

	trackedTimeToday, err := t.trackedTime(now)
	if err != nil {
		return nil, err
	}

	breakTimeToday, err := t.breakTime(now)
	if err != nil {
		return nil, err
	}

	report := &Report{
		TrackedTimeToday: trackedTimeToday,
		BreakTimeToday:   breakTimeToday,
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

func (t *Timetrace) breakTime(date time.Time) (time.Duration, error) {
	records, err := t.loadAllRecordsSortedAscending(date)
	if err != nil {
		return 0, err
	}

	// add up the time between records
	var breakTime time.Duration
	for i := 0; i < len(records)-1; i++ {
		breakTime += records[i+1].Start.Sub(*records[i].End)
	}

	return breakTime, nil
}

// Stop stops the time tracking and marks the current record as ended.
func (t *Timetrace) Stop() error {
	latestRecord, err := t.LoadLatestRecord()
	if err != nil {
		return err
	}

	if latestRecord == nil || latestRecord.End != nil {
		return ErrTrackingNotStarted
	}

	end := time.Now()
	latestRecord.End = &end

	return t.SaveRecord(*latestRecord, true)
}

func (t *Timetrace) EnsureDirectories() error {
	return t.fs.EnsureDirectories()
}

func (t *Timetrace) Config() *config.Config {
	return t.config
}

func (t *Timetrace) Formatter() *Formatter {
	return t.formatter
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

func (t *Timetrace) isDirEmpty(dir string) (bool, error) {
	openedDir, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer openedDir.Close()

	// Attempt to read 1 file's name, if it fails with
	// EOF, directory is empty
	_, err = openedDir.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err
}

func (t *Timetrace) latestNonEmptyDir(dirs []string) (string, error) {
	for i := len(dirs) - 1; i >= 0; i-- {
		isEmpty, err := t.isDirEmpty(dirs[i])
		if err != nil {
			return "", err
		}

		if !isEmpty {
			return dirs[i], nil
		}
	}

	return "", ErrAllDirectoriesEmpty
}

// RecordCollides checks if the time of a record collides
// with other records of the same day and returns a bool
func (t *Timetrace) RecordCollides(toCheck Record) (bool, error) {
	allRecords, err := t.loadAllRecords(toCheck.Start)
	if err != nil {
		return false, err
	}

	if toCheck.Start.Day() != toCheck.End.Day() {
		moreRecords, err := t.loadAllRecords(*toCheck.End)
		if err != nil {
			return false, err
		}
		for _, rec := range moreRecords {
			allRecords = append(allRecords, rec)
		}
	}

	return collides(toCheck, allRecords), nil
}

func collides(toCheck Record, allRecords []*Record) bool {
	for _, rec := range allRecords {
		if rec.Start.Before(toCheck.Start) && rec.End.After(toCheck.Start) {
			return true
		} else if rec.Start.Before(*toCheck.End) && rec.End.After(*toCheck.End) {
			return true
		} else if toCheck.Start.Before(rec.Start) && toCheck.End.After(rec.Start) {
			return true
		} else if toCheck.Start.Before(*rec.End) && toCheck.End.After(*rec.End) {
			return true
		}
	}

	return false
}
