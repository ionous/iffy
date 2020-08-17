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
func O(n string, exact bool) *core.ObjectName {
	return &core.ObjectName{&core.Text{n}, exact}
}

var True = &core.Bool{true}
var False = &core.Bool{false}

// TestExpressions single expressions.
func TestExpressions(t *testing.T) {
	t.Run("num", func(t *testing.T) {
		if e := testExpression("5", N(5)); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("txt", func(t *testing.T) {
		if e := testExpression("'5'", T("5")); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("bool", func(t *testing.T) {
		if e := testExpression("false", False); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("T cmp", func(t *testing.T) {
		if e := testExpression(
			"'a' < 'b'",
			&core.CompareText{
				T("a"), &core.LessThan{}, T("b"),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("num cmp", func(t *testing.T) {
		if e := testExpression(
			"7 >= 8",
			&core.CompareNum{
				N(7), &core.GreaterOrEqual{}, N(8),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("math", func(t *testing.T) {
		if e := testExpression(
			"(5+6)*(1+2)",
			&core.ProductOf{
				&core.SumOf{N(5), N(6)},
				&core.SumOf{N(1), N(2)},
			}); e != nil {
			t.Fatal(e)
		}
	})
	// isNot requires command parsing
	t.Run("logic", func(t *testing.T) {
		if e := testExpression(
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
			}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("global", func(t *testing.T) {
		if e := testExpression(".A",
			&core.Buffer{[]rt.Execute{
				&core.DetermineAct{"printAName",
					&core.Parameters{[]*core.Parameter{{
						"$1", &core.FromText{O("A", true)},
					}}}}}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("big dot", func(t *testing.T) {
		if e := testExpression(".A.num",
			&core.GetField{O("A", true), T("num")}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("little dot", func(t *testing.T) {
		if e := testExpression(".a.b.c",
			&core.GetField{
				&core.ObjectName{
					Name: &core.GetField{
						O("a", false),
						T("b")},
					Exactly: true,
				},
				T("c")}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("binary", func(t *testing.T) {
		if e := testExpression(".A.num * .b.num",
			&core.ProductOf{
				&core.GetField{O("A", true), T("num")},
				&core.GetField{O("b", false), T("num")},
			}); e != nil {
			t.Fatal(e)
		}
	})
}

func testExpression(str string, want interface{}) (err error) {
	if xs, e := template.ParseExpression(str); e != nil {
		err = errutil.New(e)
	} else if got, e := Convert(xs); e != nil {
		err = errutil.New(e)
	} else if diff := pretty.Diff(got, want); len(diff) > 0 {
		err = errutil.New("failed:", pretty.Sprint(got))
	}
	return
}

// test full templates
func TestTemplates(t *testing.T) {
	t.Run("cycle", func(t *testing.T) {
		if e := testTemplate("{cycle}a{or}b{or}c{end}",
			&core.CycleText{Sequence: core.Sequence{
				Seq: "autoexp1",
				Parts: []rt.TextEval{
					T("a"), T("b"), T("c"),
				},
			}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("once", func(t *testing.T) {
		if e := testTemplate("{once}a{or}b{or}c{end}",
			&core.StoppingText{Sequence: core.Sequence{
				Seq: "autoexp1",
				Parts: []rt.TextEval{
					T("a"), T("b"), T("c"),
				},
			}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("shuffle", func(t *testing.T) {
		if e := testTemplate("{shuffle}a{or}b{or}c{end}",
			&core.ShuffleText{Sequence: core.Sequence{
				Seq: "autoexp1",
				Parts: []rt.TextEval{
					T("a"), T("b"), T("c"),
				},
			}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("if", func(t *testing.T) {
		if e := testTemplate("{if 7=7}boop{else}beep{end}",
			&core.ChooseText{
				If: &core.CompareNum{
					N(7), &core.EqualTo{}, N(7),
				},
				True:  T("boop"),
				False: T("beep"),
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("unless", func(t *testing.T) {
		if e := testTemplate("{unless 7=7}boop{otherwise}beep{end}",
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
	})
	t.Run("filter", func(t *testing.T) {
		if e := testTemplate("{15|printNum!}",
			&core.PrintNum{
				Num: &core.Number{15},
			}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("span", func(t *testing.T) {
		if e := testTemplate("{15|printNum!} {if 7=7}boop{end}",
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
	})
	t.Run("indexed", func(t *testing.T) {
		if e := testTemplate("{'world'|hello!}",
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
	})
	t.Run("object", func(t *testing.T) {
		if e := testTemplate("hello {.object}",
			&core.Join{Parts: []rt.TextEval{
				T("hello "),
				&core.Buffer{[]rt.Execute{
					&core.DetermineAct{"printAName",
						&core.Parameters{[]*core.Parameter{
							&core.Parameter{"$1",
								&core.FromText{
									O("object", false),
								}}}}}}}}},
		); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("global prop", func(t *testing.T) {
		if e := testTemplate("{.Object.prop}",
			&core.GetField{O("Object", true), T("prop")},
		); e != nil {
			t.Fatal(e)
		}
	})
}
func testTemplate(str string, want interface{}) (err error) {
	if xs, e := template.Parse(str); e != nil {
		err = e
	} else if got, e := Convert(xs); e != nil {
		err = errutil.New(e, xs)
	} else if diff := pretty.Diff(got, want); len(diff) > 0 {
		err = errutil.New("mismatch:", pretty.Sprint(got))
	}
	return
}
