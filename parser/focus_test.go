package parser_test

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	. "github.com/ionous/iffy/parser"
	"github.com/ionous/sliceOf"
)

var dropGrammar = allOf(Words("drop"), anyOf(
	allOf(&Focus{Where: "held", What: things()}, &Action{"Drop"}),
))

type MyContext struct {
	MyBounds // world
	Player   map[string]Bounds
	Other    map[ident.Id]Bounds
	Log
}

func (m MyContext) GetPlayerBounds(n string) (ret Bounds, err error) {
	if s, ok := m.Player[n]; ok {
		m.Log.Log("asking for bounds", n, len(s.(MyBounds)))
		ret = s
	} else {
		ret = m
	}
	return
}

func (m MyContext) GetObjectBounds(n ident.Id) (ret Bounds, err error) {
	if s, ok := m.Other[n]; ok {
		m.Log.Log("asking for bounds", n, len(s.(MyBounds)))
		ret = s
	} else {
		err = errutil.New("unknown bounds", n)
	}
	return
}

func TestFocus(t *testing.T) {
	grammar := dropGrammar
	bounds := MyBounds{
		makeObject("red apple", "apples"),
		makeObject("apple cart", "carts"),
		makeObject("red cart", "carts"),
	}
	invBounds := MyBounds{
		makeObject("torch", "devices"),
		makeObject("crab apple", "apples"),
	}
	ctx := MyContext{
		Log:      t,
		MyBounds: bounds,
		Player:   map[string]Bounds{"held": invBounds},
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
