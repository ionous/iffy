package debug

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

func SayIt(s string) rt.TextEval {
	return &core.Text{s}
}

type MatchNumber struct {
	Val int
}

func (m MatchNumber) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.Variable(run, "num", affine.Number); e != nil {
		err = e
	} else {
		n := a.Int()
		ret = g.BoolOf(n == m.Val)
	}
	return
}

func DetermineSay(i int) *pattern.DetermineText {
	return &pattern.DetermineText{
		Pattern: "say_me", Arguments: core.NamedArgs(
			"num", &core.FromNum{
				&core.Number{float64(i)},
			}),
	}
}

var SayPattern = pattern.TextPattern{
	pattern.CommonPattern{
		Name: "say_me",
		Prologue: []term.Preparer{
			&term.Number{Name: "num"},
		}},
	[]*pattern.TextRule{
		{nil, SayIt("Not between 1 and 3.")},
		{&MatchNumber{3}, SayIt("San!")},
		{&MatchNumber{3}, SayIt("Three!")},
		{&MatchNumber{2}, SayIt("Two!")},
		{&MatchNumber{1}, SayIt("One!")},
	},
}

var SayHelloGoodbye = core.NewActivity(
	&core.ChooseAction{
		If: &core.Bool{true},
		Do: core.MakeActivity(&core.Say{
			Text: &core.Text{"hello"},
		}),
		Else: &core.ChooseNothingElse{
			core.MakeActivity(&core.Say{
				Text: &core.Text{"goodbye"},
			}),
		},
	})

var SayHelloGoodbyeData = `{
  "type": "activity",
  "value": {
    "$EXE": [{
        "type": "execute",
        "value": {
          "type": "choose_action",
          "value": {
            "$DO": {
              "type": "activity",
              "value": {
                "$EXE": [{
                    "type": "execute",
                    "value": {
                      "type": "say_text",
                      "value": {
                        "$TEXT": {
                          "type": "text_eval",
                          "value": {
                            "type": "text_value",
                            "value": {
                              "$TEXT": {
                                "type": "text",
                                "value": "hello"
                              }}}}}}}]}},
            "$ELSE": {
              "type": "brancher",
              "value": {
                "type": "choose_nothing_else",
                "value": {
                  "$DO": {
                    "type": "activity",
                    "value": {
                      "$EXE": [
                        {
                          "type": "execute",
                          "value": {
                            "type": "say_text",
                            "value": {
                              "$TEXT": {
                                "type": "text_eval",
                                "value": {
                                  "type": "text_value",
                                  "value": {
                                    "$TEXT": {
                                      "type": "text",
                                      "value": "goodbye"
                                    }}}}}}}]}}}}},
            "$IF": {
              "type": "bool_eval",
              "value": {
                "type": "bool_value",
                "value": {
                  "$BOOL": {
                    "type": "bool",
                    "value": "$TRUE"
                  }}}}}}}]}}
`
