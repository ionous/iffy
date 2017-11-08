package text_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template"
	"github.com/ionous/iffy/template/chart"
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/iffy/template/text"
	"github.com/kr/pretty"
	r "reflect"
	"strings"
	"testing"
)

// test the structure of keywords
// ( as opposed to the results of expressions )
func TestMeta(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		expect rt.TextEval
	}{
		{"if", "{if x} a {end}",
			&core.Join{[]rt.TextEval{
				&core.ChooseText{
					If:   x,
					True: a,
				},
			}}},
		{"else",
			"{if x} a {else} b {end}",
			&core.Join{[]rt.TextEval{
				&core.ChooseText{
					If:    x,
					True:  a,
					False: b,
				},
			}}},
		{"unless",
			"{unless x} a {otherwise} b {end}",
			&core.Join{[]rt.TextEval{
				&core.ChooseText{
					If:    &core.IsNot{x},
					True:  a,
					False: b,
				},
			}}},
		{"then",
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
			}}},
		{"chain",
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
			}}},
		{"sub",
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
			}}},
		{"only",
			"{once} a {end}",
			&core.Join{[]rt.TextEval{
				&core.StoppingText{
					Id: "$stoppingCounter#1",
					Values: []rt.TextEval{
						a,
					},
				},
			}}},
		{"once",
			"{once} a {or} b {or} c {end}",
			&core.Join{[]rt.TextEval{
				&core.StoppingText{
					Id: "$stoppingCounter#1",
					Values: []rt.TextEval{
						a, b, c,
					},
				},
			}}},
		{"cycle",
			"{cycle} a {or} b {or} c {end}",
			&core.Join{[]rt.TextEval{
				&core.CycleText{
					Id: "$cycleCounter#1",
					Values: []rt.TextEval{
						a, b, c,
					},
				},
			}}},
		{"shuffle",
			"{shuffle} a {or} b {or} c {end}",
			&core.Join{[]rt.TextEval{
				&core.ShuffleText{
					Id: "$shuffleCounter#1",
					Values: []rt.TextEval{
						a, b, c,
					},
				},
			}}},
	}
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)

	unique.PanicBlocks(cmds,
		(*core.Commands)(nil))

	for _, test := range tests {
		k, str, expect := test.name, test.str, test.expect
		if strings.HasPrefix(k, "x") {
			continue
		}
		t.Run(k, func(t *testing.T) {
			t.Log("testing", str)
			if dirs, e := chart.Parse(str); e != nil {
				t.Fatal(e)
			} else {
				t.Log(template.Format(dirs))
				fac := ops.NewFactory(cmds, ops.TransformFunction{core.Transform})
				cmdr := commander{t, fac, make(ident.Counters)}
				//
				if cmd, e := text.ConvertDirectives(&cmdr, dirs); e != nil {
					t.Fatal(e)
				} else {
					res := cmd.Target().Interface()
					d := pretty.Diff(res, expect)
					if len(d) > 0 {
						t.Log(d)
						t.Log("got:", pretty.Sprint(res))
						t.Log("want:", pretty.Sprint(expect))
						t.FailNow()
					}
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

var x = &core.Get{&core.TopObject{}, "x"}
var y = &core.Get{&core.TopObject{}, "y"}
var z = &core.Get{&core.TopObject{}, "z"}
var a = &core.Join{[]rt.TextEval{&core.Text{" a "}}}
var b = &core.Join{[]rt.TextEval{&core.Text{" b "}}}
var c = &core.Join{[]rt.TextEval{&core.Text{" c "}}}

// mockExpression is a mock directive parser
type commander struct {
	t *testing.T
	*ops.Factory
	gen ident.Counters
}

func (f *commander) CreateName(group string) (string, error) {
	return f.gen.NewName(group), nil
}

// our input is always one letter: xs,y,z:
// and we generate an object property command for the test.
func (f *commander) CreateExpression(xs postfix.Expression, hint r.Type) (ret *ops.Command, err error) {
	if len(xs) != 1 {
		err = errutil.New("test expected 1 element")
	} else {
		switch fn := xs[0].(type) {
		case template.Reference:
			if cmd, e := f.CreateCommand("get"); e != nil {
				err = e
			} else if e := cmd.Position("@"); e != nil {
				err = e
			} else if e := cmd.Position(fn.String()); e != nil {
				err = e
			} else {
				ret = cmd
			}
		case template.Quote:
			if cmd, e := f.CreateCommand("text"); e != nil {
				err = e
			} else if e := cmd.Position(fn.Value()); e != nil {
				err = e
			} else {
				ret = cmd
			}
		default:
			err = errutil.Fmt("test encounted unknown type %T", fn)
		}
	}
	return
}
