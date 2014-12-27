package main

import (
	"fmt"
	"github.com/mkb218/gosndfile/sndfile"
	"math"
	"os"
)

func sampleNotes(freqs *[]float32, amplitude float32, pos int) int16 {
	res := 0.0
	for _, freq := range *freqs {
		res += float64(amplitude) * math.Sin(float64(pos) / float64(freq)) / 2.0
	}
	return int16(res)
}

func main() {
	fmt.Printf("Welcome to Gorgel\n")
	var info sndfile.Info
	info.Channels = 1
	info.Samplerate = 44100
	info.Format = sndfile.SF_FORMAT_WAV|sndfile.SF_FORMAT_PCM_16
	seq := []int{16, 14, 16, 11, 7, 11, 4, 0,
				 16, 14, 16, 11, 7, 11, 4, 0}
	iterations := int(info.Samplerate / 7)
	info.Frames = int64(len(seq) * iterations)
	
	outfilename := "output.wav"
	if (len(os.Args) > 1) {
		outfilename = os.Args[1]
	}
	f, err := sndfile.Open(outfilename, sndfile.Write, &info)
	defer f.Close()
	if err != nil {
		fmt.Errorf("Error opening %s\n", os.Args[0])
		os.Exit(-1)
	}
	
	a440hz := float32(info.Samplerate) / math.Pi / 2.0 / 440
	notes := new([37]float32)
	step := math.Pow(2, 1./12)
	notes[0] = 99999
	for i := 1; i < 37; i++ {
		notes[i] = a440hz / float32(math.Pow(step, float64(i)))
	}
	
	for _, idx := range seq {
		snotes := make([]float32, 1)
		snotes[0] = notes[idx]
		for j := 0; j < iterations; j++ {
			val := sampleNotes(&snotes, 16000, j)
			f.WriteItems([]int16{val})
		}
	}
}
