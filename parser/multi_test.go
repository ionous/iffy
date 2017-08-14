package parser_test

import (
	"github.com/ionous/sliceOf"
	"strings"
	"testing"
)

func TestMulti(t *testing.T) {
	grammar := pickGrammar
	pickup := func(which string) []string {
		return sliceOf.String(
			strings.Join(sliceOf.String("pick", "up", which), " "),
			strings.Join(sliceOf.String("pick", which, "up"), " "),
		)
	}
	t.Run("all", func(t *testing.T) {
		e := parse(ctx, grammar,
			pickup("all"),
			&ActionGoal{"Take", sliceOf.String(
				"something",
				"red-apple",
				"crab-apple",
				"apple-cart",
				"red-cart",
				"apple",
				"torch")})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("some", func(t *testing.T) {
		e := parse(ctx, grammar,
			pickup("all red"),
			&ActionGoal{"Take", sliceOf.String(
				"red-apple",
				"red-cart")})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("plurals", func(t *testing.T) {
		e := parse(ctx, grammar,
			sliceOf.String("pick up apples"),
			&ActionGoal{"Take", sliceOf.String(
				"apple",
				"red-apple",
				"crab-apple")})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("one plural", func(t *testing.T) {
		e := parse(ctx, grammar,
			sliceOf.String("pick up red apples", "pick up apples red"),
			&ActionGoal{"Take", sliceOf.String(
				"red-apple")})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("empty plural", func(t *testing.T) {
		e := parse(ctx, grammar,
			sliceOf.String("pick up red apple carts"),
			&ActionGoal{"Take", sliceOf.String(
				"red-apple")})
		if e != nil {
			t.Fatal(e)
		}
	})
}
