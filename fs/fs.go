package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/dominikbraun/timetrace/config"
)

const (
	rootDirName     = ".timetrace"
	projectsDirName = "projects"
	recordsDirName  = "records"
)

const (
	recordDirLayout      = "2006-01-02"
	recordFilepathLayout = "15-04.json"
)

type Fs struct {
	config *config.Config
}

func New(config *config.Config) *Fs {
	return &Fs{config: config}
}

// ProjectFilepath returns the filepath of the project with the given key.
func (fs *Fs) ProjectFilepath(key string) string {
	name := fmt.Sprintf("%s.json", key)
	return filepath.Join(fs.projectsDir(), name)
}

// RecordFilepath returns the filepath of the record with the given name.
//
// By default, a record started at 3:00 PM will be stored in a file called
// 15-00.json. If use12hours is set in the config, it will be 03-00PM.json.
func (fs *Fs) RecordFilepath(start time.Time) string {
	name := start.Format(recordFilepathLayout)
	return filepath.Join(fs.RecordDirFromDate(start), name)
}

// RecordFilepaths returns all record filepaths within the given directory
// sorted by the given function.
//
// The directory can be obtained using functions like recordDir or RecordDirs.
// If you have a record date, use RecordDirFromDate to get the directory name.
//
// The less function allows you to sort the records. Assume three record files:
//
//	- timetrace/records/2021-05-01/08-00.json
//	- timetrace/records/2021-05-01/10-00.json
//	- timetrace/records/2021-05-01/11-30.json
//
// The following call to RecordFilepaths will return the paths of those records
// sorted from newest to oldest:
//
//	latestRecords, err := RecordFilepaths(dir, func (a, b string) bool {
//		timeA, _ := time.Parse("15-04.json", a)
//		timeB, _ := time.Parse("15-04.json", b)
//		return timeA.Before(timeB)
//	})
//
// This can be used to determine the latest record in a given record directory.
func (fs *Fs) RecordFilepaths(dir string, less func(a, b string) bool) ([]string, error) {
	items, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	filepaths := make([]string, 0)

	for _, item := range items {
		if item.IsDir() {
			continue
		}
		filepaths = append(filepaths, filepath.Join(dir, item.Name()))
	}

	sort.Slice(filepaths, func(i, j int) bool {
		return less(filepaths[i], filepaths[j])
	})

	return filepaths, nil
}

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
func (fs *Fs) RecordDirs(less func(a, b string) bool) ([]string, error) {
	items, err := ioutil.ReadDir(fs.recordsDir())
	if err != nil {
		return nil, err
	}

	var dirs []string

	for _, item := range items {
		if !item.IsDir() {
			continue
		}
		dirs = append(dirs, filepath.Join(fs.recordsDir(), item.Name()))
	}

	sort.Slice(dirs, func(i, j int) bool {
		return less(dirs[i], dirs[j])
	})

	return dirs, nil
}

func (fs *Fs) RecordDirFromDate(date time.Time) string {
	dir := date.Format(recordDirLayout)
	return fs.recordDir(dir)
}

// EnsureDirectories creates all required timetrace directories. If they already
// exist, nothing happens.
func (fs *Fs) EnsureDirectories() error {
	dirs := []string{
		fs.projectsDir(),
		fs.recordsDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return err
		}
	}

	return nil
}

func (fs *Fs) EnsureRecordDir(date time.Time) error {
	return os.MkdirAll(fs.RecordDirFromDate(date), 0777)
}

func (fs *Fs) recordDir(name string) string {
	return filepath.Join(fs.recordsDir(), name)
}

func (fs *Fs) projectsDir() string {
	return filepath.Join(fs.rootDir(), projectsDirName)
}

func (fs *Fs) recordsDir() string {
	return filepath.Join(fs.rootDir(), recordsDirName)
}

func (fs *Fs) rootDir() string {
	if fs.config.Store != "" {
		return os.ExpandEnv(fs.config.Store)
	}

	homeDir, _ := os.UserHomeDir()

	return filepath.Join(homeDir, rootDirName)
}
