package libgorgel

import (
	"testing"
)

func TestNumSamples(t *testing.T) {
	s := NewSynthesizer(120, 44100)
	// Note (22) starting at 1.5s, running for 0.5s
	s.AddCommand(NewCmdNote(22, 12, 4, ENVELOPE_DEFAULT, s))
	n := s.NumSamples()
	if n != 88200 {
		t.Errorf("expected: 88200, actual: %d", n)
	}
}

func TestReadFromFile(t *testing.T) {
	s := NewSynthesizer(120, 44100)
	err := s.ReadFromFile("testdata/popcorn.gorgel")
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(s.commands) != 29 {
		t.Errorf("expected to read 29 commands, not %d", len(s.commands))
	}
}

func TestEnvelopeRectangular(t *testing.T) {
	s := NewSynthesizer(120, 44100)
	if s.EvalEnvelope(ENVELOPE_RECTANGULAR, 0) != 1 {
		t.Errorf("rectangular envelope should yield 1 at pos 0")
	}
	if s.EvalEnvelope(ENVELOPE_RECTANGULAR, 1) != 1 {
		t.Errorf("rectangular envelope should yield 1 at pos 1")
	}
}

func TestEnvelopeLinear(t *testing.T) {
	s := NewSynthesizer(120, 44100)
	if s.EvalEnvelope(ENVELOPE_LINEAR, 0) != 1 {
		t.Errorf("linear envelope should yield 1 at pos 0")
	}
	if s.EvalEnvelope(ENVELOPE_LINEAR, .5) != .5 {
		t.Errorf("linear envelope should yield .5 at pos .5")
	}
	if s.EvalEnvelope(ENVELOPE_LINEAR, 1) != 0 {
		t.Errorf("linear envelope should yield 0 at pos 1")
	}
}

func TestDefaultEnvelope(t *testing.T) {
	s := NewSynthesizer(120, 44100)
	s.SetDefaultEnvelope(ENVELOPE_LINEAR)
	if s.EvalEnvelope(ENVELOPE_DEFAULT, .5) != .5 {
		t.Errorf("expected linear envelope")
	}
}

func TestSampleBuffer(t *testing.T) {
	s := NewSynthesizer(120, 44100)
	s.AddCommand(NewCmdNote(25, 0, 1, ENVELOPE_DEFAULT, s))
	s.AddCommand(NewCmdNote(37, 1, 1, ENVELOPE_DEFAULT, s))
	s.AddCommand(NewCmdNote(49, 480, 1, ENVELOPE_DEFAULT, s))
	buf := make([][]float32, 1)
	buf[0] = make([]float32, 10)
	s.SampleBuffer(buf)
	l := len(s.activeCommands)
	if l != 1 {
		t.Errorf("Have %d active commands, should have 1", l)
	}
	// Sample a number of buffers so we end up in the long silence between
	// notes 37 and 49.
	for i := 0; i < 2000; i++ {
		s.SampleBuffer(buf)
	}
	l = len(s.activeCommands)
	if l != 0 {
		t.Errorf("Have %d active commands, should have 0", l)
	}
}
