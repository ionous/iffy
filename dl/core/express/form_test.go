package express

import (
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestForm(t *testing.T) {
	const (
		partStr      = "{status.score}/{story.turnCount}"
		cmdStr       = "{TestThe example}"
		noneStr      = ""
		emptyStr     = "its {} empty"
		nobracketStr = "no quotes"
		escapeStr    = "its {{quoted"
		ifElseStr    = "{if x}{status.score}{else}{story.turnCount}{endif}"
	)

	t.Run("parts", func(t *testing.T) {
		x := []Token{
			{0, "status.score", false},
			{14, "/", true},
			{15, "story.turnCount", false},
		}
		res := Tokenize(partStr)
		testify.Equal(t, x, res)
	})
	t.Run("cmdStr", func(t *testing.T) {
		x := []Token{
			{0, "TestThe example", false},
		}
		res := Tokenize(cmdStr)
		testify.Equal(t, x, res)
	})
	t.Run("none", func(t *testing.T) {
		x := []Token(nil)
		res := Tokenize(noneStr)
		testify.Equal(t, x, res)
	})
	t.Run("empty", func(t *testing.T) {
		x := []Token{
			{0, "its ", true},
			{4, "", false},
			{6, " empty", true},
		}
		res := Tokenize(emptyStr)
		testify.Equal(t, x, res)
	})
	t.Run("nobrackets", func(t *testing.T) {
		res := Tokenize(nobracketStr)
		testify.Len(t, res, 1)
	})
	t.Run("escape", func(t *testing.T) {
		x := []Token{
			{0, escapeStr, true},
		}
		res := Tokenize(escapeStr)
		testify.Equal(t, x, res)
	})
	t.Run("else", func(t *testing.T) {
		x := []Token{
			{0, "if x", false},
			{6, "status.score", false},
			{20, "else", false},
			{26, "story.turnCount", false},
			{43, "endif", false},
		}
		res := Tokenize(ifElseStr)
		testify.Equal(t, x, res)
	})
}
