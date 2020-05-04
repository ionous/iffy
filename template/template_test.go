package template

import (
	"testing"

	"github.com/ionous/errutil"
)

func TestFailures(t *testing.T) {
	if e := testChart(t, "{go testScore Story.score}",
		ignoreResult); e == nil {
		t.Fatal(e)
	}
	if e := testChart(t, "{Story.score|testScore}",
		ignoreResult); e == nil {
		t.Fatal(e)
	}
}

func TestChart(t *testing.T) {
	t.Run("directives", func(t *testing.T) {
		// parse a simple string into a quoted expression.
		if e := testChart(t, "hello world",
			`"hello world"`); e != nil {
			t.Fatal(e)
		}
		// our tricky case: a single word is a key.
		if e := testChart(t, "hello {player}",
			`hello {player:}`); e == nil {
			t.Fatal(e)
		}
		if e := testChart(t, `{player!}`,
			`PLAYER/0`); e != nil {
			t.Fatal(e)
		}
		// a function and some text should be a single span
		if e := testChart(t, "{player!}, hello",
			`PLAYER/0 ", hello" Span/2`); e != nil {
			t.Fatal(e)
		}
		// text and functions with directives:
		if e := testChart(t,
			`hello {player!} to the {"world"|cap?}.`,
			`"hello " PLAYER/0 " to the " "world" CAP/1 "." Span/5`); e != nil {
			t.Fatal(e)
		}
	})

	t.Run("ifs", func(t *testing.T) {
		// an empty if statement
		if e := testChart(t, "{if test?}{end}",
			`TEST/0 IfStatement/1`); e != nil {
			t.Fatal(e)
		}
		// a more complex if statement
		if e := testChart(t, "{if 5=4}{end}",
			`5 4 EQL IfStatement/1`); e != nil {
			t.Fatal(e)
		}
		// an if statement with content
		if e := testChart(t, "{if boop?}beep{end}",
			`BOOP/0 "beep" IfStatement/2`); e != nil {
			t.Fatal(e)
		}
		// if statement with content, leading and trailing text
		if e := testChart(t, "abc{if boop?}beep{end}cba",
			`"abc" BOOP/0 "beep" IfStatement/2 "cba" Span/3`); e != nil {
			t.Fatal(e)
		}
		if e := testChart(t, "{if test?}hello {player!}{end}",
			`TEST/0 "hello " PLAYER/0 Span/2 IfStatement/2`); e != nil {
			t.Fatal(e)
		}
		// if-else
		if e := testChart(t, "{if test?}hello {player!}{else}bellow{end}",
			`TEST/0 "hello " PLAYER/0 Span/2 "bellow" IfStatement/3`); e != nil {
			t.Fatal(e)
		}
		// if-else with empty clauses
		if e := testChart(t, "{if boo?}{else}{end}",
			`BOO/0 Span/0 Span/0 IfStatement/3`); e != nil {
			t.Fatal(e)
		}
		// test unless-otherwise
		if e := testChart(t, "{unless boo?}blix{otherwise}blox{end}",
			`BOO/0 "blix" "blox" UnlessStatement/3`); e != nil {
			t.Fatal(e)
		}
		// test if-elsif chains
		if e := testChart(t, "{if boo?}a{elsif beep?}b{end}",
			`BOO/0 "a" BEEP/0 "b" IfStatement/2 IfStatement/3`); e != nil {
			t.Fatal(e)
		}
		// test if-elsif with empty leading if
		if e := testChart(t, "{if boo?}{elsif beep?}b{end}",
			`BOO/0 Span/0 BEEP/0 "b" IfStatement/2 IfStatement/3`); e != nil {
			t.Fatal(e)
		}
		// test no end; this actually works :(
		// if e:= testChart(t, "{if test}", ignoreResult); e== nil { t.Fatal(e) }
		if e := testChart(t, "{end}", ignoreResult); e == nil {
			t.Fatal(e)
		}
	})

	t.Run("sequence", func(t *testing.T) {
		if e := testChart(t, "{cycle}a{end}",
			`"a" Cycle/1`); e != nil {
			t.Fatal(e)
		}
		if e := testChart(t, "{cycle}a{or}b{or}c{end}",
			`"a" "b" "c" Cycle/3`); e != nil {
			t.Fatal(e)
		}
		if e := testChart(t, "{cycle}a{or}{player!}{or}c{end}",
			`"a" PLAYER/0 "c" Cycle/3`); e != nil {
			t.Fatal(e)
		}
		if e := testChart(t, "x{cycle}a{or}p{player!}q{end}y",
			`"x" "a" "p" PLAYER/0 "q" Span/3 Cycle/2 "y" Span/3`); e != nil {
			t.Fatal(e)
		}
		// test if statement within a cycle.
		if e := testChart(t, "{cycle}a{or}{if boop?}beep{end}{end}",
			`"a" BOOP/0 "beep" IfStatement/2 Cycle/2`); e != nil {
			t.Fatal(e)
		}
		// test mismatched keywords
		if e := testChart(t, "{cycle}a{or}{if boop?}{or}beep{end}{end}",
			ignoreResult); e == nil {
			t.Fatal(e)
		}
	})
}

func testChart(t *testing.T, str, want string) (err error) {
	t.Log("test:", str)
	if ds, e := Parse(str); e != nil {
		err = e
	} else if got := ds.String(); want == ignoreResult || got == want {
		t.Log("got", got)
	} else {
		err = mismatched(want, got)
	}
	return
}

func mismatched(want, got string) error {
	return errutil.Fmt("want(%d): %s; != got(%d): %s.", len(want), want, len(got), got)
}

const ignoreResult = "~~IGNORE~~"
