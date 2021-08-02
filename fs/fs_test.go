package fs

import (
	"fmt"
	"testing"

	"github.com/dominikbraun/timetrace/config"
	"github.com/dominikbraun/timetrace/out"
)

func TestRecordDirs(t *testing.T) {
	// test output of function
	c, err := config.FromFile()
	if err != nil {
		out.Warn("%s", err.Error())
	}

	filesystem := New(c)
	recordDirs, err := filesystem.RecordDirs()
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(recordDirs) != 2 {
		t.Fatalf("Length error: expected %v, got %v", len(recordDirs), 2)
	} else {
		fmt.Printf("%v\n", recordDirs)
	}
}
