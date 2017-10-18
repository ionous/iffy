package chart

import (
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestFilters(t *testing.T) {
	test := func(str string, name string, args ...Argument) (err error) {
		match := newFunction(name, args...)
		p := newFilterParser(mockArgFactory)
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
	args := []Argument{MockArg("a"), MockArg("b"), MockArg("c")}
	//
	x := true
	x = x && assert.Error(test("", ""))
	x = x && assert.NoError(test("go: a b c", "go", args...))
	x = x && assert.NoError(test("going! a b c", "going", args...))
	x = x && assert.NoError(test("gone? a b c", "gone", args...))
	x = x && assert.NoError(test("went!", "went"))
	x = x && fails(t, test("! a b c", ""))
	x = x && fails(t, test("go$ a b c", ""))
}

func fails(t *testing.T, e error) (okay bool) {
	if testify.Error(t, e) {
		t.Log("ok", e)
		okay = true
	}
	return
}
