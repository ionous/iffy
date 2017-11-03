package chart

import (
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestSeq(t *testing.T) {
	assert := testify.New(t)
	x := true
	x = x && assert.NoError(testSeq(t, "", "")) // arguments are optional.
	x = x && assert.NoError(testSeq(t, "a", "[a]"))
	x = x && assert.NoError(testSeq(t, "x+y", "[x][y]ADD"))
	x = x && assert.NoError(testSeq(t, "x  +  y  ", "[x][y]ADD"))
	x = x && assert.NoError(testSeq(t, "(x+y)*z", "[x][y]ADD[z]MUL"))
	x = x && assert.NoError(testSeq(t, "( x + y ) * ( z ) ", "[x][y]ADD[z]MUL"))
	x = x && assert.Error(testSeq(t, "() ", ""))
	x = x && assert.Error(testSeq(t, "( x + y ) * () ", ""))
}

func testSeq(t *testing.T, str string, want string) (err error) {
	var p SequenceParser
	t.Logf("parsing: '%s'", str)
	if e := parse(&p, str); e != nil {
		err = e
	} else if exp, e := p.GetExpression(); e != nil {
		err = e
	} else {
		res := exp.String()
		t.Logf("got: '%s'", res)
		if diff := pretty.Diff(res, want); len(diff) > 0 {
			err = errutil.Fmt("wanted: '%s'", want)
		}
	}
	return
}
