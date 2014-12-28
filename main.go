package main

import (
	"fmt"
	"gorgel/libgorgel"
	"os"
)



func main() {
	fmt.Printf("Welcome to Gorgel\n")
	g := libgorgel.NewGorgel(44100, 1, 44100/8)
	seq := []int{40, 38, 40, 35, 31, 35, 28, 0,
		         40, 38, 40, 35, 31, 35, 28, 0,
		         40, 42, 43, 42, 43, 40, 42, 40,
		         42, 38, 40, 38, 40, 35, 40}
	
	outfilename := "output.wav"
	if (len(os.Args) > 1) {
		outfilename = os.Args[1]
	}
	
	g.Synthesize(&seq, outfilename)
}
