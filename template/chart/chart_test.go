package chart

import (
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestChart(t *testing.T) {
	t.Run("directives", func(t *testing.T) {
		assert, x := testify.New(t), true
		// parse a simple string into a quoted expression.
		x = x && assert.NoError(testChart(t, "hello world",
			"`hello world`"))

		// our tricky case: a single word is a key.
		x = x && assert.Error(testChart(t, "hello {player}",
			"hello {player:}"))

		x = x && assert.NoError(testChart(t, `{player!}`,
			"PLAYER/0"))

		// a function and some text should be a single span
		x = x && assert.NoError(testChart(t, "{player!}, hello",
			"PLAYER/0 `, hello` Span/2"))

		// text and functions with directives:
		x = x && assert.NoError(testChart(t,
			`hello {player!} to the {"world"|cap?}.`,
			"`hello ` PLAYER/0 ` to the ` `world` CAP/1 `.` Span/5"))
	})

	t.Run("ifs", func(t *testing.T) {
		assert, x := testify.New(t), true
		// an empty if statement
		x = x && assert.NoError(testChart(t, "{if test?}{end}",
			"TEST/0 IfStatement/1"))

		// a more complex if statement
		x = x && assert.NoError(testChart(t, "{if 5=4}{end}",
			"5 4 EQL IfStatement/1"))

		// an if statement with content
		x = x && assert.NoError(testChart(t, "{if boop?}beep{end}",
			"BOOP/0 `beep` IfStatement/2"))

		// if statement with content, leading and trailing text
		x = x && assert.NoError(testChart(t, "abc{if boop?}beep{end}cba",
			"`abc` BOOP/0 `beep` IfStatement/2 `cba` Span/3"))

		x = x && assert.NoError(testChart(t, "{if test?}hello {player!}{end}",
			"TEST/0 `hello ` PLAYER/0 Span/2 IfStatement/2"))

		// if-else
		x = x && assert.NoError(testChart(t, "{if test?}hello {player!}{else}bellow{end}",
			"TEST/0 `hello ` PLAYER/0 Span/2 `bellow` IfStatement/3"))

		// if-else with empty clauses
		x = x && assert.NoError(testChart(t, "{if boo?}{else}{end}",
			"BOO/0 Span/0 Span/0 IfStatement/3"))

		// test unless-otherwise
		x = x && assert.NoError(testChart(t, "{unless boo?}blix{otherwise}blox{end}",
			"BOO/0 `blix` `blox` UnlessStatement/3"))

		// test if-elseIf chains
		x = x && assert.NoError(testChart(t, "{if boo?}a{elseIf beep?}b{end}",
			"BOO/0 `a` BEEP/0 `b` IfStatement/2 IfStatement/3"))

		// test if-elseIf with empty leading if
		x = x && assert.NoError(testChart(t, "{if boo?}{elseIf beep?}b{end}",
			"BOO/0 Span/0 BEEP/0 `b` IfStatement/2 IfStatement/3"))

		// test no end; this actually works :(
		// x = x && assert.Error(testChart(t, `{if test}`, ignoreResult))
		x = x && assert.Error(testChart(t, "{end}", ignoreResult))
	})

	t.Run("sequence", func(t *testing.T) {
		assert, x := testify.New(t), true
		x = x && assert.NoError(testChart(t, "{cycle}a{end}",
			"`a` Cycle/1"))
		x = x && assert.NoError(testChart(t, "{cycle}a{or}b{or}c{end}",
			"`a` `b` `c` Cycle/3"))
		x = x && assert.NoError(testChart(t, "{cycle}a{or}{player!}{or}c{end}",
			"`a` PLAYER/0 `c` Cycle/3"))
		x = x && assert.NoError(testChart(t, "x{cycle}a{or}p{player!}q{end}y",
			"`x` `a` `p` PLAYER/0 `q` Span/3 Cycle/2 `y` Span/3"))

		// test if statement within a cycle.
		x = x && assert.NoError(testChart(t, "{cycle}a{or}{if boop?}beep{end}{end}",
			"`a` BOOP/0 `beep` IfStatement/2 Cycle/2"))
		// test mismatched keywords
		x = x && assert.Error(testChart(t, "{cycle}a{or}{if boop?}{or}beep{end}{end}",
			ignoreResult))
	})
}

func testChart(t *testing.T, str, want string) (err error) {
	t.Log("test:", str)
	if ds, e := Parse(str); e != nil {
		err = e
	} else if got := ds.String(); want == ignoreResult {
		t.Log("got", got)
	} else if got != want {
		err = mismatched(want, got)
	} else {
		t.Log("ok", got)
	}
	return
}
