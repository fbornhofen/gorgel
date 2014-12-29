package libgorgel

import ()

type Command interface {
	AsString() string
	BeginQuarterBeats() int
	DurationQuarterBeats() int
	SampleFrame(f int) int16
}
