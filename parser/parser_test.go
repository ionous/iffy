package parser_test

import (
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	. "github.com/ionous/iffy/parser"
	"github.com/kr/pretty"
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

func noun(f ...Filter) Scanner {
	return &Noun{f}
}
func nouns(f ...Filter) Scanner {
	return &Multi{f}
}

// note: we use things to exclude directions
func thing() Scanner {
	return noun(&HasClass{"things"})
}

func things() Scanner {
	return nouns(&HasClass{"things"})
}

var lookGrammar = allOf(Words("look/l"), anyOf(
	allOf(&Action{"Look"}),
	allOf(Words("at"), noun(), &Action{"Examine"}),
	// before "look inside", since inside is also direction.
	allOf(noun(&HasClass{"directions"}), &Action{"Examine"}),
	allOf(Words("to"), noun(&HasClass{"directions"}), &Action{"Examine"}),
	allOf(Words("inside/in/into/through/on"), noun(), &Action{"Search"}),
	allOf(Words("under"), noun(), &Action{"LookUnder"}),
))

var pickGrammar = allOf(Words("pick"), anyOf(
	allOf(Words("up"), things(), &Action{"Take"}),
	allOf(things(), Words("up"), &Action{"Take"}),
))

func makeObject(s ...string) *MyObject {
	name, s := s[0], s[1:]
	names := strings.Fields(name)
	s = append(s, "things")
	id := ident.IdOf(strings.Join(names, "-"))
	return &MyObject{Id: id, Names: names, Classes: s}
}

var ctx = func() (ret MyBounds) {
	ret = MyBounds{
		makeObject("something"),
		makeObject("red apple", "apples"),
		makeObject("crab apple", "apples"),
		makeObject("apple cart", "carts"),
		makeObject("red cart", "carts"),
		makeObject("torch", "devices"),
	}
	return append(ret, Directions...)
}()

type Goal interface {
	Goal() Goal // marker: returns self
}

type ActionGoal struct {
	Action string
	Nouns  []string
}

func StringIds(strs []string) (ret []ident.Id) {
	for _, str := range strs {
		ret = append(ret, ident.IdOf(str))
	}
	return
}
func IdStrings(ids []ident.Id) (ret []string) {
	for _, id := range ids {
		ret = append(ret, id.String())
	}
	return
}

type ClarifyGoal struct {
	// do we print the text here or not?
	// it might be nice for testing sake --
	// What do you want to examine
	// What do you want to look at?
	// and note, yu eed the matched "verb"?
	NounInstance string
}

type ErrorGoal struct {
	Error string
}

func (a *ActionGoal) Goal() Goal  { return a }
func (a *ClarifyGoal) Goal() Goal { return a }
func (a *ErrorGoal) Goal() Goal   { return a }

type Log interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

func parse(log Log, ctx Context, match Scanner, phrases []string, goals ...Goal) (err error) {
	for _, in := range phrases {
		fields := strings.Fields(in)
		if e := innerParse(log, ctx, match, fields, goals); e != nil {
			err = errutil.Fmt("%v for '%s'", e, in)
			break
		}
	}
	return
}

func innerParse(log Log, ctx Context, match Scanner, in []string, goals []Goal) (err error) {
	if len(goals) == 0 {
		err = errutil.New("expected some goals")
	} else {
		goal, goals := goals[0], goals[1:]
		if bounds, e := ctx.GetPlayerBounds(""); e != nil {
			err = e
		} else if res, e := match.Scan(ctx, bounds, Cursor{Words: in}); e != nil {
			// on error:
			switch g := goal.(type) {
			case *ErrorGoal:
				if e.Error() != g.Error {
					err = errutil.Fmt("mismatched error want:'%s' got:'%s'", g, e)
				} else {
					log.Log("matched error", []error{e})
				}
			case *ClarifyGoal:
				clarify := g
				switch e := e.(type) {
				case MissingObject:
					extend := append(in, clarify.NounInstance)
					err = innerParse(log, ctx, match, extend, goals)
				case AmbiguousObject:
					// println(strings.Join(in, "/"))
					// insert resolution into input.
					i, s := e.Depth, append(in, "")
					copy(s[i+1:], s[i:])
					s[i] = clarify.NounInstance
					// println(strings.Join(s, "\\"))
					err = innerParse(log, ctx, match, s, goals)
				default:
					err = errutil.Fmt("clarification not implemented for %T", e)
				}
			default:
				err = errutil.New("unexpected failure:", e)
			}
		} else if goal == nil {
			err = errutil.New("unexpected success")
		} else if g, ok := goal.(*ActionGoal); !ok {
			err = errutil.Fmt("unexpected goal %s %T for result %v", in, goal, pretty.Sprint(res))
		} else if list, ok := res.(*ResultList); !ok {
			err = errutil.New("expected result list %T", res)
		} else if last, ok := list.Last(); !ok {
			err = errutil.New("result list was empty")
		} else if act, ok := last.(ResolvedAction); !ok {
			err = errutil.New("expected resolved action %T", last)
		} else if !strings.EqualFold(act.Name, g.Action) {
			err = errutil.New("expected action", act, "got", g.Action)
		} else {
			nouns := StringIds(g.Nouns)
			objs := list.Objects()
			if want, have := strings.Join(IdStrings(nouns), ","),
				strings.Join(IdStrings(objs), ","); want != have {
				err = errutil.Fmt("expected nouns %q got %q", want, have)
			} else {
				log.Logf("matched %v", in)
			}
		}
	}
	return
}
