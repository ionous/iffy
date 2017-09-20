package template

import (
	"github.com/kr/pretty"
	"testing"
)

func TestForm(t *testing.T) {
	const (
		partStr      = "{status.score}/{story.turnCount}"
		cmdStr       = "{TestThe example}"
		emptyStr     = "its {} empty"
		noneStr      = ""
		nobracketStr = "no quotes"
		escapeStr    = "its {{quoted"
		ifElseStr    = "{if x}{status.score}{else}{story.turnCount}{endif}"
	)
	none := []Token(nil)
	t.Run("parts", func(t *testing.T) {
		x := []Token{
			{0, "status.score", false},
			{14, "/", true},
			{15, "story.turnCount", false},
		}
		ts := Tokenize(partStr)
		if d := pretty.Diff(x, ts); len(d) > 0 {
			t.Fatal(d)
		}
	})
	t.Run("cmdStr", func(t *testing.T) {
		x := []Token{
			{0, "TestThe example", false},
		}
		ts := Tokenize(cmdStr)
		if d := pretty.Diff(x, ts); len(d) > 0 {
			t.Fatal(d)
		}
	})
	t.Run("empty", func(t *testing.T) {
		x := []Token{
			{0, "its ", true},
			{4, "", false},
			{6, " empty", true},
		}
		ts := Tokenize(emptyStr)
		if d := pretty.Diff(x, ts); len(d) > 0 {
			t.Fatal(d)
		}
	})
	t.Run("none", func(t *testing.T) {
		ts := Tokenize(noneStr)
		if d := pretty.Diff(none, ts); len(d) > 0 {
			t.Fatal(d)
		}
	})
	t.Run("nobrackets", func(t *testing.T) {
		ts := Tokenize(nobracketStr)
		if d := pretty.Diff(none, ts); len(d) > 0 {
			t.Fatal(d)
		}
	})
	t.Run("escape", func(t *testing.T) {
		ts := Tokenize(escapeStr)
		if d := pretty.Diff(none, ts); len(d) > 0 {
			t.Fatal(d)
		}
	})
	t.Run("else", func(t *testing.T) {
		x := []Token{
			{0, "if x", false},
			{6, "status.score", false},
			{20, "else", false},
			{26, "story.turnCount", false},
			{43, "endif", false},
		}
		ts := Tokenize(ifElseStr)
		if d := pretty.Diff(x, ts); len(d) > 0 {
			t.Fatal(d)
		}
	})
}
