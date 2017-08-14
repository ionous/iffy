package parser_test

import (
	. "github.com/ionous/iffy/parser"
	"github.com/ionous/sliceOf"
	"testing"
)

// note, in reality burn would use only held things.
var burnGrammar = allOf(Words("burn/light"), anyOf(
	allOf(things(), &Action{"Burn"}),
	allOf(things(), Words("with"), things(), &Action{"Burn"}),
))

func TestDual(t *testing.T) {
	grammar := burnGrammar
	t.Run("one", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("burn/light cart"),
			&ClarifyGoal{"red"},
			&ActionGoal{"Burn", sliceOf.String("red-cart")})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("two", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			sliceOf.String("burn red cart with torch"),
			&ActionGoal{"Burn", sliceOf.String("red-cart", "torch")})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("clarify", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("light cart with"),
			&ClarifyGoal{"red"},
			&ClarifyGoal{"torch"},
			&ActionGoal{"Burn", sliceOf.String("red-cart", "torch")})
		if e != nil {
			t.Fatal(e)
		}
	})
}
