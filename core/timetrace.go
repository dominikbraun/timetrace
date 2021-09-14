package core

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dominikbraun/timetrace/config"
	"github.com/dominikbraun/timetrace/out"
)

const (
	BakFileExt = ".bak"
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
	ProjectBackupFilepaths() ([]string, error)
	RecordFilepath(start time.Time) string
	RecordBackupFilepath(start time.Time) string
	RecordFilepaths(dir string, less func(a, b string) bool) ([]string, error)
	RecordDirs() ([]string, error)
	ReportDir() string
	RecordDirFromDate(date time.Time) string
	EnsureDirectories() error
	EnsureRecordDir(date time.Time) error
	WriteReport(path string, data []byte) error
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
			useDecimalHours: config.UseDecimalHours,
			use12Hours:      config.Use12Hours,
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
	trackedTimeCurrent := latestRecord.Duration()
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

// Report generates a report of tracked times
//
// The report can be filtered by the given Filter* funcs. By default
// all records of all projects will be collected. Interaction with the report
// can be done via the Reporter instance
func (t *Timetrace) Report(filter ...func(*Record) bool) (*Reporter, error) {
	recordDirs, err := t.fs.RecordDirs()
	if err != nil {
		return nil, err
	}
	// collect records
	var result = make([]*Record, 0)
	for _, dir := range recordDirs {
		records, err := t.loadFromRecordDir(dir, filter...)
		if err != nil {
			return nil, err
		}
		result = append(result, records...)
	}

	var reporter = Reporter{
		t:      t,
		report: make(map[string][]*Record),
		totals: make(map[string]time.Duration),
	}
	// prepare data  for serialization
	reporter.sortAndMerge(result)
	return &reporter, nil
}

// WriteReport forwards the byte slice to the fs but checks in prior for
// the correct output path. If the user has not provided one the config.ReportPath
// will be used if not set path falls-back to $HOME/.timetrace/reports/report-<time.unix>
func (t *Timetrace) WriteReport(path string, data []byte) error {
	return t.fs.WriteReport(path, data)
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
		trackedTime += record.Duration()
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

func printCollisions(t *Timetrace, records []*Record) {
	out.Err("collides with these records :")

	rows := make([][]string, len(records))

	for i, record := range records {
		end := "still running"
		if record.End != nil {
			end = t.Formatter().TimeString(*record.End)
		}

		billable := "no"

		if record.IsBillable {
			billable = "yes"
		}

		rows[i] = make([]string, 6)
		rows[i][0] = strconv.Itoa(i + 1)
		rows[i][1] = t.Formatter().RecordKey(record)
		rows[i][2] = record.Project.Key
		rows[i][3] = t.Formatter().TimeString(record.Start)
		rows[i][4] = end
		rows[i][5] = billable
	}

	out.Table([]string{"#", "Key", "Project", "Start", "End", "Billable"}, rows, []string{})

	out.Warn(" start and end of the record should not overlap with others")
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

	collide, collidingRecords := collides(toCheck, allRecords)
	if collide {
		printCollisions(t, collidingRecords)
	}

	return collide, nil
}

func collides(toCheck Record, allRecords []*Record) (bool, []*Record) {
	collide := false
	collidingRecords := make([]*Record, 0)
	for _, rec := range allRecords {

		if rec.End != nil && rec.Start.Before(*toCheck.End) && rec.End.After(toCheck.Start) {
			collidingRecords = append(collidingRecords, rec)
			collide = true
		}

		if rec.End == nil && toCheck.End.After(rec.Start) {
			collidingRecords = append(collidingRecords, rec)
			collide = true
		}
	}

	return collide, collidingRecords
}

// isBackFile checks if a given filename is a backup-file
func isBakFile(filename string) bool {
	return filepath.Ext(filename) == BakFileExt
}
