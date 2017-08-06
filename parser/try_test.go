package parser_test

import (
	. "github.com/ionous/iffy/parser"
	testify "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestWord(t *testing.T) {
	match := func(input, goal string) bool {
		ctx := Context{
			Words: strings.Fields(input),
			Match: &Match{Scanner: &Word{goal}},
		}
		return ctx.Advance()
	}
	t.Run("match", func(t *testing.T) {
		testify.True(t, match("Beep", "beep"))
	})
	t.Run("mismatch", func(t *testing.T) {
		testify.False(t, match("Boop", "beep"))
	})
}

func TestVariousOf(t *testing.T) {
	words := func(goal ...string) (ret []Scanner) {
		for _, g := range goal {
			ret = append(ret, &Word{g})
		}
		return
	}

	match := func(input string, goal Scanner) bool {
		ctx := Context{
			Words: strings.Fields(input),
			Match: &Match{Scanner: goal},
		}
		return ctx.Advance()
	}
	t.Run("any", func(t *testing.T) {
		wordList := words("beep", "blox")
		t.Run("match", func(t *testing.T) {
			testify.True(t, match("Beep", &AnyOf{wordList}))
			testify.True(t, match("Blox", &AnyOf{wordList}))
		})
		t.Run("mismatch", func(t *testing.T) {
			testify.False(t, match("Boop", &AnyOf{wordList}))
			testify.False(t, match("Beep", &AnyOf{}))
			testify.False(t, match("", &AnyOf{}))
		})
	})
	t.Run("all", func(t *testing.T) {
		wordList := words("beep", "blox")
		t.Run("match", func(t *testing.T) {
			testify.True(t, match("Beep BLOX", &AllOf{wordList}))
		})
		t.Run("mismatch", func(t *testing.T) {
			testify.False(t, match("BLOX Beep", &AllOf{wordList}))
			testify.False(t, match("Beep", &AllOf{wordList}))
			testify.False(t, match("BLOX", &AllOf{wordList}))
			testify.False(t, match("Boop", &AllOf{wordList}))
			testify.False(t, match("", &AllOf{wordList}))
			testify.False(t, match("", &AllOf{}))
		})
	})
}

// 	func TestAllOf(t *testing.T) {
// 	}
