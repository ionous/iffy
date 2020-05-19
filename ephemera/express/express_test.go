package express

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/template"
	"github.com/kr/pretty"
)

// test the operators and quotes
func TestOperators(t *testing.T) {

	T := func(s string) rt.TextEval {
		return &core.Text{s}
	}
	N := func(n float64) rt.NumberEval {
		return &core.Number{n}
	}
	True := &core.Bool{true}
	False := &core.Bool{false}
	//
	tests := []struct {
		name string
		str  string
		want interface{}
	}{
		{"num", "5", N(5)},
		{"txt", "'5'", T("5")},
		{"bool", "false", False},
		{"T cmp", "'a' < 'b'",
			&core.CompareText{
				T("a"),
				&core.LessThan{},
				T("b"),
			},
		},
		{"num cmp", "7 >= 8",
			&core.CompareNum{
				N(7),
				&core.GreaterOrEqual{},
				N(8),
			},
		},
		{"math", "(5+6)*(1+2)",
			&core.ProductOf{
				&core.SumOf{N(5), N(6)},
				&core.SumOf{N(1), N(2)},
			},
		},
		// isNot requires command parsing
		{"logic",
			"true and (false or {isNot: true})",
			&core.AllTrue{[]rt.BoolEval{
				True,
				&core.AnyTrue{[]rt.BoolEval{
					False,
					// isNot requires command parsing
					&core.IsNot{
						True,
					},
				}},
			}},
		},
		{"global", "A",
			T("A"),
		},
		{"big dot", "A.num",
			&core.GetField{T("A"), T("num")},
		},
		{"little dot", "a.b.c",
			&core.GetField{
				&core.GetField{&core.GetVar{"a"}, T("b")},
				T("c"),
			},
		},
		{"binary", "A.num * b.num",
			&core.ProductOf{
				&core.GetField{T("A"), T("num")},
				&core.GetField{&core.GetVar{"b"}, T("num")},
			},
		},
	}

	for _, test := range tests {
		if xs, e := template.ParseExpression(test.str); e != nil {
			t.Fatal(test.name, e)
		} else if got, e := Convert(xs); e != nil {
			t.Fatal(test.name, e)
		} else if diff := pretty.Diff(got, test.want); len(diff) > 0 {
			t.Fatal("failed:", test.name, pretty.Sprint(got))
		} else {
			t.Log("ok:", test.name, pretty.Sprint(got))
		}

	}

}
