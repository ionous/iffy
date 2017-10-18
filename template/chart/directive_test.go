package chart

import (
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestDirectives(t *testing.T) {
	test := func(str string, match *Directive) (err error) {
		p := newSubParser()
		r := parse(p, str)
		if n, e := p.GetBlock(); e != nil {
			err = errutil.New(e, r)
		} else if res, ok := n.(*Directive); !ok {
			err = errutil.Fmt("unexpected block %T", n)
		} else {
			t.Log("test:", str)
			t.Log("output:", res)
			if diff := pretty.Diff(match, res); len(diff) > 0 {
				t.Log("wanted:", match)
				err = errutil.New("mismatched results")
			}
		}
		return
	}
	assert, x := testify.New(t), true
	x = x && assert.NoError(test(`A | B: x {y} | C!`, newDir(
		newRef("A"),
		newFunction("B", newRef("x"), newDir(newRef("y"))),
		newFunction("C"))))
	x = x && assert.NoError(test(`"adam!" | capitalize! | prepend: "Hello "`, newDir(
		newQuote(`"adam!"`),
		newFunction("capitalize"),
		newFunction("prepend", newQuote(`"Hello "`)))))
	x = x && assert.NoError(test(`a.b + c | D?`, newExp(
		newRef("a", "b"), "+ c",
		newFunction("D"))))
}

func newDir(s Argument, filters ...*Function) *Directive {
	return newExp(s, "", filters...)
}
func newExp(s Argument, exp string, filters ...*Function) *Directive {
	d := &Directive{s, exp, nil}
	for _, f := range filters {
		d.Filters = append(d.Filters, *f)
	}
	return d
}
