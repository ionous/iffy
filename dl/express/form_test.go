package express

import (
	testify "github.com/stretchr/testify/assert"
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

	t.Run("parts", func(t *testing.T) {
		x := Template{
			{0, "status.score", false},
			{14, "/", true},
			{15, "story.turnCount", false},
		}
		if res, ok := Tokenize(partStr); testify.True(t, ok) {
			testify.Equal(t, x, res)
		}
	})
	t.Run("cmdStr", func(t *testing.T) {
		x := Template{
			{0, "TestThe example", false},
		}
		if res, ok := Tokenize(cmdStr); testify.True(t, ok) {
			testify.Equal(t, x, res)
		}
	})
	t.Run("empty", func(t *testing.T) {
		x := Template{
			{0, "its ", true},
			{4, "", false},
			{6, " empty", true},
		}
		if res, ok := Tokenize(emptyStr); testify.True(t, ok) {
			testify.Equal(t, x, res)
		}
	})
	t.Run("none", func(t *testing.T) {
		if _, ok := Tokenize(noneStr); ok {
			t.Fatal("should be false")
		}
	})
	t.Run("nobrackets", func(t *testing.T) {
		if _, ok := Tokenize(nobracketStr); ok {
			t.Fatal("should be false")
		}
	})
	t.Run("escape", func(t *testing.T) {
		if _, ok := Tokenize(escapeStr); ok {
			t.Fatal("should be false")
		}
	})
	t.Run("else", func(t *testing.T) {
		x := Template{
			{0, "if x", false},
			{6, "status.score", false},
			{20, "else", false},
			{26, "story.turnCount", false},
			{43, "endif", false},
		}
		if res, ok := Tokenize(ifElseStr); testify.True(t, ok) {
			testify.Equal(t, x, res)
		}
	})
}
