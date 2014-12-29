package libgorgel

type CmdDummy struct {
	beginQuarterBeats    int
	durationQuarterBeats int
}

func NewCmdDummy(begin int, duration int) *CmdDummy {
	d := new(CmdDummy)
	d.beginQuarterBeats = begin
	d.durationQuarterBeats = duration
	return d
}

func (d *CmdDummy) AsString() string {
	return "D"
}

func (d *CmdDummy) BeginQuarterBeats() int {
	return d.beginQuarterBeats
}

func (d *CmdDummy) DurationQuarterBeats() int {
	return d.durationQuarterBeats
}

func (d *CmdDummy) SampleFrame(f int) int16 {
	return 0
}
