package core

import "time"

const (
	defaultTimeLayout        = "15:04"
	default12HoursTimeLayout = "03:04PM"
)

type Formatter struct {
	use12Hours bool
}

func (f *Formatter) timeLayout() string {
	if f.use12Hours {
		return default12HoursTimeLayout
	}
	return defaultTimeLayout
}

func (f *Formatter) TimeString(input time.Time) string {
	return input.Format(f.timeLayout())
}

const (
	defaultRecordKeyLayout        = "2006-01-02-15-04"
	default12HoursRecordKeyLayout = "2006-01-02-03-04PM"
)

func (f *Formatter) RecordKeyLayout() string {
	if f.use12Hours {
		return default12HoursRecordKeyLayout
	}
	return defaultRecordKeyLayout
}

func (f *Formatter) ParseRecordKeyString(recordKey string) (time.Time, error) {
	return time.Parse(f.RecordKeyLayout(), recordKey)
}

func (f *Formatter) RecordKeyString(record Record) string {
	return record.Start.Format(f.RecordKeyLayout())
}
