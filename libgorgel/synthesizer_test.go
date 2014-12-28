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
