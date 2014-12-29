package libgorgel

import (
	"testing"
)

func TestRead(t *testing.T) {
	g := NewGorgelFile("testdata/popcorn.gorgel", nil)
	err := g.Read()
	if err != nil {
		t.Errorf("error reading GorgelFile during test: %s", err)
	}
	n := len(g.Commands())
	if n != 29 {
		t.Errorf("expected: 29 commands, actual: %d", n)
	}
	c1 := g.Commands()[0]
	if c1.AsString() != "N 40, 0, 2" {
		t.Errorf("expected commands[0] to be \"N 40, 0, 2\", not \"%s\"", c1.AsString())
	}
	c2 := g.Commands()[28]
	if c2.AsString() != "N 43, 90, 2" {
		t.Errorf("expected commands[28] to be \"N 43, 90, 2\", not \"%s\"", c2.AsString())
	}
}
