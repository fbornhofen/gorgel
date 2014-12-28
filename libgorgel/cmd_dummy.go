package libgorgel

type CmdDummy struct {
	beginQuarterBeats    int
	durationQuarterBeats int
}

func NewCmdDummy(begin int, duration int) *CmdDummy {
	n := new(CmdDummy)
	n.beginQuarterBeats = begin
	n.durationQuarterBeats = duration
	return n
}

func (n *CmdDummy) BeginQuarterBeats() int {
	return n.beginQuarterBeats
}

func (n *CmdDummy) DurationQuarterBeats() int {
	return n.durationQuarterBeats
}

func (n *CmdDummy) SampleFrame(s *Synthesizer, f int) int16 {
	return 0
}
