package parser_test

import (
	. "github.com/ionous/iffy/parser"
	"github.com/ionous/sliceOf"
	"testing"
)

var dropGrammar = allOf(Words("drop"), anyOf(
	allOf(&Focus{Where: "held", What: things()}, &Action{"Drop"}),
))

type MyContext struct {
	MyScope // world
	Inv     MyScope
}

func (m MyContext) GetPlayerScope(n string) (ret Scope) {
	switch n {
	case "held":
		ret = m.Inv
	default:
		ret = m
	}
	return
}

func TestFocus(t *testing.T) {

	grammar := dropGrammar

	scope := MyScope{
		makeObject("red apple", "apples"),
		makeObject("apple cart", "carts"),
		makeObject("red cart", "carts"),
	}
	invScope := MyScope{
		makeObject("torch", "devices"),
		makeObject("crab apple", "apples"),
	}
	ctx := MyContext{scope, invScope}

	t.Run("drop one", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("drop"),
			&ClarifyGoal{"apple"},
			&ActionGoal{"Drop", sliceOf.String("crab-apple")})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("drop all", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("drop everything"),
			&ActionGoal{"Drop", sliceOf.String("torch", "crab-apple")})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("drop error", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("drop cart"),
			&ErrorGoal{"you can't see any such thing"})
		if e != nil {
			t.Fatal("expected an error")
		}
	})
}
