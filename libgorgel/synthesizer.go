package libgorgel

import (
	"fmt"
	"github.com/mkb218/gosndfile/sndfile"
	"math"
	"os"
)

type Synthesizer struct {
	BeatsPerMin int
	SampleRate  int
	Channels    int
	scale       []float32
	commands    []Command
	envelopes	[]EnvelopeFunc
}

func (s *Synthesizer) createScale() {
	a110hz := float32(s.SampleRate) / math.Pi / 2.0 / 110
	s.scale = make([]float32, 49)
	step := math.Pow(2, 1./12)
	s.scale[0] = 99999
	for i := 1; i < 49; i++ {
		s.scale[i] = a110hz / float32(math.Pow(step, float64(i)))
	}
	fillEnvelopes(&s.envelopes)
}

func NewSynthesizer(bpm int, sampleRate int) *Synthesizer {
	s := new(Synthesizer)
	s.BeatsPerMin = bpm
	s.SampleRate = sampleRate
	s.createScale()
	s.Channels = 1
	return s
}

func (s *Synthesizer) AddCommand(c Command) {
	s.commands = append(s.commands, c)
}

func (s *Synthesizer) NumSamples() int {
	maxQuarterBeat := 0
	for _, c := range s.commands {
		end := c.BeginQuarterBeats() + c.DurationQuarterBeats()
		if end > maxQuarterBeat {
			maxQuarterBeat = end
		}
	}
	bps := float32(s.BeatsPerMin) / 60.0
	return int(float32(maxQuarterBeat) * float32(s.SampleRate) / bps / 4.0)
}

func (s *Synthesizer) openAndWriteHeader(filename string) *sndfile.File {
	var info sndfile.Info
	fmt.Printf("Writing header ...")
	info.Channels = int32(s.Channels)
	info.Samplerate = int32(s.SampleRate)
	info.Format = sndfile.SF_FORMAT_WAV | sndfile.SF_FORMAT_PCM_16
	file, err := sndfile.Open(filename, sndfile.Write, &info)
	if err != nil {
		fmt.Errorf("Error opening %s, error: \n", filename, err)
		os.Exit(-1)
	}
	fmt.Printf(" done\n")
	return file
}

func (s *Synthesizer) WriteWaveFile(filename string) {
	f := s.openAndWriteHeader(filename)
	numSamples := s.NumSamples()
	fmt.Printf("Sampling %d frames\n", numSamples)
	for i := 0; i < numSamples; i++ {
		var val int16 = 0
		for _, c := range s.commands {
			val += c.SampleFrame(i)
		}
		f.WriteItems([]int16{val})
	}
}

func (s *Synthesizer) ReadFromFile(filename string) error {
	g := NewGorgelFile(filename, s)
	err := g.Read()
	s.commands = g.Commands()
	return err
}

func (s *Synthesizer) EvalEnvelope(e Envelope, relPos float64) float64 {
	return s.envelopes[e](relPos)
}
