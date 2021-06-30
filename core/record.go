package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	recordLayout = "15-04"
)

var (
	ErrRecordNotFound       = errors.New("record not found")
	ErrBackupRecordNotFound = errors.New("backup record not found")
	ErrRecordAlreadyExists  = errors.New("record already exists")
)

type Record struct {
	Start      time.Time  `json:"start"`
	End        *time.Time `json:"end"`
	Project    *Project   `json:"project"`
	IsBillable bool       `json:"is_billable"`
}

func unwrapRecordPointerResult(record *Record, err error) (Record, error) {
	if err != nil {
		return Record{}, err
	}
	return *record, err
}

// LoadRecord loads the record with the given start time. Returns
// ErrRecordNotFound if the record cannot be found.
//goplug:generate
func (t *Timetrace) LoadRecord(start time.Time) (Record, error) {
	path := t.fs.RecordFilepath(start)
	return unwrapRecordPointerResult(t.loadRecord(path))
}

//goplug:generate
func (t *Timetrace) LoadBackupRecord(start time.Time) (Record, error) {
	path := t.fs.RecordBackupFilepath(start)
	return unwrapRecordPointerResult(t.loadRecord(path))
}

// ListRecords loads and returns all records from the given date. If no records
// are found, an empty slice and no error will be returned.
//goplug:generate
func (t *Timetrace) ListRecords(date time.Time) ([]Record, error) {
	dir := t.fs.RecordDirFromDate(date)
	paths, err := t.fs.RecordFilepaths(dir, func(_, _ string) bool {
		return true
	})
	if err != nil {
		return nil, err
	}

	records := make([]Record, 0)

	for _, path := range paths {
		record, err := t.loadRecord(path)
		if err != nil {
			return nil, err
		}
		records = append(records, *record)
	}

	return records, nil
}

// SaveRecord persists the given record. Returns ErrRecordAlreadyExists if the
// record already exists and saving isn't forced.
//goplug:generate
func (t *Timetrace) SaveRecord(record Record, force bool) error {
	path := t.fs.RecordFilepath(record.Start)

	if _, err := os.Stat(path); err == nil && !force {
		return ErrRecordAlreadyExists
	}

	if err := t.fs.EnsureRecordDir(record.Start); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(&record, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}

// BackupRecord creates a backup of the given record file
//goplug:generate
func (t *Timetrace) BackupRecord(recordKey time.Time) error {
	path := t.fs.RecordFilepath(recordKey)
	record, err := t.loadRecord(path)
	if err != nil {
		return err
	}
	// create a new .bak filepath from the record struct
	backupPath := t.fs.RecordBackupFilepath(recordKey)

	backupFile, err := os.OpenFile(backupPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(&record, "", "\t")
	if err != nil {
		return err
	}

	_, err = backupFile.Write(bytes)

	return err
}

//goplug:generate
func (t *Timetrace) RevertRecord(recordKey time.Time) error {
	record, err := t.LoadBackupRecord(recordKey)
	if err != nil {
		return err
	}

	path := t.fs.RecordFilepath(recordKey)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(&record, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}

// DeleteRecord removes the given record. Returns ErrRecordNotFound if the
// project doesn't exist.
//goplug:generate
func (t *Timetrace) DeleteRecord(record Record) error {
	path := t.fs.RecordFilepath(record.Start)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrRecordNotFound
	}

	return os.Remove(path)
}

// EditRecordManual opens the record file in the preferred or default editor.
//goplug:generate
func (t *Timetrace) EditRecordManual(recordTime time.Time) error {
	path := t.fs.RecordFilepath(recordTime)

	if _, err := t.loadRecord(path); err != nil {
		return err
	}

	editor := t.editorFromEnvironment()
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// EditRecord loads the record internally, applies the option values and saves the record
//goplug:generate
func (t *Timetrace) EditRecord(recordTime time.Time, plus string, minus string) error {
	path := t.fs.RecordFilepath(recordTime)

	record, err := t.loadRecord(path)
	if err != nil {
		return err
	}

	err = t.editRecord(record, plus, minus)
	if err != nil {
		return err
	}

	err = t.SaveRecord(*record, true)
	if err != nil {
		return err
	}

	return nil
}

func (t *Timetrace) loadAllRecords(date time.Time) ([]*Record, error) {
	dir := t.fs.RecordDirFromDate(date)

	recordFilepaths, err := t.fs.RecordFilepaths(dir, func(_, _ string) bool {
		return true
	})
	if err != nil {
		return nil, err
	}

	var records []*Record

	for _, recordFilepath := range recordFilepaths {
		record, err := t.loadRecord(recordFilepath)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (t *Timetrace) loadAllRecordsSortedAscending(date time.Time) ([]*Record, error) {
	dir := t.fs.RecordDirFromDate(date)

	recordFilepaths, err := t.fs.RecordFilepaths(dir, func(a, b string) bool {
		timeA, _ := time.Parse(recordLayout, a)
		timeB, _ := time.Parse(recordLayout, b)
		return timeA.Before(timeB)
	})
	if err != nil {
		return nil, err
	}

	var records []*Record

	for _, recordFilepath := range recordFilepaths {
		record, err := t.loadRecord(recordFilepath)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// LoadLatestRecord loads the youngest record. This may also be a record from
// another day. If there is no latest record, nil and no error will be returned.
//goplug:generate
func (t *Timetrace) LoadLatestRecord() (Record, error) {
	latestDirs, err := t.fs.RecordDirs()
	if err != nil {
		return Record{}, err
	}

	if len(latestDirs) == 0 {
		return Record{}, nil
	}

	dir, err := t.latestNonEmptyDir(latestDirs)
	if err != nil {
		return Record{}, err
	}

	latestRecords, err := t.fs.RecordFilepaths(dir, func(a, b string) bool {
		timeA, _ := time.Parse(recordLayout, a)
		timeB, _ := time.Parse(recordLayout, b)
		return timeA.Before(timeB)
	})
	if err != nil {
		return Record{}, err
	}

	if len(latestRecords) == 0 {
		return Record{}, nil
	}

	path := latestRecords[len(latestRecords)-1]

	return unwrapRecordPointerResult(t.loadRecord(path))
}

// loadOldestRecord returns the oldest record of the given date. If there is no
// oldest record, nil and no error will be returned.
func (t *Timetrace) loadOldestRecord(date time.Time) (*Record, error) {
	dir := t.fs.RecordDirFromDate(date)

	oldestRecords, err := t.fs.RecordFilepaths(dir, func(a, b string) bool {
		timeA, _ := time.Parse(recordLayout, a)
		timeB, _ := time.Parse(recordLayout, b)
		return timeA.After(timeB)
	})
	if err != nil {
		return nil, err
	}

	if len(oldestRecords) == 0 {
		return nil, nil
	}

	path := oldestRecords[0]

	return t.loadRecord(path)
}

// loadFromRecordDir loads all records for one directory and returns them. The slice can be filtered
// through the filter options.
func (t *Timetrace) loadFromRecordDir(recordDir string, filter ...func(Record) bool) ([]Record, error) {
	filesInfo, err := ioutil.ReadDir(recordDir)
	if err != nil {
		return nil, err
	}
	var foundRecords = make([]Record, 0)

outer:
	for _, info := range filesInfo {
		record, err := t.loadRecord(filepath.Join(recordDir, info.Name()))
		if err != nil {
			return nil, err
		}
		// apply all filter on record to check if Records should be used
		for _, f := range filter {
			if !f(*record) {
				// if either filter returns false
				// skip record
				continue outer
			}
		}
		foundRecords = append(foundRecords, *record)
	}
	return foundRecords, nil
}

func (t *Timetrace) loadRecord(path string) (*Record, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if strings.HasSuffix(path, ".bak") {
				return nil, ErrBackupRecordNotFound
			}
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	var record Record

	if err := json.Unmarshal(file, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (t *Timetrace) editRecord(record *Record, plus string, minus string) error {

	if record.End == nil {
		return errors.New("record is still in progress")
	}

	var dur time.Duration
	var err error
	if plus != "" {
		dur, err = time.ParseDuration(plus)
	} else {
		dur, err = time.ParseDuration(minus)
		dur = -dur
	}
	if err != nil {
		return err
	}

	newEnd := record.End.Add(dur)
	if newEnd.Before(record.Start) {
		return errors.New("new ending time is before start time of record")
	}
	record.End = &newEnd

	return nil
}
