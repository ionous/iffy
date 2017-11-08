package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
	"github.com/ionous/iffy/template/chart"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestExpr(t *testing.T) {
	// A.num
	bigDot := &Render{&core.Object{"A"}, "num"}
	// A.num * B.num
	binary := &core.Mul{
		bigDot,
		&Render{&core.Object{"B"}, "num"},
	}
	//
	tests := []struct {
		name string
		str  string
		want interface{}
	}{
		{"literal", "5", &core.Num{5}},
		{"no dot", "A", &core.Object{"A"}},
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
			&core.Add{
				&core.Num{5},
				binary,
			},
		},
		{"math", "(5+6)*(1+2)",
			&core.Mul{
				&core.Add{
					&core.Num{5},
					&core.Num{6},
				},
				&core.Add{
					&core.Num{1},
					&core.Num{2},
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
	fac := ops.NewFactory(cmds, ops.TransformFunction{core.Transform})
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if xs, e := chart.ParseExpression(test.str); e != nil {
				t.Fatal(e)
			} else if res, e := Convert(fac, xs); e != nil {
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
