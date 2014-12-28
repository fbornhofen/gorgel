package main

import (
	"fmt"
	"gorgel/libgorgel"
	"os"
)

func main() {
	fmt.Printf("Welcome to Gorgel\n")
	outfilename := "output.wav"
	if len(os.Args) > 1 {
		outfilename = os.Args[1]
	}
	fmt.Printf("Writing %s\n", outfilename)
	s := libgorgel.NewSynthesizer(120, 44100)
	seq := []*libgorgel.CmdNote{
		libgorgel.NewCmdNote(40, 0, 2, s),
		libgorgel.NewCmdNote(38, 3, 2, s),
		libgorgel.NewCmdNote(40, 6, 2, s),
		libgorgel.NewCmdNote(35, 9, 2, s),
		libgorgel.NewCmdNote(31, 12, 2, s),
		libgorgel.NewCmdNote(35, 15, 2, s),
		libgorgel.NewCmdNote(28, 18, 2, s),

		libgorgel.NewCmdNote(40, 24, 2, s),
		libgorgel.NewCmdNote(38, 27, 2, s),
		libgorgel.NewCmdNote(40, 30, 2, s),
		libgorgel.NewCmdNote(35, 33, 2, s),
		libgorgel.NewCmdNote(31, 36, 2, s),
		libgorgel.NewCmdNote(35, 39, 2, s),
		libgorgel.NewCmdNote(28, 42, 2, s),

		libgorgel.NewCmdNote(40, 48, 2, s),
		libgorgel.NewCmdNote(42, 51, 2, s),
		libgorgel.NewCmdNote(43, 54, 2, s),
		libgorgel.NewCmdNote(42, 57, 2, s),
		libgorgel.NewCmdNote(43, 60, 2, s),
		libgorgel.NewCmdNote(40, 63, 2, s),
		libgorgel.NewCmdNote(42, 66, 2, s),
		libgorgel.NewCmdNote(40, 69, 2, s),
		libgorgel.NewCmdNote(42, 72, 2, s),
		libgorgel.NewCmdNote(38, 75, 2, s),
		libgorgel.NewCmdNote(40, 78, 2, s),
		libgorgel.NewCmdNote(38, 81, 2, s),
		libgorgel.NewCmdNote(40, 84, 2, s),
		libgorgel.NewCmdNote(42, 87, 2, s),
		libgorgel.NewCmdNote(43, 90, 2, s),
	}
	for _, n := range seq {
		s.AddCommand(n)
	}
	s.WriteWaveFile(outfilename)
}
