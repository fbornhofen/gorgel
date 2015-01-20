package libgorgel

import (
	"testing"
)

func newSynthesizer() *Synthesizer {
	return NewSynthesizer(120, 44100)
}

func TestAsString(t *testing.T) {
	n := NewCmdNote(24, 4, 4, ENVELOPE_RECTANGULAR, nil)
	if n.AsString() != "N 24, 4, 4" {
		t.Errorf("is \"%s\", should be \"N 24, 4, 4\"")
	}
}

func TestFirstAndLastFrame(t *testing.T) {
	s := newSynthesizer()
	// Note starting at .5s, runnin .5s
	n := NewCmdNote(24, 4, 4, ENVELOPE_RECTANGULAR, s)
	if n.BeginFrame != 22050 {
		t.Errorf("note should begin at 22050 (actual: %d)", n.BeginFrame)
	}
	if n.EndFrame != 44100 {
		t.Errorf("note should end at 44100 (actual: %d)", n.EndFrame)
	}
}

func TestSampleAt(t *testing.T) {
	s := newSynthesizer()
	// Note starting at .5s, runnin .5s
	n := NewCmdNote(24, 4, 4, ENVELOPE_RECTANGULAR, s)
	var val float32
	val = n.SampleFrame(100)
	if val > 0.0 {
		t.Errorf("note starts too early")
	}
	val = n.SampleFrame(23000)
	if val < 0.0001 {
		t.Errorf("note should produce nonzero value while playing")
	}
	val = n.SampleFrame(44101)
	if val > 0.0 {
		t.Errorf("note should stop after 1s")
	}
}
