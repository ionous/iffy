package parser_test

import (
	"github.com/ionous/errutil"
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
func nouns(f ...Filter) Scanner {
	return &Multi{f}
}

// note: we use things to exclude directions
func things() Scanner {
	return nouns(&HasClass{"things"})
}

var lookGrammar = allOf(words("look/l"), anyOf(
	allOf(&Action{"Look"}),
	allOf(words("at"), noun(), &Action{"Examine"}),
	// before "look inside", since inside is also direction.
	allOf(noun(&HasClass{"directions"}), &Action{"Examine"}),
	allOf(words("to"), noun(&HasClass{"directions"}), &Action{"Examine"}),
	allOf(words("inside/in/into/through/on"), noun(), &Action{"Search"}),
	allOf(words("under"), noun(), &Action{"LookUnder"}),
))

var pickGrammar = allOf(words("pick"), anyOf(
	allOf(words("up"), things(), &Action{"Take"}),
	allOf(things(), words("up"), &Action{"Take"}),
))

func makeObject(s ...string) *MyObject {
	name, s := s[0], s[1:]
	names := strings.Fields(name)
	s = append(s, "things")
	return &MyObject{Id: strings.Join(names, "-"), Names: names, Classes: s}
}

var ctx = func() (ret MyScope) {
	ret = MyScope{
		makeObject("something"),
		makeObject("red apple", "apples"),
		makeObject("crab apple", "apples"),
		makeObject("apple cart", "carts"),
		makeObject("red cart", "carts"),
		makeObject("apple", "apples"),
	}
	return append(ret, Directions...)
}()

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
				"apple")})
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

func TestDisambiguation(t *testing.T) {
	grammar := lookGrammar
	t.Run("trailing noun", func(t *testing.T) {
		e := parse(ctx, grammar,
			Phrases("look at"),
			&ClarifyGoal{"something"},
			&ActionGoal{"Examine", sliceOf.String("something")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("shared names", func(t *testing.T) {
		e := parse(ctx, grammar,
			Phrases("look at red cart"),
			&ActionGoal{"Examine", sliceOf.String("red-cart")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("ambiguous shared names", func(t *testing.T) {
		e := parse(ctx, grammar,
			Phrases("look at cart"),
			&ClarifyGoal{"red"},
			&ActionGoal{"Examine", sliceOf.String("red-cart")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("ambiguous loops", func(t *testing.T) {
		e := parse(ctx, grammar,
			Phrases("look at cart"),
			&ClarifyGoal{"cart"},
			&ClarifyGoal{"cart"},
			&ClarifyGoal{"red"},
			&ActionGoal{"Examine", sliceOf.String("red-cart")})
		if e != nil {
			t.Fatal(e)
		}
	})

	// even though it doesn't during normal play.
	t.Run("exact name works during disambiguation", func(t *testing.T) {
		e := parse(ctx, grammar,
			Phrases("look at apple"),
			&ClarifyGoal{"apple"},
			&ActionGoal{"Examine", sliceOf.String("apple")})
		if e != nil {
			t.Fatal(e)
		}
	})

	t.Run("doubled names dont match incorrectly", func(t *testing.T) {
		e := parse(ctx, grammar,
			Phrases("look at apple apple apple cart"),
			&ActionGoal{"Examine", sliceOf.String("apple-cart")})
		if e != nil {
			t.Fatal(e)
		}
	})
}

func TestParser(t *testing.T) {
	grammar := lookGrammar

	t.Run("look", func(t *testing.T) {
		e := parse(ctx, grammar,
			Phrases("look/l"),
			&ActionGoal{"Look", nil})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("examine", func(t *testing.T) {
		e := parse(ctx, grammar,
			Phrases("look/l at something"),
			&ActionGoal{
				"Examine", sliceOf.String("something"),
			})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("search", func(t *testing.T) {
		e := parse(ctx, grammar,
			Phrases("look/l inside/in/into/through/on something"),
			&ActionGoal{
				"Search", sliceOf.String("something"),
			})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("look under", func(t *testing.T) {
		e := parse(ctx, grammar,
			Phrases("look/l under something"),
			&ActionGoal{
				"LookUnder", sliceOf.String("something"),
			})
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("look dir", func(t *testing.T) {
		look := Phrases("look/l")
		for _, d := range directions {
			d := sliceOf.String(d)
			if e := parse(ctx, grammar,
				permute(look, d),
				&ActionGoal{"Examine", d}); e != nil {
				t.Fatal(e)
				break
			}
		}
	})
	t.Run("look no dir", func(t *testing.T) {
		e := parse(ctx, grammar,
			Phrases("look something"),
			nil)
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("look to dir", func(t *testing.T) {
		lookTo := Phrases("look/l to")
		for _, d := range directions {
			d := sliceOf.String(d)
			if e := parse(ctx, grammar,
				permute(lookTo, d),
				&ActionGoal{"Examine", d}); e != nil {
				t.Fatal(e)
				break
			}
		}
	})
}

type Goal interface {
	Goal() Goal // marker: retuns self
}

type ActionGoal struct {
	Action string
	Nouns  []string
}

type ClarifyGoal struct {
	// do we print the text here or not?
	// it might be nice for testing sake --
	// What do you want to examine
	// What do you want to look at?
	// and note, yu eed the matched "verb"?
	Noun string
}

func (a *ActionGoal) Goal() Goal {
	return a
}

func (a *ClarifyGoal) Goal() Goal {
	return a
}

func parse(ctx Context, match Scanner, phrases []string, goals ...Goal) (err error) {
	for _, in := range phrases {
		fields := strings.Fields(in)
		if e := innerParse(ctx, match, fields, goals); e != nil {
			err = errutil.Fmt("%v for '%s'", e, in)
			break
		}
	}
	return
}

// FIX: will need a "GetScope(actor)" empty *my* box, empty chairman's box.
func innerParse(ctx Context, match Scanner, in []string, goals []Goal) (err error) {
	if len(goals) == 0 {
		err = errutil.New("expected some goals")
	} else {
		goal, goals := goals[0], goals[1:]
		if res, e := match.Scan(ctx, Cursor{Words: in}); e != nil {

			if clarify, clarifies := goal.(*ClarifyGoal); !clarifies {
				if goal != nil {
					err = errutil.New("unexpected failure", e)
				}
			} else {
				switch e := e.(type) {
				case MissingObject:
					extend := append(in, clarify.Noun)
					err = innerParse(ctx, match, extend, goals)
				case AmbiguousObject:
					//
					// println(strings.Join(in, "/"))
					next := append(in[:e.Depth], clarify.Noun)
					next = append(next, in[e.Depth:]...)
					// println(strings.Join(next, "\\"))

				default:
					err = errutil.Fmt("clarification not implemented for %T", e)
				}
			}
		} else if goal == nil {
			err = errutil.New("unexpected success")
		} else if g, ok := goal.(*ActionGoal); !ok {
			err = errutil.New("unexpected goal", in, goal)
		} else if list, ok := res.(*ResultList); !ok {
			err = errutil.New("expected result list %T", res)
		} else if act, ok := list.Last().(ResolvedAction); !ok {
			err = errutil.New("expected resolved action %T", list.Last())
		} else if !strings.EqualFold(act.Name, g.Action) {
			err = errutil.New("expected action", act, "got", g.Action)
		} else if objs := list.Objects(); !testify.ObjectsAreEqual(g.Nouns, objs) {
			err = errutil.New("expected nouns (", strings.Join(g.Nouns, ","), ") got (", strings.Join(objs, ","), ")")
		}
	}
	return
}
