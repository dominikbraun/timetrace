package main

type ProjectFS interface {
	Load(key string) (*Project, error)
	Save(key string, b []byte) error
	Backup(key string, b []byte) error
}

type RecordFS interface {
	Load(path string) (*Record, error)
	Save(path, b []byte) error
	Backup(key string, b []byte) error
	Delete(path string) error
}

type Loader interface {
	Load(key string, v interface{}) error
}

type Saver interface {
	Save(key string, v interface{}) error
}

type Backuper interface {
	Backup(key string, v interface{}) error
	Revert(key string, v interface{}) error
}
