package parser_test

import (
	"github.com/ionous/iffy/ident"
	. "github.com/ionous/iffy/parser"
	"github.com/ionous/sliceOf"
	"testing"
)

// note, in reality burn would use only held things.
var takeGrammar = allOf(
	Words("get"),
	&Target{[]Scanner{things(), Words("from/off"), thing()}},
	&Action{"Remove"},
)

func TestTarget(t *testing.T) {
	grammar := takeGrammar
	scope := MyScope{
		makeObject("green apple", "apples"),
		makeObject("apple cart", "carts"),
		makeObject("red cart", "carts"),
	}
	appleCart := MyScope{
		makeObject("crab apple", "apples"),
		makeObject("red apple", "apples"),
	}
	redCart := MyScope{
		makeObject("yellow apple", "apples"),
	}
	//
	ctx := MyContext{
		Log:     t,
		MyScope: scope,
		Other: Scopes{
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
