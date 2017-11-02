package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestArg(t *testing.T) {
	test := func(str string, match ...string) (err error) {
		p := MakeArgParser(&MatchFactory{match})
		t.Logf("parsing: '%s'", str)
		if e := parse(&p, str); e != nil {
			err = e
		} else if res, e := p.GetExpression(); e != nil {
			err = e
		} else {
			var want postfix.Expression
			for _, str := range match {
				want = append(want, Quote(str))
			}
			if diff := pretty.Diff(res, want); len(diff) > 0 {
				err = errutil.New("mismatched results", diff)
			}
		}
		return
	}
	assert := testify.New(t)
	x := true
	x = x && assert.NoError(test("")) // arguments are optional.
	x = x && assert.NoError(test("a", "a"))
	x = x && assert.NoError(test("a b c", "a", "b", "c"))
	x = x && assert.NoError(test("a  b		c", "a", "b", "c"))
	x = x && assert.NoError(test("a b c  ", "a", "b", "c"))
}
