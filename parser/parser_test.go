package parser_test

import (
	. "github.com/ionous/iffy/parser"
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func anyOf(s ...Scanner) (ret Scanner) {
	if len(s) == 1 {
		ret = s[0]
	} else {
		ret = &AnyOf{s}
	}
	return
}

func allOf(s ...Scanner) (ret Scanner) {
	if len(s) == 1 {
		ret = s[0]
	} else {
		ret = &AllOf{s}
	}
	return
}

func words(s string) (ret Scanner) {
	if split := strings.Split(s, "/"); len(split) == 1 {
		ret = &Word{s}
	} else {
		var words []Scanner
		for _, g := range split {
			words = append(words, &Word{g})
		}
		ret = &AnyOf{words}
	}
	return
}

func noun(f ...Filter) Scanner {
	return &Object{f}
}

func TestParser(t *testing.T) {
	obj := func(names ...string) MyObject {
		return MyObject{Id: names[0], Names: names}
	}
	scope := MyScope{
		obj("something"),
	}
	scope = append(scope, Directions()...)

	grammar :=
		allOf(words("look/l"), anyOf(
			allOf(&Action{"Look"}),
			allOf(words("at"), noun(), &Action{"Examine"}),
			allOf(words("inside/in/into/through/on"), noun(), &Action{"Search"}),
			allOf(words("under"), noun(), &Action{"LookUnder"}),
			allOf(noun(&HasClass{"direction"}), &Action{"Examine"}),
			allOf(words("to"), noun(&HasClass{"direction"}), &Action{"Examine"}),
		))

		// first, we want to test a simple set of example actions,
		// all of which start the same way, but end with different actions.
		// later, we will test disambiguation; errors; multiple objects: etc.
	t.Run("look", func(t *testing.T) {
		parse(t, scope, grammar,
			Phrases("look/l"),
			&Result{
				Action: "Look",
			})
	})
	t.Run("examine", func(t *testing.T) {
		parse(t, scope, grammar,
			Phrases("look/l at something"),
			&Result{
				Action: "Examine",
				Nouns:  sliceOf.String("something"),
			})
	})
	t.Run("search", func(t *testing.T) {
		parse(t, scope, grammar,
			Phrases("look/l inside/in/into/through/on something"),
			&Result{
				Action: "Search",
				Nouns:  sliceOf.String("something"),
			})
	})
	t.Run("look under", func(t *testing.T) {
		parse(t, scope, grammar,
			Phrases("look/l under something"),
			&Result{
				Action: "LookUnder",
				Nouns:  sliceOf.String("something"),
			})
	})
	t.Run("look dir", func(t *testing.T) {
		look := Phrases("look/l")
		for _, d := range directions {
			d := sliceOf.String(d)
			parse(t, scope, grammar,
				permute(look, d),
				&Result{
					Action: "Examine",
					Nouns:  d,
				})
		}
	})
	t.Run("look no dir", func(t *testing.T) {
		parse(t, scope, grammar,
			Phrases("look something"),
			nil)
	})
	t.Run("look to dir", func(t *testing.T) {
		lookTo := Phrases("look/l to")
		for _, d := range directions {
			d := sliceOf.String(d)
			parse(t, scope, grammar,
				permute(lookTo, d),
				&Result{
					Action: "Examine",
					Nouns:  d,
				})
		}
	})
}

type Result struct {
	Action string
	Nouns  []string
}

func parse(t *testing.T, scope Scope, match Scanner, phrases []string, goal *Result) {
	assert := testify.New(t)
	for _, in := range phrases {
		// Parse:
		res, ok := Parse(scope, match, in)
		if goal == nil && ok {
			t.Fatal("expected failure:", in)
		} else if goal != nil && !ok {
			t.Fatal("expected:", in)
		} else if ok {
			if !assert.Equal(goal.Action, res.Action) {
				break
			}
			var nouns []string
			for _, rank := range res.Matches {
				nouns = append(nouns, rank.Nouns...)
			}
			if !assert.Equal(goal.Nouns, nouns) {
				break
			}
		}
	}
}

// // var ADirection = append(shortDirections, directions...)

// // 		t.Run("examine dir", func(t *testing.T) {
// // 			for _, dir := range ADirection {
// // 				parse(t, scope, grammar,
// // 					   "look/l",
// // 					dir,
// // 					&Result{
// // 						Action: "Examine",
// // 						Nouns:  (dir),
// // 					},
// // 				}
// // 			}
// //
// // 		})
// // 		t.Run("examine to dir", func(t *testing.T) {
// // 			for _, dir := range ADirection {
// // 				parse(t, scope, grammar,
// // 					   ("look to", "l to"),
// // 					dir,
// // 					&Result{
// // 						Action: "Examine",
// // 						Nouns:  (dir),
// // 					},
// // 				}
// //
// // 			}
// // 		})
// // 	})
