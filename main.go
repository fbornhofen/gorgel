package main

import (
	"fmt"
	"gorgel/libgorgel"
	"os"
)

func main() {
	fmt.Printf("Welcome to Gorgel\n")
	infilename := ""
	if len(os.Args) > 1 {
		infilename = os.Args[1]
	} else {
		fmt.Errorf("usage: gorgel INFILE [OUTFILE]\n")
		os.Exit(-1)
	}
	s := libgorgel.NewSynthesizer(120, 44100)
	err := s.ReadFromFile(infilename)
	if err != nil {
		fmt.Errorf("%s", err)
		os.Exit(-1)
	}
	if len(os.Args) > 2 {
		outfilename := os.Args[2]
		fmt.Printf("Writing %s\n", outfilename)
		s.WriteWaveFile(outfilename)
	} else {
		fmt.Printf("Playing\n")
		s.Play()
	}
}
