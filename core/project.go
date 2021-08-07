package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	defaultEditor = "vi"
)

var (
	ErrProjectNotFound       = errors.New("project not found")
	ErrBackupProjectNotFound = errors.New("backup project not found")
	ErrProjectAlreadyExists  = errors.New("project already exists")
	ErrParentlessModule      = errors.New("no parent project for module exists, please create parent first")
)

type Project struct {
	Key string `json:"key"`
}

// Parent returns the parent project of the current project or an empty string
// if there is no parent. If it has a parent, the current project is a module.
func (p *Project) Parent() string {
	tokens := strings.Split(p.Key, "@")

	if len(tokens) < 2 {
		return ""
	}

	return tokens[1]
}

func (p *Project) IsModule() bool {
	return p.Parent() != ""
}

// LoadProject loads the project with the given key. Returns ErrProjectNotFound
// if the project cannot be found.
func (t *Timetrace) LoadProject(key string) (*Project, error) {
	path := t.fs.ProjectFilepath(key)
	return t.loadProject(path)
}

func (t *Timetrace) LoadBackupProject(key string) (*Project, error) {
	path := t.fs.ProjectBackupFilepath(key)
	return t.loadProject(path)
}

// ListProjectModules loads all modules for a project and returns their keys as a concatenated string
func (t *Timetrace) ListProjectModules(project *Project) (string, error) {
	allModules, err := t.loadProjectModules(project)
	if err != nil {
		return "", err
	}

	if len(allModules) == 0 {
		return "-", nil
	}

	var mList string
	for i, p := range allModules {
		// get the name of the module without the prefix
		mList += strings.Split(p.Key, "@")[0]
		// append comma if it is not the last element
		if i+1 != len(allModules) {
			mList += ","
		}
	}

	return mList, nil
}

// ListProjects loads and returns all stored projects sorted by their filenames.
// If no projects are found, an empty slice and no error will be returned.
func (t *Timetrace) ListProjects() ([]*Project, error) {
	paths, err := t.fs.ProjectFilepaths()
	if err != nil {
		return nil, err
	}

	projects := make([]*Project, 0)

	for _, path := range paths {
		project, err := t.loadProject(path)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// SaveProject persists the given project. Returns ErrProjectAlreadyExists if
// the project already exists and saving isn't forced.
func (t *Timetrace) SaveProject(project Project, force bool) error {
	if project.IsModule() {
		err := t.assertParent(project)
		if err != nil {
			return err
		}
	}

	path := t.fs.ProjectFilepath(project.Key)
	if _, err := os.Stat(path); err == nil && !force {
		return ErrProjectAlreadyExists
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(&project, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}

// BackupProject creates a backup of the given project file.
func (t *Timetrace) BackupProject(projectKey string) error {
	project, err := t.LoadProject(projectKey)
	if err != nil {
		return err
	}

	path := t.fs.ProjectBackupFilepath(projectKey)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(&project, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}

// RevertProject reverts the given project to its latest backup
func (t *Timetrace) RevertProject(projectKey string) error {
	// get all backup filepaths
	backups, err := t.fs.ProjectBackupFilepaths()
	if err != nil {
		return err
	}

	// get submodules associated with projectKey
	var submodules []string
	for _, backup := range backups {
		if strings.HasSuffix(backup, fmt.Sprintf("@%s.json.bak", projectKey)) {
			submoduleKey := strings.TrimSuffix(filepath.Base(backup), ".json.bak")
			submodules = append(submodules, submoduleKey)
		}
	}

	// revert submodules
	for _, key := range submodules {
		if err := t.revert(key); err != nil {
			return err
		}
	}

	// revert parent project
	err = t.revert(projectKey)
	return err
}

// EditProject opens the project file in the preferred or default editor .
func (t *Timetrace) EditProject(projectKey string) error {
	if _, err := t.LoadProject(projectKey); err != nil {
		return err
	}

	editor := t.editorFromEnvironment()
	path := t.fs.ProjectFilepath(projectKey)

	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// DeleteProject removes the given project and any associated submodules. Returns ErrProjectNotFound if the
// project doesn't exist.
func (t *Timetrace) DeleteProject(project Project) error {
	// check if project has submodules
	modules, err := t.loadProjectModules(&project)
	if err != nil {
		return err
	}

	// delete all submodules if they exist
	for _, module := range modules {
		if err := t.delete(module.Key); err != nil {
			return err
		}
	}

	// delete parent project
	return t.delete(project.Key)
}

func (t *Timetrace) loadProject(path string) (*Project, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if strings.HasSuffix(path, ".bak") {
				return nil, ErrBackupProjectNotFound
			}
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	var project Project

	if err := json.Unmarshal(file, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

// loadProjectModules loads all modules of the given project.
//
// Since project modules are projects with the name <module>@<project>, this
// function simply loads all "projects" suffixed with @<key>.
func (t *Timetrace) loadProjectModules(project *Project) ([]*Project, error) {
	projects, err := t.ListProjects()
	if err != nil {
		return nil, err
	}

	var modules []*Project

	for _, p := range projects {
		if p.Parent() == project.Key {
			modules = append(modules, p)
		}
	}

	return modules, nil
}

func (t *Timetrace) editorFromEnvironment() string {
	if t.config.Editor != "" {
		return t.config.Editor
	}

	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	return defaultEditor
}

func (t *Timetrace) assertParent(project Project) error {
	allP, err := t.ListProjects()
	if err != nil {
		return err
	}

	for _, p := range allP {
		if p.Key == project.Parent() {
			return nil
		}
	}

	return ErrParentlessModule
}

func (t *Timetrace) revert(key string) error {
	project, err := t.LoadBackupProject(key)
	if err != nil {
		return err
	}

	// get filepath of reverted file
	path := t.fs.ProjectFilepath(key)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(&project, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}

func (t *Timetrace) delete(key string) error {
	path := t.fs.ProjectFilepath(key)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrProjectNotFound
	}

	return os.Remove(path)
}
