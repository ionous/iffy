package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
	"github.com/ionous/sliceOf"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestOperand(t *testing.T) {
	test := func(str string, match postfix.Function) (err error) {
		var p OperandParser
		if e := parse(&p, str); e != nil {
			err = errutil.New("parse error", str, e)
		} else if res, e := p.GetOperand(); e != nil {
			err = errutil.New("operand error", str, e)
		} else if diff := pretty.Diff(match, res); len(diff) > 0 {
			err = errutil.New("mismatched", str, res)
		} else {
			t.Log(str, res)
		}
		return
	}
	assert := testify.New(t)

	x := true
	x = x && assert.NoError(test("'hello'", Quote("'hello'")))
	x = x && assert.NoError(test("1.2", Number(1.2)))
	x = x && assert.NoError(test("object", Reference(sliceOf.String("object"))))
	x = x && assert.NoError(test("object.property", Reference(sliceOf.String("object", "property"))))
	x = x && fails(t, test("#", nil))
}

func fails(t *testing.T, e error) (okay bool) {
	if testify.Error(t, e) {
		t.Log("ok", e)
		okay = true
	}
	return
}
