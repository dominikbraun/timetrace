package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/dominikbraun/timetrace/config"
)

const (
	rootDirectory     = ".timetrace"
	projectsDirectory = "projects"
	recordsDirectory  = "records"
)

// RecordDirs returns all record directories sorted by the given function.
//
// For example, assume three directories containing record files:
//
//	- timetrace/records/2021-05-01
//	- timetrace/records/2021-05-02
//	- timetrace/records/2021-05-03
//
// The following call to RecordDirs will return those directories sorted from
// newest to oldest:
//
//	latestRecordDirs, err := RecordDirs(func (a, b string) bool {
//		dateA, _ := time.Parse("2006-01-02", a)
//		dateB, _ := time.Parse("2006-01-02", b)
//		return dateA.Before(dateB)
//	})
//
// This can be used to determine the latest record directory and obtain the
// latest record within that directory using RecordFilepaths.
//
// Note that all timetrace directories have to exist for RecordDirs to work.
func RecordDirs(less func(a, b string) bool) ([]string, error) {
	items, err := ioutil.ReadDir(recordsDir())
	if err != nil {
		return nil, err
	}

	var dirs []string

	for _, item := range items {
		if !item.IsDir() {
			continue
		}
		dirs = append(dirs, filepath.Join(recordsDir(), item.Name()))
	}

	sort.Slice(dirs, func(i, j int) bool {
		return less(dirs[i], dirs[j])
	})

	return dirs, nil
}

// EnsureDirectories creates all required timetrace directories. If they already
// exist, nothing happens.
func EnsureDirectories() error {
	dirs := []string{
		projectsDir(),
		recordsDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}

	return nil
}

func EnsureRecordDir(date time.Time) error {
	return os.MkdirAll(RecordDirFromDate(date), 0777)
}

func RecordDirFromDate(date time.Time) string {
	dir := date.Format("2006-01-02")
	return RecordDir(dir)
}

func RecordDir(name string) string {
	return filepath.Join(recordsDir(), name)
}

func projectsDir() string {
	return filepath.Join(rootDir(), projectsDirectory)
}

func recordsDir() string {
	return filepath.Join(rootDir(), recordsDirectory)
}

func rootDir() string {
	configuredRoot := config.Get().Root

	if configuredRoot != "" {
		return os.ExpandEnv(configuredRoot)
	}

	homeDir, _ := os.UserHomeDir()

	return filepath.Join(homeDir, rootDirectory)
}
