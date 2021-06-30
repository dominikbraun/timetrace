package actions

import (
	"github.com/aligator/goplug/goplug"
	core0 "github.com/dominikbraun/timetrace/core"
	time0 "time"
)

// HostActions contains the host-implementations of actions.
type HostActions struct {
	Core0TimetraceRef *core0.Timetrace
}

type ClientActions struct {
	client *goplug.Client
}

func NewClientActions(plugin *goplug.Client) ClientActions {
	return ClientActions{
		client: plugin,
	}
}

// Make some plugin-methods available to the plugins.

func (c *ClientActions) Print(text string) error {
	return c.client.Print(text)
}

// Action implementations for host and client.

type LoadProjectRequest struct {
	Key string `json:"key"`
}

type LoadProjectResponse struct {
	Res0 core0.Project `json:"res0"`
}

// LoadProject loads the project with the given key. Returns ErrProjectNotFound
// if the project cannot be found.
func (h *HostActions) LoadProject(args LoadProjectRequest, reply *LoadProjectResponse) error {
	// Host implementation.
	res0, err := h.Core0TimetraceRef.LoadProject(
		args.Key,
	)

	if err != nil {
		return err
	}

	*reply = LoadProjectResponse{
		Res0: res0,
	}

	return nil
}

// LoadProject loads the project with the given key. Returns ErrProjectNotFound
// if the project cannot be found.
func (c *ClientActions) LoadProject(
	key string,
) (res0 core0.Project, err error) {
	// Calling from the plugin.
	response := LoadProjectResponse{}
	err = c.client.Call("LoadProject", LoadProjectRequest{
		Key: key,
	}, &response)
	return response.Res0, err
}

type LoadBackupProjectRequest struct {
	Key string `json:"key"`
}

type LoadBackupProjectResponse struct {
	Res0 core0.Project `json:"res0"`
}

func (h *HostActions) LoadBackupProject(args LoadBackupProjectRequest, reply *LoadBackupProjectResponse) error {
	// Host implementation.
	res0, err := h.Core0TimetraceRef.LoadBackupProject(
		args.Key,
	)

	if err != nil {
		return err
	}

	*reply = LoadBackupProjectResponse{
		Res0: res0,
	}

	return nil
}

func (c *ClientActions) LoadBackupProject(
	key string,
) (res0 core0.Project, err error) {
	// Calling from the plugin.
	response := LoadBackupProjectResponse{}
	err = c.client.Call("LoadBackupProject", LoadBackupProjectRequest{
		Key: key,
	}, &response)
	return response.Res0, err
}

type ListProjectModulesRequest struct {
	Project core0.Project `json:"project"`
}

type ListProjectModulesResponse struct {
	Res0 string `json:"res0"`
}

// ListProjectModules loads all modules for a project and returns their keys as a concatenated string
func (h *HostActions) ListProjectModules(args ListProjectModulesRequest, reply *ListProjectModulesResponse) error {
	// Host implementation.
	res0, err := h.Core0TimetraceRef.ListProjectModules(
		args.Project,
	)

	if err != nil {
		return err
	}

	*reply = ListProjectModulesResponse{
		Res0: res0,
	}

	return nil
}

// ListProjectModules loads all modules for a project and returns their keys as a concatenated string
func (c *ClientActions) ListProjectModules(
	project core0.Project,
) (res0 string, err error) {
	// Calling from the plugin.
	response := ListProjectModulesResponse{}
	err = c.client.Call("ListProjectModules", ListProjectModulesRequest{
		Project: project,
	}, &response)
	return response.Res0, err
}

type ListProjectsRequest struct {
}

type ListProjectsResponse struct {
	Res0 []core0.Project `json:"res0"`
}

// ListProjects loads and returns all stored projects sorted by their filenames.
// If no projects are found, an empty slice and no error will be returned.
func (h *HostActions) ListProjects(args ListProjectsRequest, reply *ListProjectsResponse) error {
	// Host implementation.
	res0, err := h.Core0TimetraceRef.ListProjects()

	if err != nil {
		return err
	}

	*reply = ListProjectsResponse{
		Res0: res0,
	}

	return nil
}

// ListProjects loads and returns all stored projects sorted by their filenames.
// If no projects are found, an empty slice and no error will be returned.
func (c *ClientActions) ListProjects() (res0 []core0.Project, err error) {
	// Calling from the plugin.
	response := ListProjectsResponse{}
	err = c.client.Call("ListProjects", ListProjectsRequest{}, &response)
	return response.Res0, err
}

type SaveProjectRequest struct {
	Project core0.Project `json:"project"`
	Force   bool          `json:"force"`
}

type SaveProjectResponse struct {
}

// SaveProject persists the given project. Returns ErrProjectAlreadyExists if
// the project already exists and saving isn't forced.
func (h *HostActions) SaveProject(args SaveProjectRequest, reply *SaveProjectResponse) error {
	// Host implementation.
	err := h.Core0TimetraceRef.SaveProject(
		args.Project,
		args.Force,
	)

	if err != nil {
		return err
	}

	return nil
}

// SaveProject persists the given project. Returns ErrProjectAlreadyExists if
// the project already exists and saving isn't forced.
func (c *ClientActions) SaveProject(
	project core0.Project,
	force bool,
) error {
	// Calling from the plugin.
	response := SaveProjectResponse{}
	err := c.client.Call("SaveProject", SaveProjectRequest{
		Project: project,
		Force:   force,
	}, &response)
	return err
}

type BackupProjectRequest struct {
	ProjectKey string `json:"projectKey"`
}

type BackupProjectResponse struct {
}

// BackupProject creates a backup of the given project file.
func (h *HostActions) BackupProject(args BackupProjectRequest, reply *BackupProjectResponse) error {
	// Host implementation.
	err := h.Core0TimetraceRef.BackupProject(
		args.ProjectKey,
	)

	if err != nil {
		return err
	}

	return nil
}

// BackupProject creates a backup of the given project file.
func (c *ClientActions) BackupProject(
	projectKey string,
) error {
	// Calling from the plugin.
	response := BackupProjectResponse{}
	err := c.client.Call("BackupProject", BackupProjectRequest{
		ProjectKey: projectKey,
	}, &response)
	return err
}

type RevertProjectRequest struct {
	ProjectKey string `json:"projectKey"`
}

type RevertProjectResponse struct {
}

func (h *HostActions) RevertProject(args RevertProjectRequest, reply *RevertProjectResponse) error {
	// Host implementation.
	err := h.Core0TimetraceRef.RevertProject(
		args.ProjectKey,
	)

	if err != nil {
		return err
	}

	return nil
}

func (c *ClientActions) RevertProject(
	projectKey string,
) error {
	// Calling from the plugin.
	response := RevertProjectResponse{}
	err := c.client.Call("RevertProject", RevertProjectRequest{
		ProjectKey: projectKey,
	}, &response)
	return err
}

type EditProjectRequest struct {
	ProjectKey string `json:"projectKey"`
}

type EditProjectResponse struct {
}

// EditProject opens the project file in the preferred or default editor .
func (h *HostActions) EditProject(args EditProjectRequest, reply *EditProjectResponse) error {
	// Host implementation.
	err := h.Core0TimetraceRef.EditProject(
		args.ProjectKey,
	)

	if err != nil {
		return err
	}

	return nil
}

// EditProject opens the project file in the preferred or default editor .
func (c *ClientActions) EditProject(
	projectKey string,
) error {
	// Calling from the plugin.
	response := EditProjectResponse{}
	err := c.client.Call("EditProject", EditProjectRequest{
		ProjectKey: projectKey,
	}, &response)
	return err
}

type DeleteProjectRequest struct {
	Project core0.Project `json:"project"`
}

type DeleteProjectResponse struct {
}

// DeleteProject removes the given project. Returns ErrProjectNotFound if the
// project doesn't exist.
func (h *HostActions) DeleteProject(args DeleteProjectRequest, reply *DeleteProjectResponse) error {
	// Host implementation.
	err := h.Core0TimetraceRef.DeleteProject(
		args.Project,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteProject removes the given project. Returns ErrProjectNotFound if the
// project doesn't exist.
func (c *ClientActions) DeleteProject(
	project core0.Project,
) error {
	// Calling from the plugin.
	response := DeleteProjectResponse{}
	err := c.client.Call("DeleteProject", DeleteProjectRequest{
		Project: project,
	}, &response)
	return err
}

type LoadRecordRequest struct {
	Start time0.Time `json:"start"`
}

type LoadRecordResponse struct {
	Res0 core0.Record `json:"res0"`
}

// LoadRecord loads the record with the given start time. Returns
// ErrRecordNotFound if the record cannot be found.
func (h *HostActions) LoadRecord(args LoadRecordRequest, reply *LoadRecordResponse) error {
	// Host implementation.
	res0, err := h.Core0TimetraceRef.LoadRecord(
		args.Start,
	)

	if err != nil {
		return err
	}

	*reply = LoadRecordResponse{
		Res0: res0,
	}

	return nil
}

// LoadRecord loads the record with the given start time. Returns
// ErrRecordNotFound if the record cannot be found.
func (c *ClientActions) LoadRecord(
	start time0.Time,
) (res0 core0.Record, err error) {
	// Calling from the plugin.
	response := LoadRecordResponse{}
	err = c.client.Call("LoadRecord", LoadRecordRequest{
		Start: start,
	}, &response)
	return response.Res0, err
}

type LoadBackupRecordRequest struct {
	Start time0.Time `json:"start"`
}

type LoadBackupRecordResponse struct {
	Res0 core0.Record `json:"res0"`
}

func (h *HostActions) LoadBackupRecord(args LoadBackupRecordRequest, reply *LoadBackupRecordResponse) error {
	// Host implementation.
	res0, err := h.Core0TimetraceRef.LoadBackupRecord(
		args.Start,
	)

	if err != nil {
		return err
	}

	*reply = LoadBackupRecordResponse{
		Res0: res0,
	}

	return nil
}

func (c *ClientActions) LoadBackupRecord(
	start time0.Time,
) (res0 core0.Record, err error) {
	// Calling from the plugin.
	response := LoadBackupRecordResponse{}
	err = c.client.Call("LoadBackupRecord", LoadBackupRecordRequest{
		Start: start,
	}, &response)
	return response.Res0, err
}

type ListRecordsRequest struct {
	Date time0.Time `json:"date"`
}

type ListRecordsResponse struct {
	Res0 []core0.Record `json:"res0"`
}

// ListRecords loads and returns all records from the given date. If no records
// are found, an empty slice and no error will be returned.
func (h *HostActions) ListRecords(args ListRecordsRequest, reply *ListRecordsResponse) error {
	// Host implementation.
	res0, err := h.Core0TimetraceRef.ListRecords(
		args.Date,
	)

	if err != nil {
		return err
	}

	*reply = ListRecordsResponse{
		Res0: res0,
	}

	return nil
}

// ListRecords loads and returns all records from the given date. If no records
// are found, an empty slice and no error will be returned.
func (c *ClientActions) ListRecords(
	date time0.Time,
) (res0 []core0.Record, err error) {
	// Calling from the plugin.
	response := ListRecordsResponse{}
	err = c.client.Call("ListRecords", ListRecordsRequest{
		Date: date,
	}, &response)
	return response.Res0, err
}

type SaveRecordRequest struct {
	Record core0.Record `json:"record"`
	Force  bool         `json:"force"`
}

type SaveRecordResponse struct {
}

// SaveRecord persists the given record. Returns ErrRecordAlreadyExists if the
// record already exists and saving isn't forced.
func (h *HostActions) SaveRecord(args SaveRecordRequest, reply *SaveRecordResponse) error {
	// Host implementation.
	err := h.Core0TimetraceRef.SaveRecord(
		args.Record,
		args.Force,
	)

	if err != nil {
		return err
	}

	return nil
}

// SaveRecord persists the given record. Returns ErrRecordAlreadyExists if the
// record already exists and saving isn't forced.
func (c *ClientActions) SaveRecord(
	record core0.Record,
	force bool,
) error {
	// Calling from the plugin.
	response := SaveRecordResponse{}
	err := c.client.Call("SaveRecord", SaveRecordRequest{
		Record: record,
		Force:  force,
	}, &response)
	return err
}

type BackupRecordRequest struct {
	RecordKey time0.Time `json:"recordKey"`
}

type BackupRecordResponse struct {
}

// BackupRecord creates a backup of the given record file
func (h *HostActions) BackupRecord(args BackupRecordRequest, reply *BackupRecordResponse) error {
	// Host implementation.
	err := h.Core0TimetraceRef.BackupRecord(
		args.RecordKey,
	)

	if err != nil {
		return err
	}

	return nil
}

// BackupRecord creates a backup of the given record file
func (c *ClientActions) BackupRecord(
	recordKey time0.Time,
) error {
	// Calling from the plugin.
	response := BackupRecordResponse{}
	err := c.client.Call("BackupRecord", BackupRecordRequest{
		RecordKey: recordKey,
	}, &response)
	return err
}

type RevertRecordRequest struct {
	RecordKey time0.Time `json:"recordKey"`
}

type RevertRecordResponse struct {
}

func (h *HostActions) RevertRecord(args RevertRecordRequest, reply *RevertRecordResponse) error {
	// Host implementation.
	err := h.Core0TimetraceRef.RevertRecord(
		args.RecordKey,
	)

	if err != nil {
		return err
	}

	return nil
}

func (c *ClientActions) RevertRecord(
	recordKey time0.Time,
) error {
	// Calling from the plugin.
	response := RevertRecordResponse{}
	err := c.client.Call("RevertRecord", RevertRecordRequest{
		RecordKey: recordKey,
	}, &response)
	return err
}

type DeleteRecordRequest struct {
	Record core0.Record `json:"record"`
}

type DeleteRecordResponse struct {
}

// DeleteRecord removes the given record. Returns ErrRecordNotFound if the
// project doesn't exist.
func (h *HostActions) DeleteRecord(args DeleteRecordRequest, reply *DeleteRecordResponse) error {
	// Host implementation.
	err := h.Core0TimetraceRef.DeleteRecord(
		args.Record,
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteRecord removes the given record. Returns ErrRecordNotFound if the
// project doesn't exist.
func (c *ClientActions) DeleteRecord(
	record core0.Record,
) error {
	// Calling from the plugin.
	response := DeleteRecordResponse{}
	err := c.client.Call("DeleteRecord", DeleteRecordRequest{
		Record: record,
	}, &response)
	return err
}

type EditRecordManualRequest struct {
	RecordTime time0.Time `json:"recordTime"`
}

type EditRecordManualResponse struct {
}

// EditRecordManual opens the record file in the preferred or default editor.
func (h *HostActions) EditRecordManual(args EditRecordManualRequest, reply *EditRecordManualResponse) error {
	// Host implementation.
	err := h.Core0TimetraceRef.EditRecordManual(
		args.RecordTime,
	)

	if err != nil {
		return err
	}

	return nil
}

// EditRecordManual opens the record file in the preferred or default editor.
func (c *ClientActions) EditRecordManual(
	recordTime time0.Time,
) error {
	// Calling from the plugin.
	response := EditRecordManualResponse{}
	err := c.client.Call("EditRecordManual", EditRecordManualRequest{
		RecordTime: recordTime,
	}, &response)
	return err
}

type EditRecordRequest struct {
	RecordTime time0.Time `json:"recordTime"`
	Plus       string     `json:"plus"`
	Minus      string     `json:"minus"`
}

type EditRecordResponse struct {
}

// EditRecord loads the record internally, applies the option values and saves the record
func (h *HostActions) EditRecord(args EditRecordRequest, reply *EditRecordResponse) error {
	// Host implementation.
	err := h.Core0TimetraceRef.EditRecord(
		args.RecordTime,
		args.Plus,
		args.Minus,
	)

	if err != nil {
		return err
	}

	return nil
}

// EditRecord loads the record internally, applies the option values and saves the record
func (c *ClientActions) EditRecord(
	recordTime time0.Time,
	plus string,
	minus string,
) error {
	// Calling from the plugin.
	response := EditRecordResponse{}
	err := c.client.Call("EditRecord", EditRecordRequest{
		RecordTime: recordTime,
		Plus:       plus,
		Minus:      minus,
	}, &response)
	return err
}

type LoadLatestRecordRequest struct {
}

type LoadLatestRecordResponse struct {
	Res0 core0.Record `json:"res0"`
}

// LoadLatestRecord loads the youngest record. This may also be a record from
// another day. If there is no latest record, nil and no error will be returned.
func (h *HostActions) LoadLatestRecord(args LoadLatestRecordRequest, reply *LoadLatestRecordResponse) error {
	// Host implementation.
	res0, err := h.Core0TimetraceRef.LoadLatestRecord()

	if err != nil {
		return err
	}

	*reply = LoadLatestRecordResponse{
		Res0: res0,
	}

	return nil
}

// LoadLatestRecord loads the youngest record. This may also be a record from
// another day. If there is no latest record, nil and no error will be returned.
func (c *ClientActions) LoadLatestRecord() (res0 core0.Record, err error) {
	// Calling from the plugin.
	response := LoadLatestRecordResponse{}
	err = c.client.Call("LoadLatestRecord", LoadLatestRecordRequest{}, &response)
	return response.Res0, err
}
