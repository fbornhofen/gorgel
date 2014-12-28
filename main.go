package main

import (
	"fmt"
	"gorgel/libgorgel"
	"os"
)



func main() {
	fmt.Printf("Welcome to Gorgel\n")
	g := libgorgel.NewGorgel(44100, 1, 44100/7)
	seq := []int{16, 14, 16, 11, 7, 11, 4, 0,
				 16, 14, 16, 11, 7, 11, 4, 0,
				 16, 18, 19, 18, 19, 16, 18, 16,
				 18, 14, 16, 14, 16, 11, 16}
	
	outfilename := "output.wav"
	if (len(os.Args) > 1) {
		outfilename = os.Args[1]
	}
	
	g.Synthesize(&seq, outfilename)
}
