package express

import (
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/template"
	"github.com/kr/pretty"
)

// test the operators and quotes
func TestOperators(t *testing.T) {
	//
	tests := []struct {
		name string
		str  string
		want interface{}
	}{
		{"num", "5", &core.Number{5}},
		{"txt", "'5'", &core.Text{"5"}},
		{"text cmp", "'a' < 'b'",
			&core.CompareText{
				&core.Text{"a"},
				&core.LessThan{},
				&core.Text{"b"},
			},
		},
		{"num cmp", "7 >= 8",
			&core.CompareNum{
				&core.Number{7},
				&core.GreaterOrEqual{},
				&core.Number{8},
			},
		},
		{"math", "(5+6)*(1+2)",
			&core.ProductOf{
				&core.SumOf{
					&core.Number{5},
					&core.Number{6},
				},
				&core.SumOf{
					&core.Number{1},
					&core.Number{2},
				},
			},
		},
	}

	for _, test := range tests {
		if xs, e := template.ParseExpression(test.str); e != nil {
			t.Fatal(test.name, e)
		} else if got, e := Convert(xs); e != nil {
			t.Fatal(test.name, e)
		} else if diff := pretty.Diff(got, test.want); len(diff) > 0 {
			t.Fatal(test.name, pretty.Sprint(got))
		} else {
			t.Log(test.name, pretty.Sprint(got))
		}

	}

}
