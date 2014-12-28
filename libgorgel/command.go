package libgorgel

import ()

type Command interface {
	BeginQuarterBeats() int
	DurationQuarterBeats() int
	SampleFrame(f int) int16
}
