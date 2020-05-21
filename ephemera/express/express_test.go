package express

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/template"
	"github.com/kr/pretty"
)

func T(s string) rt.TextEval {
	return &core.Text{s}
}
func N(n float64) rt.NumberEval {
	return &core.Number{n}
}

var True = &core.Bool{true}
var False = &core.Bool{false}

// test single expressions.
func TestExpressions(t *testing.T) {
	expressions := []struct {
		name string
		str  string
		want interface{}
	}{
		{"num", "5", N(5)},
		{"txt", "'5'", T("5")},
		{"bool", "false", False},
		{"T cmp", "'a' < 'b'",
			&core.CompareText{
				T("a"), &core.LessThan{}, T("b"),
			},
		},
		{"num cmp", "7 >= 8",
			&core.CompareNum{
				N(7), &core.GreaterOrEqual{}, N(8),
			},
		},
		{"math", "(5+6)*(1+2)",
			&core.ProductOf{
				&core.SumOf{N(5), N(6)},
				&core.SumOf{N(1), N(2)},
			},
		},
		// isNot requires command parsing
		{"logic", "true and (false or {isNot: true})",
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

	for _, test := range expressions {
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

// test full templates
func TestTemplates(t *testing.T) {
	var c Converter

	test := func(str string, want interface{}) (err error) {
		if xs, e := template.Parse(str); e != nil {
			err = e
		} else if got, e := c.Convert(xs); e != nil {
			err = errutil.New(e, xs)
		} else if diff := pretty.Diff(got, want); len(diff) > 0 {
			err = errutil.New("mismatch:", pretty.Sprint(got))
		}
		return
	}
	if e := test("{cycle}a{or}b{or}c{end}",
		&core.CycleText{Sequence: core.Sequence{
			Seq: "autoexp1",
			Parts: []rt.TextEval{
				T("a"), T("b"), T("c"),
			},
		}}); e != nil {
		t.Fatal(e)
	}
	if e := test("{once}a{or}b{or}c{end}",
		&core.StoppingText{Sequence: core.Sequence{
			Seq: "autoexp2",
			Parts: []rt.TextEval{
				T("a"), T("b"), T("c"),
			},
		}}); e != nil {
		t.Fatal(e)
	}
	if e := test("{shuffle}a{or}b{or}c{end}",
		&core.ShuffleText{Sequence: core.Sequence{
			Seq: "autoexp3",
			Parts: []rt.TextEval{
				T("a"), T("b"), T("c"),
			},
		}}); e != nil {
		t.Fatal(e)
	}
	if e := test("{if 7=7}boop{else}beep{end}",
		&core.ChooseText{
			If: &core.CompareNum{
				N(7), &core.EqualTo{}, N(7),
			},
			True:  T("boop"),
			False: T("beep"),
		}); e != nil {
		t.Fatal(e)
	}
	if e := test("{unless 7=7}boop{otherwise}beep{end}",
		&core.ChooseText{
			If: &core.IsNot{
				&core.CompareNum{
					N(7), &core.EqualTo{}, N(7),
				}},
			True:  T("boop"),
			False: T("beep"),
		}); e != nil {
		t.Fatal(e)
	}
	if e := test("{15|printNum!}",
		&core.PrintNum{
			Num: &core.Number{15},
		}); e != nil {
		t.Fatal(e)
	}
	if e := test("{15|printNum!} {if 7=7}boop{end}",
		&core.Join{
			Parts: []rt.TextEval{
				&core.PrintNum{N(15)},
				T(" "),
				&core.ChooseText{
					If: &core.CompareNum{
						N(7), &core.EqualTo{}, N(7),
					},
					True: T("boop"),
				},
			},
		}); e != nil {
		t.Fatal(e)
	}
	if e := test("{'world'|hello!}",
		&core.DetermineText{
			Pattern: "hello",
			Parameters: &core.Parameters{[]*core.Parameter{
				&core.Parameter{
					Name: "$1",
					From: &core.FromText{
						Val: T("world"),
					}}}}}); e != nil {
		t.Fatal(e)
	}
}
