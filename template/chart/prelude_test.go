package chart

import (
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestPrelude(t *testing.T) {
	mockBlockFactory := func() subBlockParser {
		return mockBlockParser{}
	}
	test := func(str string, prelude Argument) (err error) {
		p := newCustomPrelude(mockBlockFactory, mockArgFactory)

		if end := parse(p, str); end > 0 {
			err = errutil.New(str, endpointError(end))
		} else if res, e := p.GetArg(); e != nil {
			err = errutil.New(str, e)
		} else if diff := pretty.Diff(prelude, res); len(diff) > 0 {
			err = errutil.New(str, res)
		} else {
			t.Log(str, res)
		}
		return
	}
	assert := testify.New(t)

	x := true
	x = x && assert.NoError(test("'hello'", newQuote("'hello'")))
	x = x && assert.NoError(test("1.2", newNumber(1.2)))
	x = x && assert.NoError(test("fun!", newFunction("fun")))
	x = x && assert.NoError(test("args: mock1 mock2", newFunction("args", newMockArg("mock1"), newMockArg("mock2"))))
	x = x && assert.NoError(test("object", newRef("object")))
	x = x && assert.NoError(test("object.property", newRef("object", "property")))
	x = x && assert.NoError(test("{}", newMockDirective()))
	x = x && fails(t, test("#", nil))
}

func newQuote(s string) Argument                    { return &Quote{s} }
func newNumber(v float64) Argument                  { return &Number{v} }
func newRef(names ...string) Argument               { return &Reference{names} }
func newFunction(n string, a ...Argument) *Function { return &Function{n, a} }
func newMockDirective() MockDirective               { return MockDirective{} }
func newMockArg(x string) MockArg                   { return MockArg(x) }
