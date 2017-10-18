package chart

import (
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestHead(t *testing.T) {
	var mockBlockFactory blockFactory = func() subBlockParser {
		return mockBlockParser{}
	}
	test := func(str string, head Argument) (err error) {
		p := newHeadParser(mockBlockFactory, mockArgFactory)

		if end := parse(p, str); end > 0 {
			err = errutil.New(str, endpointError(end))
		} else if res, e := p.GetArg(); e != nil {
			err = errutil.New(str, e)
		} else if diff := pretty.Diff(head, res); len(diff) > 0 {
			err = errutil.New(str, "mismatched results", diff)
		} else {
			t.Log(str, res)
		}
		return
	}
	assert := testify.New(t)

	x := true
	x = x && assert.NoError(test("'hello'", &QuotedArg{"'hello'"}))
	x = x && assert.NoError(test("1.2", &NumberArg{1.2}))
	x = x && assert.NoError(test("fun!", &FunctionArg{"fun", nil}))
	x = x && assert.NoError(test("args: mock1 mock2", &FunctionArg{"args", []Argument{TestArg("mock1"), TestArg("mock2")}}))
	x = x && assert.NoError(test("object.property", &ReferenceArg{[]string{"object", "property"}}))
	x = x && assert.NoError(test("{}", TestDirective{}))
	x = x && fails(t, test("#", nil))
}
