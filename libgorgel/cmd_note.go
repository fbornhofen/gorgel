package libgorgel

import (
	"fmt"
	"math"
)

type CmdNote struct {
	ScaleIndex           int
	synthesizer          *Synthesizer
	beginQuarterBeats    int
	durationQuarterBeats int
	BeginFrame           int
	EndFrame             int
	Envelope             Envelope
}

func NewCmdNote(scaleIndex int, begin int, duration int, e Envelope, s *Synthesizer) *CmdNote {
	n := new(CmdNote)
	n.ScaleIndex = scaleIndex
	n.durationQuarterBeats = duration
	n.beginQuarterBeats = begin
	n.Envelope = e
	if s != nil {
		n.SetSynthesizer(s)
	}
	return n
}

func (n *CmdNote) AsString() string {
	return fmt.Sprintf("N %d, %d, %d", n.ScaleIndex, n.beginQuarterBeats, n.durationQuarterBeats)
}

func (n *CmdNote) SetSynthesizer(s *Synthesizer) {
	n.synthesizer = s
	bps := float32(s.BeatsPerMin) / 60
	fpb := float32(s.SampleRate) / bps
	n.BeginFrame = int(float32(n.beginQuarterBeats) / 4 * fpb)
	n.EndFrame = int(float32(n.durationQuarterBeats)/4*fpb) + n.BeginFrame
}

func (n *CmdNote) BeginQuarterBeats() int {
	return n.beginQuarterBeats
}

func (n *CmdNote) DurationQuarterBeats() int {
	return n.durationQuarterBeats
}

func (n *CmdNote) FirstSample() int {
	return n.BeginFrame
}

func (n *CmdNote) LastSample() int {
	return n.EndFrame
}

func (n *CmdNote) sampleNote(freq float32, amplitude float32, pos int) int16 {
	val := float64(amplitude) * math.Sin(float64(pos)/float64(freq)) / 2.0
	relPos := float64(pos) / (float64(n.EndFrame) - float64(n.BeginFrame))
	return int16(n.synthesizer.EvalEnvelope(n.Envelope, relPos) * val)
}

func (n *CmdNote) SampleFrame(f int) int16 {
	if f < n.BeginFrame || f > n.EndFrame {
		return 0
	}
	pos := f - n.BeginFrame
	return n.sampleNote(n.synthesizer.scale[n.ScaleIndex], 8000, pos)
}
