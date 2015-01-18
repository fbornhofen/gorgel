package libgorgel

import ()

type Command interface {
	AsString() string
	BeginQuarterBeats() int
	DurationQuarterBeats() int
	FirstSample() int
	LastSample() int
	SampleFrame(f int) int16
}
