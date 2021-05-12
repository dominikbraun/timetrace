package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
)

const (
	defaultEditor = "vi"
)

var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectAlreadyExists = errors.New("project already exists")
)

type Project struct {
	Key string `json:"key"`
}

// LoadProject loads the project with the given key. Returns ErrProjectNotFound
// if the project cannot be found.
func (t *Timetrace) LoadProject(key string) (*Project, error) {
	path := t.fs.ProjectFilepath(key)

	file, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
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

// SaveProject persists the given project. Returns ErrProjectAlreadyExists if
// the project already exists and saving isn't forced.
func (t *Timetrace) SaveProject(project Project, force bool) error {
	path := t.fs.ProjectFilepath(project.Key)

	if _, err := os.Stat(path); os.IsExist(err) && !force {
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

// EditProject opens the project file in the preferred or default editor.
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

// DeleteProject removes the given project. Returns ErrProjectNotFound if the
// project doesn't exist.
func (t *Timetrace) DeleteProject(project Project) error {
	path := t.fs.ProjectFilepath(project.Key)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrProjectNotFound
	}

	return os.Remove(path)
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
