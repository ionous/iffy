package parser_test

import (
	"testing"

	"github.com/ionous/iffy/ident"
	. "github.com/ionous/iffy/parser"
	"github.com/ionous/sliceOf"
)

var takeGrammar = allOf(
	Words("get"),
	&Target{[]Scanner{things(), Words("from/off"), thing()}},
	&Action{"Remove"},
)

func TestTarget(t *testing.T) {
	grammar := takeGrammar
	bounds := MyBounds{
		makeObject("green apple", "apples"),
		makeObject("apple cart", "carts"),
		makeObject("red cart", "carts"),
	}
	appleCart := MyBounds{
		makeObject("crab apple", "apples"),
		makeObject("red apple", "apples"),
	}
	redCart := MyBounds{
		makeObject("yellow apple", "apples"),
	}
	//
	ctx := MyContext{
		Log:      t,
		MyBounds: bounds,
		Other: map[ident.Id]Bounds{
			ident.IdOf("apple-cart"): appleCart,
			ident.IdOf("red-cart"):   redCart},
	}

	t.Run("take exact", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("get apple from/off red cart"),
			&ActionGoal{"Remove", sliceOf.String("red-cart", "yellow-apple")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("clarify", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("get apple from/off cart"),
			&ClarifyGoal{"apple"},
			&ClarifyGoal{"red"},
			&ActionGoal{"Remove", sliceOf.String("apple-cart", "red-apple")})
		if e != nil {
			t.Fatal(e)
		}
	})
}
