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

func TestParser(t *testing.T) {
	obj := func(names ...string) MyObject {
		return MyObject{Id: names[0], Names: names}
	}
	scope := MyScope{
		obj("something"),
	}
	scope = append(scope, Directions()...)

	grammar := allOf(words("look/l"),
		anyOf(
			allOf(&Action{"Look"}),
			allOf(words("at"), noun(), &Action{"Examine"}),
			// before "look into", since into is also direction.
			allOf(noun(&HasClass{"direction"}), &Action{"Examine"}),
			allOf(words("to"), noun(&HasClass{"direction"}), &Action{"Examine"}),
			allOf(words("inside/in/into/through/on"), noun(), &Action{"Search"}),
			allOf(words("under"), noun(), &Action{"LookUnder"}),
		))

	t.Run("look", func(t *testing.T) {
		parse(t, scope, grammar,
			Phrases("look/l"),
			&ActionGoal{"Look", nil})
	})
	t.Run("examine", func(t *testing.T) {
		parse(t, scope, grammar,
			Phrases("look/l at something"),
			&ActionGoal{
				"Examine", sliceOf.String("something"),
			})
	})
	t.Run("search", func(t *testing.T) {
		parse(t, scope, grammar,
			Phrases("look/l inside/in/into/through/on something"),
			&ActionGoal{
				"Search", sliceOf.String("something"),
			})
	})
	t.Run("look under", func(t *testing.T) {
		parse(t, scope, grammar,
			Phrases("look/l under something"),
			&ActionGoal{
				"LookUnder", sliceOf.String("something"),
			})
	})
	t.Run("look dir", func(t *testing.T) {
		look := Phrases("look/l")
		for _, d := range directions {
			d := sliceOf.String(d)
			parse(t, scope, grammar,
				permute(look, d),
				&ActionGoal{"Examine", d})
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
				&ActionGoal{"Examine", d})
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

func parse(t *testing.T, scope Scope, match Scanner, phrases []string, goals ...Goal) {

	for _, in := range phrases {
		in := strings.Fields(in)
		if e := innerParse(scope, match, in, goals); e != nil {
			e = errutil.Fmt("%v for '%s'", e, in)
			t.Fatal(e)
			break
		}
	}
	return
}

// FIX: will need a "GetScope(actor)" empty *my* box, empty chairman's box.
func innerParse(scope Scope, match Scanner, in []string, goals []Goal) (err error) {
	if len(goals) == 0 {
		err = errutil.New("expected some goals")
	} else {
		goal, goals := goals[0], goals[1:]
		res, ok := Parse(scope, match, in)
		if !ok {
			if goal != nil {
				err = errutil.New("unexpected failure")
			}
		} else if goal == nil {
			err = errutil.New("unexpected success")
		} else if !res.Complete() {
			if clarify, ok := goal.(*ClarifyGoal); !ok {
				err = errutil.New("expected clarification")
			} else {
				if res.NeedsNoun {
					// option 1: reparse the tere
					// option 2: keep a pointer, and reparse from there.
					// for good or for ill, object already matched:
					// (maybe it should have returned a partial?)
					// so reparse is the only method.
					extend := append(in, clarify.Noun)
					err = innerParse(scope, match, extend, goals)
				} else {
					err = errutil.New("not implemented")
				}
			}
		} else if act, ok := goal.(*ActionGoal); !ok {
			err = errutil.New("unexpected result:", in, res)
		} else {
			if !strings.EqualFold(act.Action, res.Action) {
				err = errutil.New("expected action", act.Action, "got", res.Action)
			} else {
				var nouns []string
				for _, rank := range res.Matches {
					nouns = append(nouns, rank.Nouns...)
				}
				if !testify.ObjectsAreEqual(act.Nouns, nouns) {
					err = errutil.New("expected nouns", strings.Join(act.Nouns, ","), "got", strings.Join(nouns, ","))
				}
			}
		}
	}
	return
}
