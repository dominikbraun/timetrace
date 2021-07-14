package core

// Provider is a slightly presumptious interface for how third party
// integrations might be used. The idea is that any integration need only
// fufill the two methods defined here, and it can be called via the push
// command
type Provider interface {
	CheckRecordsExist(records []*Record) ([]*Record, error)
	UploadRecord(records *Record) error
	Name() string
}
