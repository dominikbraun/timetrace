package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dominikbraun/timetrace/config"
)

const (
	rootDirName     = ".timetrace"
	projectsDirName = "projects"
	recordsDirName  = "records"
	reportDirName   = "reports"
)

const (
	recordDirLayout            = "2006-01-02"
	recordFilepathLayout       = "15-04.json"
	recordBackupFilepathLayout = "15-04.json.bak"
)

type Fs struct {
	config    *config.Config
	sanitizer *strings.Replacer
}

func New(config *config.Config) *Fs {
	return &Fs{
		config:    config,
		sanitizer: strings.NewReplacer("/", "-", "\\", "-"),
	}
}

// ProjectFilepath returns the filepath of the project with the given key.
func (fs *Fs) ProjectFilepath(key string) string {
	key = fs.sanitizer.Replace(key)
	name := fmt.Sprintf("%s.json", key)
	return filepath.Join(fs.projectsDir(), name)
}

// ProjectBackupFilepath return the filepath of the backup project with the
// given key.
func (fs *Fs) ProjectBackupFilepath(key string) string {
	key = fs.sanitizer.Replace(key)
	name := fmt.Sprintf("%s.json.bak", key)
	return filepath.Join(fs.projectsDir(), name)
}

// ProjectFilepaths returns all non-backup project filepaths sorted alphabetically.
func (fs *Fs) ProjectFilepaths() ([]string, error) {
	dir := fs.projectsDir()

	items, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filepaths []string

	for _, item := range items {
		if item.IsDir() {
			continue
		}
		itemName := item.Name()
		if strings.HasPrefix(itemName, ".") {
			continue
		}
		if strings.HasSuffix(itemName, ".bak") {
			continue
		}

		filepaths = append(filepaths, filepath.Join(dir, itemName))
	}
	sort.Strings(filepaths)

	return filepaths, nil
}

// ProjectBackupFilepaths returns all backup project filepaths sorted alphabetically.
func (fs *Fs) ProjectBackupFilepaths() ([]string, error) {
	dir := fs.projectsDir()

	items, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var filepaths []string

	for _, item := range items {
		if item.IsDir() {
			continue
		}
		itemName := item.Name()
		if strings.HasPrefix(itemName, ".") {
			continue
		}
		if !strings.HasSuffix(itemName, ".bak") {
			continue
		}

		filepaths = append(filepaths, filepath.Join(dir, itemName))
	}
	sort.Strings(filepaths)

	return filepaths, nil
}

// RecordFilepath returns the filepath of the record with the given start time.
//
// Note that the start time also has to contain the date as this determines the
// directory the project is stored in.
func (fs *Fs) RecordFilepath(start time.Time) string {
	name := start.Format(recordFilepathLayout)
	return filepath.Join(fs.RecordDirFromDate(start), name)
}

func (fs *Fs) RecordBackupFilepath(start time.Time) string {
	name := start.Format(recordBackupFilepathLayout)
	return filepath.Join(fs.RecordDirFromDate(start), name)
}

// RecordFilepaths returns all non-backup record filepaths within the given
// directory sorted by the given function.
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
		itemName := item.Name()
		if strings.HasPrefix(itemName, ".") {
			continue
		}
		if strings.HasSuffix(itemName, ".bak") {
			continue
		}

		filepaths = append(filepaths, filepath.Join(dir, itemName))
	}

	sort.Slice(filepaths, func(i, j int) bool {
		return less(filepaths[i], filepaths[j])
	})

	return filepaths, nil
}

// RecordDirs returns all record directories sorted alphabetically. This can be
// used to determine the latest record directory and obtain the latest record
// within that directory using RecordFilepaths.
//
// Note that all timetrace directories have to exist for RecordDirs to work.
func (fs *Fs) RecordDirs() ([]string, error) {
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

	sort.Strings(dirs)

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
		fs.recordsInitSubDir(),
		fs.ReportDir(),
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

func (fs *Fs) recordsInitSubDir() string {
	return fs.RecordDirFromDate(time.Now())
}

func (fs *Fs) ReportDir() string {
	return path.Join(fs.rootDir(), reportDirName)
}

func (fs *Fs) rootDir() string {
	if fs.config.Store != "" {
		return os.ExpandEnv(fs.config.Store)
	}

	homeDir, _ := os.UserHomeDir()

	return filepath.Join(homeDir, rootDirName)
}

func (fs *Fs) WriteReport(filepath string, data []byte) error {
	reportPath := filepath
	// no out flag provided
	// TODO: looks a little irritating could be re-written
	if reportPath == "" {
		reportPath = fs.config.ReportPath
		if reportPath == "" {
			// fileName -> report-<time.unix>
			fileName := strings.Join([]string{"report", strconv.FormatInt(time.Now().Unix(), 10)}, "-")
			reportPath = path.Join(fs.ReportDir(), fileName)
		}
	}

	if err := ioutil.WriteFile(reportPath, data, 0644); err != nil {
		return err
	}
	return nil
}
