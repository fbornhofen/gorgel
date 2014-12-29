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
