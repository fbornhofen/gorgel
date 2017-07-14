package libgorgel

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
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
	commandStarts   map[int][]Command
	commandEnds     map[int][]Command
	activeCommands  []Command
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
	s.activeCommands = make([]Command, 0)
	s.commandStarts = make(map[int][]Command)
	s.commandEnds = make(map[int][]Command)
	return s
}

func (s *Synthesizer) AddCommand(c Command) {
	s.commands = append(s.commands, c)
	f := c.FirstSample()
	if s.commandStarts[f] == nil {
		s.commandStarts[f] = make([]Command, 0)
	}
	s.commandStarts[f] = append(s.commandStarts[f], c)
	l := c.LastSample()
	if s.commandEnds[l] == nil {
		s.commandEnds[l] = make([]Command, 0)
	}
	s.commandEnds[l] = append(s.commandEnds[l], c)
}

func (s *Synthesizer) NumSamples() int {
	max := 0
	for _, c := range s.commands {
		if c.LastSample() > max {
			max = c.LastSample()
		}
	}
	return max
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
	buf := make([][]float32, 1)
	buf[0] = make([]float32, numSamples)
	writeBuf := make([]int16, numSamples)
	go func() {
		s.SampleBuffer(buf)
	}()
	<-s.notifications
	for i := 0; i < numSamples; i++ {
		writeBuf[i] = int16(buf[0][i] * 0x7FFF)
	}
	f.WriteItems(writeBuf)
}

func (s *Synthesizer) Play() {
	portaudio.Initialize()
	defer portaudio.Terminate()
	stream, err := portaudio.OpenDefaultStream(0, 1, float64(s.SampleRate), 0, s.SampleBuffer)
	if err != nil {
		panic(err)
	}
	s.curSample = 0
	stream.Start()
	<-s.notifications
	stream.Stop()
}

func (s *Synthesizer) activateCommand(c Command) {
	s.activeCommands = append(s.activeCommands, c)
}

func (s *Synthesizer) activateUpcoming() {
	upcoming := s.commandStarts[s.curSample]
	if upcoming == nil {
		return
	}
	for _, c := range upcoming {
		s.activateCommand(c)
	}
}

func (s *Synthesizer) deactivateCommand(c Command) {
	idx := -1
	for i, a := range s.activeCommands {
		if a == c {
			idx = i
			break
		}
	}
	if idx < 0 {
		return
	}
	s.activeCommands = append(s.activeCommands[:idx], s.activeCommands[idx+1:]...)
}

func (s *Synthesizer) deactivateCompleted() {
	completed := s.commandEnds[s.curSample]
	if completed == nil {
		return
	}
	for _, c := range completed {
		s.deactivateCommand(c)
	}
}

func (s *Synthesizer) SampleBuffer(out [][]float32) {
	for i := range out[0] {
		var val float32 = 0
		s.activateUpcoming()
		for _, c := range s.activeCommands {
			val += float32(c.SampleFrame(s.curSample))
		}
		s.deactivateCompleted()
		if val > 1.0 {
			val = 1.0
		}
		out[0][i] = val
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
	for _, c := range g.Commands() {
		s.AddCommand(c)
	}
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
