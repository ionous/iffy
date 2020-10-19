package debug

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

func SayIt(s string) rt.TextEval {
	return &core.Text{s}
}

type MatchNumber struct {
	Val int
}

func (m MatchNumber) GetBool(run rt.Runtime) (okay bool, err error) {
	if a, e := run.GetField(object.Variables, "num"); e != nil {
		err = e
	} else if v, e := a.GetNumber(); e != nil {
		err = e
	} else {
		n := int(v)
		okay = n == int(m.Val)
	}
	return
}

func DetermineSay(i int) *pattern.DetermineText {
	return &pattern.DetermineText{
		"sayMe", pattern.NewNamedParams(
			"num", &core.FromNum{
				&core.Number{float64(i)},
			}),
	}
}

var SayPattern = pattern.TextPattern{
	pattern.CommonPattern{
		Name: "sayMe",
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
	&core.Choose{
		If: &core.Bool{true},
		True: core.NewActivity(&core.Say{
			Text: &core.Text{"hello"},
		}),
		False: core.NewActivity(&core.Say{
			Text: &core.Text{"goodbye"},
		}),
	})

var SayHelloGoodbyeData = reader.Map{
	"type": "activity",
	"value": map[string]interface{}{
		"$EXE": []interface{}{
			map[string]interface{}{
				"type": "execute",
				"value": map[string]interface{}{
					"type": "choose",
					"value": map[string]interface{}{
						"$FALSE": map[string]interface{}{
							"type": "activity",
							"value": map[string]interface{}{
								"$EXE": []interface{}{
									map[string]interface{}{
										"type": "execute",
										"value": map[string]interface{}{
											"type": "say_text",
											"value": map[string]interface{}{
												"$TEXT": map[string]interface{}{
													"type": "text_eval",
													"value": map[string]interface{}{
														"type": "text_value",
														"value": map[string]interface{}{
															"$TEXT": map[string]interface{}{
																"type":  "text",
																"value": "goodbye",
															}}}}}}}}}},
						"$IF": map[string]interface{}{
							"type": "bool_eval",
							"value": map[string]interface{}{
								"type": "bool_value",
								"value": map[string]interface{}{
									"$BOOL": map[string]interface{}{
										"type":  "bool",
										"value": "$TRUE",
									}}}},
						"$TRUE": map[string]interface{}{
							"type": "activity",
							"value": map[string]interface{}{
								"$EXE": []interface{}{
									map[string]interface{}{
										"type": "execute",
										"value": map[string]interface{}{
											"type": "say_text",
											"value": map[string]interface{}{
												"$TEXT": map[string]interface{}{
													"type": "text_eval",
													"value": map[string]interface{}{
														"type": "text_value",
														"value": map[string]interface{}{
															"$TEXT": map[string]interface{}{
																"type":  "text",
																"value": "hello",
															}}}}}}}}}}}}}}},
}
