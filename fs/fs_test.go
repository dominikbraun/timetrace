package fs

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/dominikbraun/timetrace/config"
)

func TestProjectFilepaths(t *testing.T) {

	tt := []struct {
		name string
		mock *Fs
		list []string
		want []string
		err  error
	}{
		{
			name: "sorted project lists",
			mock: createMockFS(t),
			list: []string{"project_1", "prject_2"},
			want: []string{"project_1", "prject_2"},
			err:  nil,
		},
	}

	for _, tc := range tt {
		insertProjects(t, tc.mock, tc.list...)

		projects, err := tc.mock.ProjectFilepaths()
		if err != tc.err {
			t.Fatalf("[%s] want-err: %v, got-err: %v", tc.name, tc.err, err)
		}
		if !reflect.DeepEqual(tc.want, projects) {
			t.Fatalf("[%s] want-projects: %v, got-projects: %v", tc.name, tc.want, projects)
		}
	}
}

func createMockFS(t *testing.T) *Fs {
	c, err := config.FromFile()
	if err != nil {
		t.Fatalf("[createMockFS] could not create config: %v", err)
	}
	mock := New(c)
	if err := mock.EnsureDirectories(); err != nil {
		t.Fatalf("[createMockFS] could not ensure dirs: %v", err)
	}
	return mock
}

func insertProjects(t *testing.T, fs *Fs, projects ...string) {
	for _, project := range projects {
		if err := fs.Wrapper.MkdirAll(filepath.Join(fs.projectsDir(), project), 0777); err != nil {
			t.Fatalf("[insertProject] could not insert project: %v", err)
		}
	}
}
