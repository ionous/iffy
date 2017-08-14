package parser_test

import (
	"github.com/ionous/sliceOf"
	"testing"
)

func TestDisambiguation(t *testing.T) {
	grammar := lookGrammar
	t.Run("trailing noun", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("look at"),
			&ClarifyGoal{"something"},
			&ActionGoal{"Examine", sliceOf.String("something")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("shared names", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("look at red cart"),
			&ActionGoal{"Examine", sliceOf.String("red-cart")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("ambiguous shared names", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("look at cart"),
			&ClarifyGoal{"red"},
			&ActionGoal{"Examine", sliceOf.String("red-cart")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("ambiguous loops", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("look at cart"),
			&ClarifyGoal{"cart"},
			&ClarifyGoal{"cart"},
			&ClarifyGoal{"red"},
			&ActionGoal{"Examine", sliceOf.String("red-cart")})
		if e != nil {
			t.Fatal(e)
		}
	})

	// FIX? names which are "subsets" of other names, dont play well in inform
	// nor do they here. might consider adding tests for that in "compliation"
	// even though it doesn't during normal play.
	// t.Run("exact name works during disambiguation", func(t *testing.T) {
	// 	e := parse(t, ctx, grammar,
	// 		Phrases("look at apple"),
	// 		&ClarifyGoal{"apple"},
	// 		&ActionGoal{"Examine", sliceOf.String("apple")})
	// 	if e != nil {
	// 		t.Fatal(e)
	// 	}
	// })

	t.Run("doubled names dont match incorrectly", func(t *testing.T) {
		e := parse(t, ctx, grammar,
			Phrases("look at apple apple apple cart"),
			&ActionGoal{"Examine", sliceOf.String("apple-cart")})
		if e != nil {
			t.Fatal(e)
		}
	})
}
