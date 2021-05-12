package fs

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"time"

	"github.com/dominikbraun/timetrace/config"
)

const (
	defaultRecordFilepathLayout = "15-04.json"
)

// RecordFilepaths returns all record filepaths within the given directory
// sorted by the given function.
//
// The directory can be obtained using functions like RecordDir or RecordDirs.
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
func RecordFilepaths(dir string, less func(a, b string) bool) ([]string, error) {
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

// ProjectFilepath returns the filepath of the project with the given key.
func ProjectFilepath(key string) string {
	name := fmt.Sprintf("%s.json", key)
	return filepath.Join(projectsDir(), name)
}

// RecordFilepath returns the filepath of the record with the given name.
//
// By default, a record started at 3:00 PM will be stored in a file called
// 15-00.json. If use12hours is set in the config, it will be 03-00PM.json.
func RecordFilepath(start time.Time) string {
	layout := defaultRecordFilepathLayout

	if config.Get().Use12Hours {
		layout = "03-04PM.json"
	}

	name := start.Format(layout)
	return filepath.Join(RecordDirFromDate(start), name)
}
