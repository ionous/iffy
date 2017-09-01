package parser_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	. "github.com/ionous/iffy/parser"
	"github.com/ionous/sliceOf"
	"testing"
)

var dropGrammar = allOf(Words("drop"), anyOf(
	allOf(&Focus{Where: ident.IdOf("held"), What: things()}, &Action{"Drop"}),
))

type MyContext struct {
	MyScope       // world
	Player, Other Scopes
	Log
}

type Scopes map[ident.Id]MyScope

func (m MyContext) GetPlayerScope(n ident.Id) (ret Scope, err error) {
	if s, ok := m.Player[n]; ok {
		m.Log.Log("asking for scope", n, len(s))
		ret = s
	} else {
		ret = m
	}
	return
}

func (m MyContext) GetOtherScope(n ident.Id) (ret Scope, err error) {
	if s, ok := m.Other[n]; ok {
		m.Log.Log("asking for scope", n, len(s))
		ret = s
	} else {
		err = errutil.New("unknown scope", n)
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
	ctx := MyContext{
		Log:     t,
		MyScope: scope,
		Player:  Scopes{ident.IdOf("held"): invScope},
	}

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
