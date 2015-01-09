package libgorgel

import (
	"testing"
)

func TestNumSamples(t *testing.T) {
	s := NewSynthesizer(120, 44100)
	// Dummy starting at 1.5s, running for 0.5s
	s.AddCommand(NewCmdDummy(12, 4))
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