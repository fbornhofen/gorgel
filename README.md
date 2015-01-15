gorgel
======

Synthesizer written in Go. The idea is to learn Go and sound synthesis at the same time.

Generating WAVE Files
---------------------

Usage:

    gorgel infile.gorgel outfile.wav

See libgorgel/testdata for sample input files. At the moment, the only supported output format is WAV with one channel and 16 bit sample width. However, libsndfile has more to offer, so this limitation may be temporary.

Real-time Synthesis
-------------------

Usage:

    gorgel infile.gorgel
    
When no output file is specified, gorgel just plays the file using PortAudio.

Ideas/Plans
-----------

- More expressive input format
- More options and effects in the synthesizer
- Live synthesis
- (G)UI for live synthesis
- Improve performance: low-hanging fruit, but it does not hurt enough yet ...
