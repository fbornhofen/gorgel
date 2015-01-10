package libgorgel

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type GorgelFile struct {
	filename    string
	commands    []Command
	synthesizer *Synthesizer
}

func NewGorgelFile(filename string, s *Synthesizer) *GorgelFile {
	f := new(GorgelFile)
	f.filename = filename
	f.synthesizer = s
	return f
}

func (g *GorgelFile) createCommand(t rune, params []string) (Command, error) {
	numParams := len(params)
	switch t {
	case 'N':
		idx, _ := strconv.ParseInt(params[0], 10, 32)
		begin, _ := strconv.ParseInt(params[1], 10, 32)
		duration, _ := strconv.ParseInt(params[2], 10, 32)
		var envelopeName string
		envelope := ENVELOPE_DEFAULT
		if numParams > 3 {
			envelopeName = params[3]
		}
		switch envelopeName {
		case "lin":
			envelope = ENVELOPE_LINEAR
		case "rect":
			envelope = ENVELOPE_RECTANGULAR
		case "adsr":
			envelope = ENVELOPE_ADSR
		}
		return NewCmdNote(int(idx), int(begin), int(duration), envelope, g.synthesizer), nil
	}
	return nil, errors.New("invalid command")
}

func (g *GorgelFile) applyHeaderCommand(params []string) {
	if len(params) < 2 {
		return
	}
	switch params[0] {
	case "bpm":
		bpm, _ := strconv.ParseInt(params[1], 10, 32)
		g.synthesizer.BeatsPerMin = int(bpm)
	case "envelope":
		envType := params[1]
		switch envType {
		case "rect":
			g.synthesizer.defaultEnvelope = ENVELOPE_RECTANGULAR
		case "lin":
			g.synthesizer.defaultEnvelope = ENVELOPE_LINEAR
		case "adsr":
			g.synthesizer.defaultEnvelope = ENVELOPE_ADSR
		}
	}
}

func (g *GorgelFile) process(r *bufio.Reader) {
	eof := false
	isHeader := true
	lineCount := 0
	for eof == false {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				eof = true
			} else {
				fmt.Errorf("%s\n", err)
				os.Exit(-1)
			}
		}
		if len(line) == 0 {
			continue
		}
		var cmdType rune = rune(line[0])
		toks := strings.Split(line[1:], ",")
		for i, t := range toks {
			toks[i] = strings.TrimSpace(t)
		}
		if cmdType == 'H' {
			if isHeader {
				g.applyHeaderCommand(toks)
			} else {
				fmt.Errorf("line %d: warning: header command not in file header, ignoring")
			}
		} else {
			isHeader = false
			c, err := g.createCommand(cmdType, toks)
			if err == nil {
				g.commands = append(g.commands, c)
			}
		}
		lineCount += 1
	}
}

func (g *GorgelFile) Read() error {
	f, err := os.Open(g.filename)
	if err != nil {
		fmt.Errorf("error opening %s: %s\n", g.filename, err)
		return err
	}
	r := bufio.NewReader(f)
	g.process(r)
	return nil
}

func (g *GorgelFile) Commands() []Command {
	return g.commands
}
