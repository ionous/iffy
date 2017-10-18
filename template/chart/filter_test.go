package chart

import (
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestFilters(t *testing.T) {
	test := func(str string, name string, args ...Spec) (err error) {
		match := &FunctionSpec{name, args}
		p := newFilterParser(mockSpecFactory)
		if end := parse(p, str); end > 0 {
			err = errutil.New(str, endpointError(end))
		} else if res, e := p.GetFunction(); e != nil {
			err = errutil.New(str, e)
		} else if diff := pretty.Diff(match, res); len(diff) > 0 {
			err = errutil.New(str, "mismatched results", diff)
		} else {
			t.Log(str, res)
		}
		return
	}
	assert := testify.New(t)
	args := []Spec{testSpec("a"), testSpec("b"), testSpec("c")}
	//
	x := assert.Error(test("", ""))
	x = x && assert.NoError(test("go: a b c", "go", args...))
	x = x && assert.NoError(test("going! a b c", "going", args...))
	x = x && assert.NoError(test("gone? a b c", "gone", args...))
	x = x && assert.NoError(test("went!", "went"))
	x = x && assert.Error(test("go$ a b c", ""))
}
