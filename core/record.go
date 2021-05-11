package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/dominikbraun/timetrace/fs"
)

var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrRecordAlreadyExists = errors.New("record already exists")
)

type Record struct {
	Start      time.Time  `json:"start"`
	End        *time.Time `json:"end"`
	Project    *Project   `json:"project"`
	IsBillable bool       `json:"is_billable"`
}

// LoadRecord loads the record with the given start time. Returns
// ErrRecordNotFound if the record cannot be found.
func LoadRecord(start time.Time) (*Record, error) {
	path := fs.RecordFilepath(start)
	return loadRecord(path)
}

// SaveRecord persists the given record. Returns ErrRecordAlreadyExists if the
// record already exists and saving isn't forced.
func SaveRecord(record Record, force bool) error {
	path := fs.RecordFilepath(record.Start)

	if _, err := os.Stat(path); os.IsExist(err) && !force {
		return ErrRecordAlreadyExists
	}

	if err := fs.EnsureRecordDir(record.Start); err != nil {
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

// DeleteRecord removes the given record. Returns ErrRecordNotFound if the
// project doesn't exist.
func DeleteRecord(record Record) error {
	path := fs.RecordFilepath(record.Start)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrRecordNotFound
	}

	return os.Remove(path)
}

// loadLatestRecord loads the youngest record. This may also be a record from
// another day. If there is no latest record, nil and no error will be returned.
func loadLatestRecord() (*Record, error) {
	latestDirs, err := fs.RecordDirs(func(a, b string) bool {
		dateA, _ := time.Parse("2006-01-02", a)
		dateB, _ := time.Parse("2006-01-02", a)
		return dateA.Before(dateB)
	})
	if err != nil {
		return nil, err
	}

	if len(latestDirs) == 0 {
		return nil, nil
	}

	dir := latestDirs[0]

	latestRecords, err := fs.RecordFilepaths(dir, func(a, b string) bool {
		timeA, _ := time.Parse("15-04", a)
		timeB, _ := time.Parse("15-04", b)
		return timeA.Before(timeB)
	})
	if err != nil {
		return nil, err
	}

	if len(latestRecords) == 0 {
		return nil, nil
	}

	path := latestRecords[0]

	return loadRecord(path)
}

// loadOldestRecord returns the oldest record of the given date. If there is no
// oldest record, nil and no error will be returned.
func loadOldestRecord(date time.Time) (*Record, error) {
	dir := fs.RecordDirFromDate(date)

	oldestRecords, err := fs.RecordFilepaths(dir, func(a, b string) bool {
		timeA, _ := time.Parse("15-04", a)
		timeB, _ := time.Parse("15-04", b)
		return timeA.After(timeB)
	})
	if err != nil {
		return nil, err
	}

	if len(oldestRecords) == 0 {
		return nil, nil
	}

	path := oldestRecords[0]

	return loadRecord(path)
}

func loadRecord(path string) (*Record, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
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
