package chart

import (
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestArgs(t *testing.T) {
	var mockSpecFactory specFactory = func() specParser {
		return &mockSpecParser{}
	}
	test := func(str string, match ...Spec) (err error) {
		p := argParser{newSpecParser: mockSpecFactory}
		if end := parse(&p, str); end > 0 {
			err = endpointError(end)
		} else if res, e := p.GetSpecs(); e != nil {
			err = e
		} else if diff := pretty.Diff(match, res); len(diff) > 0 {
			err = errutil.New(str, "mismatched results", pretty.Sprint(diff))
		} else {
			t.Log(str, res)
		}
		return
	}
	assert := testify.New(t)
	x := assert.Error(test(""))
	x = x && assert.NoError(test("a", testSpec("a")))
	x = x && assert.NoError(test("a b c", testSpec("a"), testSpec("b"), testSpec("c")))
	x = x && assert.NoError(test("a  b		c", testSpec("a"), testSpec("b"), testSpec("c")))
}

// stands in for head spec
type testSpec string

// implements spec:
func (testSpec) specNode() {}

// generates test blocks
type mockSpecParser struct{ identParser }

func (m mockSpecParser) GetSpec() (ret Spec, err error) {
	if n, e := m.GetName(); e != nil {
		err = e
	} else {
		ret = testSpec(n)
	}
	return
}
