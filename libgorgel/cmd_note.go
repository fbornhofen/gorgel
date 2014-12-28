package libgorgel

import (
	"math"
)

type CmdNote struct {
	ScaleIndex           int
	synthesizer          *Synthesizer
	beginQuarterBeats    int
	durationQuarterBeats int
	BeginFrame           int
	EndFrame             int
}

func NewCmdNote(scaleIndex int, begin int, duration int, s *Synthesizer) *CmdNote {
	n := new(CmdNote)
	n.synthesizer = s
	n.ScaleIndex = scaleIndex
	n.durationQuarterBeats = duration
	n.beginQuarterBeats = begin
	bps := float32(s.BeatsPerMin) / 60
	fpb := float32(s.SampleRate) / bps
	n.BeginFrame = int(float32(n.beginQuarterBeats) / 4 * fpb)
	n.EndFrame = int(float32(n.durationQuarterBeats)/4*fpb) + n.BeginFrame
	return n
}

func (n *CmdNote) BeginQuarterBeats() int {
	return n.beginQuarterBeats
}

func (n *CmdNote) DurationQuarterBeats() int {
	return n.durationQuarterBeats
}

func sampleNote(freq float32, amplitude float32, pos int) int16 {
	val := float64(amplitude) * math.Sin(float64(pos)/float64(freq)) / 2.0
	return int16(val)
}

func (n *CmdNote) SampleFrame(f int) int16 {
	if f < n.BeginFrame || f > n.EndFrame {
		return 0
	}
	pos := f - n.BeginFrame
	return sampleNote(n.synthesizer.scale[n.ScaleIndex], 8000, pos)
}
