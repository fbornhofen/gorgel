package libgorgel

import (
	"code.google.com/p/portaudio-go/portaudio"
	"fmt"
	"github.com/mkb218/gosndfile/sndfile"
	"math"
	"os"
)

type Synthesizer struct {
	BeatsPerMin     int
	SampleRate      int
	Channels        int
	scale           []float32
	commands        []Command
	envelopes       []EnvelopeFunc
	defaultEnvelope Envelope
	curSample       int
	notifications   chan int
}

func (s *Synthesizer) createScale() {
	a110hz := float32(s.SampleRate) / math.Pi / 2.0 / 110
	s.scale = make([]float32, 61)
	step := math.Pow(2, 1./12)
	s.scale[0] = 99999
	for i := 1; i < 61; i++ {
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
	s.defaultEnvelope = ENVELOPE_RECTANGULAR
	s.notifications = make(chan int)
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
	for s.curSample = 0; s.curSample < numSamples; s.curSample++ {
		var val int16 = 0
		for _, c := range s.commands {
			val += c.SampleFrame(s.curSample)
		}
		f.WriteItems([]int16{val})
	}
}

func (s *Synthesizer) Play() {
	portaudio.Initialize()
	defer portaudio.Terminate()
	stream, err := portaudio.OpenDefaultStream(0, 1, float64(s.SampleRate), 0, s.SampleBuffer)
	if err != nil {
		panic(err)
	}
	stream.Start()
	<-s.notifications
	stream.Stop()
}

func (s *Synthesizer) SampleBuffer(out [][]float32) {
	for i := range out[0] {
		var val float32 = 0
		for _, c := range s.commands {
			val += float32(c.SampleFrame(s.curSample))
		}
		out[0][i] = val / 0x7FFF
		s.curSample++
		if s.curSample == s.NumSamples() {
			s.notifications <- 1
			return
		}
	}
}

func (s *Synthesizer) ReadFromFile(filename string) error {
	g := NewGorgelFile(filename, s)
	err := g.Read()
	s.commands = g.Commands()
	return err
}

func (s *Synthesizer) EvalEnvelope(e Envelope, relPos float64) float64 {
	if e == ENVELOPE_DEFAULT {
		return s.envelopes[s.defaultEnvelope](relPos)
	}
	return s.envelopes[e](relPos)
}

func (s *Synthesizer) SetDefaultEnvelope(e Envelope) {
	s.defaultEnvelope = e
}
