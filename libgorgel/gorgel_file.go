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
		envelope := ENVELOPE_RECTANGULAR
		if numParams > 3 {
			envelopeName = params[3]
		}
		switch envelopeName {
		case "lin":
			envelope = ENVELOPE_LINEAR
		case "rect":
			envelope = ENVELOPE_RECTANGULAR
		}
		return NewCmdNote(int(idx), int(begin), int(duration), envelope, g.synthesizer), nil
	}
	return nil, errors.New("invalid command")
}

func (g *GorgelFile) process(r *bufio.Reader) {
	eof := false
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
		c, err := g.createCommand(cmdType, toks)
		if err == nil {
			g.commands = append(g.commands, c)
		}
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
