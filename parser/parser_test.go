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

func noun() Scanner {
	return &Object{}
}

func TestParser(t *testing.T) {
	obj := func(names ...string) MyObject {
		return MyObject{Id: names[0], Names: names}
	}
	scope := MyScope{
		obj("something"),
	}

	grammar :=
		allOf(words("look/l"), anyOf(
			allOf(&Action{"Look"}),
			allOf(words("at"), noun(), &Action{"Examine"}),
			allOf(words("inside/in/into/through/on"), noun(), &Action{"Search"}),
			allOf(words("under"), noun(), &Action{"LookUnder"}),
		))

		// first, we want to test a simple set of example actions,
		// all of which start the same way, but end with different actions.
		// later, we will test disambiguation; errors; multiple objects: etc.
	t.Run("look", func(t *testing.T) {
		parse(t, scope, grammar,
			"look/l", &Result{
				Action: "Look",
			})
	})
	t.Run("examine", func(t *testing.T) {
		parse(t, scope, grammar,
			"look/l at something", &Result{
				Action: "Examine",
				Nouns:  sliceOf.String("something"),
			})
	})
	t.Run("search", func(t *testing.T) {
		parse(t, scope, grammar,
			"look/l inside/in/into/through/on something",
			&Result{
				Action: "Search",
				Nouns:  sliceOf.String("something"),
			})
	})
	t.Run("look under", func(t *testing.T) {
		parse(t, scope, grammar,
			"look/l under something",
			&Result{
				Action: "LookUnder",
				Nouns:  sliceOf.String("something"),
			})
	})

}

type Result struct {
	Action string
	Nouns  []string
}

func parse(t *testing.T, scope Scope, match Scanner, phrase string, goal *Result) {
	assert := testify.New(t)
	for _, in := range MakePhrases(phrase) {
		// Parse:
		if res, ok := Parse(scope, match, in); assert.True(ok, in) {
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

// // var shortDirections = (
// // 	"n", "s", "e", "w", "ne", "nw", "se", "sw",
// // 	"u", "up", "ceiling", "above", "sky",
// // 	"d", "down", "floor", "below", "ground",
// // )
// // var directions = (
// // 	"north",
// // 	"south",
// // 	"east",
// // 	"west",
// // 	"northeast",
// // 	"northwest",
// // 	"southeast",
// // 	"southwest",
// // 	"up",   // "up above",
// // 	"down", // "ground",
// // 	"inside",
// // 	"outside")

// // var ADirection = append(shortDirections, directions...)

// // 		t.Run("consult", func(t *testing.T) {
// // 			parse(t, scope, grammar,
// // 				   "look/l",
// // 				"under something",
// // 				&Result{
// // 					Action: "LookUnder",
// // 					Nouns:  ("something"),
// // 				},
// // 			}
// //
// // 		})

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
