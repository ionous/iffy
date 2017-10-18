package chart

import (
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestCall(t *testing.T) {
	test := func(str string, match ...Argument) (err error) {
		args := newCallParser(mockArgFactory)
		if end := parse(args, str); end > 0 {
			err = endpointError(end)
		} else if res, e := args.GetArgs(); e != nil {
			err = e
		} else if diff := pretty.Diff(match, res); len(diff) > 0 {
			err = errutil.New(str, "mismatched results", pretty.Sprint(diff))
		} else {
			t.Log(str, res)
		}
		return
	}
	assert := testify.New(t)
	x := assert.NoError(test("")) // arguments are optional.
	x = x && assert.NoError(test("a", MockArg("a")))
	x = x && assert.NoError(test("a b c", MockArg("a"), MockArg("b"), MockArg("c")))
	x = x && assert.NoError(test("a  b		c", MockArg("a"), MockArg("b"), MockArg("c")))
	x = x && assert.NoError(test("a b c  ", MockArg("a"), MockArg("b"), MockArg("c")))
}

// stands in for prelude arg
type MockArg string

// implements arg:
func (MockArg) argNode() {}

// generates test blocks
type mockArgParser struct{ identParser }

func (m mockArgParser) GetArg() (ret Argument, err error) {
	if n, e := m.GetName(); e != nil {
		err = e
	} else {
		ret = MockArg(n)
	}
	return
}

func mockArgFactory() argParser {
	return &mockArgParser{}
}
