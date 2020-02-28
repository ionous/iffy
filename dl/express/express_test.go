package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestExpress(t *testing.T) {
	// A.num
	bigDot := &Render{&core.ObjectName{"A"}, "num"}
	// A.num * B.num
	binary := &core.ProductOf{
		bigDot,
		&Render{&core.ObjectName{"B"}, "num"},
	}
	//
	tests := []struct {
		name string
		str  string
		want interface{}
	}{
		{"literal", "5", &core.NumValue{5}},
		{"no dot", "A", &core.ObjectName{"A"}},
		{"little dot", "a.b.c",
			&Render{
				Obj: &Render{
					Obj:  &GetAt{"a"},
					Prop: "b",
				},
				Prop: "c",
			},
		},
		{"big dot", "A.num", bigDot},
		{"binary", "A.num * B.num", binary},
		{"chain", "5 + A.num * B.num",
			&core.SumOf{
				&core.NumValue{5},
				binary,
			},
		},
		{"text cmp", "'a' < 'b'",
			&core.CompareText{
				&core.TextValue{"a"},
				&core.LessThan{},
				&core.TextValue{"b"},
			},
		},
		{"num cmp", "7 >= 8",
			&core.CompareNum{
				&core.NumValue{7},
				&core.GreaterOrEqual{},
				&core.NumValue{8},
			},
		},
		{"math", "(5+6)*(1+2)",
			&core.ProductOf{
				&core.SumOf{
					&core.NumValue{5},
					&core.NumValue{6},
				},
				&core.SumOf{
					&core.NumValue{1},
					&core.NumValue{2},
				},
			},
		},
		{"logic",
			"a and (b or {isNot: c})", &core.AllTrue{[]rt.BoolEval{
				&GetAt{"a"},
				&core.AnyTrue{[]rt.BoolEval{
					&GetAt{"b"},
					&core.IsNot{
						&GetAt{"c"},
					},
				}},
			}},
		},
	}
	cmds := ops.NewOps(nil)
	unique.PanicBlocks(cmds,
		(*core.Commands)(nil),
		(*Commands)(nil))
	fac := ops.NewFactory(cmds, ops.Transformer(core.Transform))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if xs, e := template.ParseExpression(test.str); e != nil {
				t.Fatal(e)
			} else if res, e := Convert(fac, xs, nil); e != nil {
				t.Fatal(e)
			} else {
				got := res.Target().Interface()
				if want := test.want; !testify.ObjectsAreEqualValues(want, got) {
					// got != want
					t.Log(pretty.Diff(got, want))
					t.Log("got:", pretty.Sprint(got))
					t.Log("want:", pretty.Sprint(want))
					t.FailNow()
				}
			}
		})
	}
}
