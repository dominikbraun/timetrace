package core

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/dominikbraun/timetrace/config"
	"github.com/dominikbraun/timetrace/fs"
)

func TestRenameRecordFileIfNeeded(t *testing.T) {
	dir, err := os.MkdirTemp("", "timetrace-record-rename-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	cfg := &config.Config{Store: dir}
	filesystem := fs.New(cfg)
	timetrace := New(cfg, filesystem)

	loc := time.Local
	originalStart := time.Date(2024, 2, 5, 16, 57, 24, 0, loc)
	updatedStart := time.Date(2024, 2, 5, 16, 55, 24, 0, loc)
	end := time.Date(2024, 2, 5, 16, 57, 34, 0, loc)

	record := Record{
		Start: originalStart,
		End:   &end,
		Project: &Project{
			Key: "dummy",
		},
	}

	if err := timetrace.SaveRecord(record, false); err != nil {
		t.Fatalf("SaveRecord: %v", err)
	}

	if err := timetrace.BackupRecord(originalStart); err != nil {
		t.Fatalf("BackupRecord: %v", err)
	}

	oldPath := filesystem.RecordFilepath(originalStart)
	updated := record
	updated.Start = updatedStart
	bytes, err := json.MarshalIndent(&updated, "", "\t")
	if err != nil {
		t.Fatalf("marshal record: %v", err)
	}
	if err := os.WriteFile(oldPath, bytes, 0600); err != nil {
		t.Fatalf("write updated record: %v", err)
	}

	if err := timetrace.renameRecordFileIfNeeded(originalStart); err != nil {
		t.Fatalf("renameRecordFileIfNeeded: %v", err)
	}

	newPath := filesystem.RecordFilepath(updatedStart)
	if _, err := os.Stat(newPath); err != nil {
		t.Fatalf("expected renamed record at %s: %v", newPath, err)
	}
	if _, err := os.Stat(oldPath); !os.IsNotExist(err) {
		t.Fatalf("expected old record path removed: %s", oldPath)
	}

	backupPath := filesystem.RecordBackupFilepath(updatedStart)
	if _, err := os.Stat(backupPath); err != nil {
		t.Fatalf("expected backup renamed to %s: %v", backupPath, err)
	}

	loaded, err := timetrace.LoadRecord(updatedStart)
	if err != nil {
		t.Fatalf("LoadRecord with updated key: %v", err)
	}
	if !loaded.Start.Equal(updatedStart) {
		t.Fatalf("expected start %v, got %v", updatedStart, loaded.Start)
	}
}

func TestRenameRecordFileIfNeededNoOpWhenUnchanged(t *testing.T) {
	dir, err := os.MkdirTemp("", "timetrace-record-rename-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	cfg := &config.Config{Store: dir}
	filesystem := fs.New(cfg)
	timetrace := New(cfg, filesystem)

	start := time.Date(2024, 2, 5, 16, 57, 0, 0, time.Local)
	end := start.Add(2 * time.Minute)
	record := Record{
		Start:   start,
		End:     &end,
		Project: &Project{Key: "dummy"},
	}

	if err := timetrace.SaveRecord(record, false); err != nil {
		t.Fatalf("SaveRecord: %v", err)
	}

	if err := timetrace.renameRecordFileIfNeeded(start); err != nil {
		t.Fatalf("renameRecordFileIfNeeded: %v", err)
	}

	path := filesystem.RecordFilepath(start)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("record should remain at %s: %v", path, err)
	}
}