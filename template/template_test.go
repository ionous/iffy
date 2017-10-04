package template_test

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template"
	"github.com/kr/pretty"
	r "reflect"
	"strings"
	"testing"
)

func TestStates(t *testing.T) {
	tests := map[string]struct {
		str    string
		expect rt.TextEval
	}{
		"if": {
			"{if x} a {end}",
			&core.Join{[]rt.TextEval{
				&core.ChooseText{
					If:   x,
					True: a,
				},
			}},
		},
		"else": {
			"{if x} a {else} b {end}",
			&core.Join{[]rt.TextEval{
				&core.ChooseText{
					If:    x,
					True:  a,
					False: b,
				},
			}},
		},
		"unless": {
			"{unless x} a {otherwise} b {end}",
			&core.Join{[]rt.TextEval{
				&core.ChooseText{
					If:    &core.IsNot{x},
					True:  a,
					False: b,
				},
			}},
		},
		"then": {
			"{if x} a {else y} b {end}",
			&core.Join{[]rt.TextEval{
				&core.ChooseText{
					If:   x,
					True: a,
					False: &core.ChooseText{
						If:   y,
						True: b,
					},
				},
			}},
		},
		"chain": {
			"{if x} a {else y} b {else z} c {end}",
			&core.Join{[]rt.TextEval{
				&core.ChooseText{
					If:   x,
					True: a,
					False: &core.ChooseText{
						If:   y,
						True: b,
						False: &core.ChooseText{
							If:   z,
							True: c,
						},
					},
				},
			}},
		},
		"sub": {
			"{if x} a {else}{if y} b {else} c {end}{end}",
			&core.Join{[]rt.TextEval{
				&core.ChooseText{
					If:   x,
					True: a,
					False: &core.Join{[]rt.TextEval{
						&core.ChooseText{
							If:    y,
							True:  b,
							False: c,
						},
					}},
				},
			}},
		},
		"only": {
			"{once} a {end}",
			&core.Join{[]rt.TextEval{
				&core.StoppingText{
					Id: "$stoppingCounter#1",
					Values: []rt.TextEval{
						a,
					},
				},
			}},
		},
		"once": {
			"{once} a {or} b {or} c {end}",
			&core.Join{[]rt.TextEval{
				&core.StoppingText{
					Id: "$stoppingCounter#1",
					Values: []rt.TextEval{
						a, b, c,
					},
				},
			}},
		},
		"cycle": {
			"{cycle} a {or} b {or} c {end}",
			&core.Join{[]rt.TextEval{
				&core.CycleText{
					Id: "$cycleCounter#1",
					Values: []rt.TextEval{
						a, b, c,
					},
				},
			}},
		},
		"shuffle": {
			"{shuffle} a {or} b {or} c {end}",
			&core.Join{[]rt.TextEval{
				&core.ShuffleText{
					Id: "$shuffleCounter#1",
					Values: []rt.TextEval{
						a, b, c,
					},
				},
			}},
		},
	}
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)

	unique.PanicBlocks(cmds,
		(*core.Commands)(nil))

	for k, test := range tests {
		str, expect := test.str, test.expect
		if strings.HasPrefix(k, "x") {
			continue
		}
		t.Run(k, func(t *testing.T) {
			ts := template.MakeFactory(make(ident.Counters), directives)
			var root struct{ rt.TextEval }
			c := cmds.NewBuilder(&root, core.Xform{})
			if e := ts.Templatize(c, str); e != nil {
				t.Fatal(e)
			} else if e := c.Build(); e != nil {
				t.Fatal(e)
			} else {
				res := root.TextEval
				d := pretty.Diff(res, expect)
				if len(d) > 0 {
					t.Log(d)
					t.Log("got:", pretty.Sprint(res))
					t.Log("want:", pretty.Sprint(expect))
					t.FailNow()
				}
			}

		})
	}
}

type tstlog struct {
	t *testing.T
}

func (t tstlog) Write(p []byte) (int, error) {
	if len(p) > 0 {
		t.t.Log(string(p[:len(p)-1]))
	}
	return len(p), nil
}

var x = &core.Get{Obj: &core.TopObject{}, Prop: "x"}
var y = &core.Get{Obj: &core.TopObject{}, Prop: "y"}
var z = &core.Get{Obj: &core.TopObject{}, Prop: "z"}
var a = &core.Join{[]rt.TextEval{&core.Text{Text: " a "}}}
var b = &core.Join{[]rt.TextEval{&core.Text{Text: " b "}}}
var c = &core.Join{[]rt.TextEval{&core.Text{Text: " c "}}}

// directives is a mock directive parser
// our input is always one letter: x,y,z:
// and we just generate an object property command for the test.
func directives(c spec.Block, in []string, hint r.Type) error {
	c.Cmd("get", "@", in[0])
	return nil
}
