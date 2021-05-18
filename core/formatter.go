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
