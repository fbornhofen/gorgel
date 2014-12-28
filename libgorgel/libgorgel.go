package libgorgel

import (
	"fmt"
	"github.com/mkb218/gosndfile/sndfile"
	"math"
	"os"
)

type Gorgel struct {
	Samplerate int
	Channels int
	SamplesPerUnit int
	format int
	headerWritten bool
	file *sndfile.File
	
	notes []float32
	amplitude float32
}

func NewGorgel(samplerate int, channels int, samplesPerUnit int) *Gorgel {
	res := new(Gorgel)
	res.Samplerate = samplerate
	res.Channels = channels
	res.SamplesPerUnit = samplesPerUnit
	res.amplitude = 16000
	res.init()
	return res
}

func (g *Gorgel) init() {
	a440hz := float32(g.Samplerate) / math.Pi / 2.0 / 440
	g.notes = make([]float32, 37)
	step := math.Pow(2, 1./12)
	g.notes[0] = 99999
	for i := 1; i < 37; i++ {
		g.notes[i] = a440hz / float32(math.Pow(step, float64(i)))
	}
}

func (g *Gorgel) openAndWriteHeader(filename string) {
	var info sndfile.Info
	info.Channels = int32(g.Channels)
	info.Samplerate = int32(g.Samplerate)
	info.Format = sndfile.SF_FORMAT_WAV|sndfile.SF_FORMAT_PCM_16
	var err error
	g.file, err = sndfile.Open(filename, sndfile.Write, &info)
	if err != nil {
		fmt.Errorf("Error opening %s\n", filename)
		os.Exit(-1)
	}
}

func sampleNotes(freqs *[]float32, amplitude float32, pos int) int16 {
	res := 0.0
	for _, freq := range *freqs {
		res += float64(amplitude) * math.Sin(float64(pos) / float64(freq)) / 2.0
	}
	return int16(res)
}

func (g *Gorgel) Synthesize(seq *[]int, filename string) {
	g.openAndWriteHeader(filename)
	for _, idx := range *seq {
		snotes := make([]float32, 1)
		snotes[0] = g.notes[idx]
		for j := 0; j < g.SamplesPerUnit; j++ {
			val := sampleNotes(&snotes, g.amplitude, j)
			g.file.WriteItems([]int16{val})
		}
	}
}
