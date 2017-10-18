package chart

import (
	"github.com/ionous/errutil"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestEpilogue(t *testing.T) {
	canTrim, cnt := true, 0
	test := func(str, m string, r rune) (err error) {
		p := newEpilogueParser(canTrim)
		if end := parse(p, str); end > 0 {
			err = errutil.New("test", cnt, str, endpointError(end))
		} else if v, x, e := p.GetResult(); e != nil {
			err = errutil.New("test", cnt, str, e)
		} else if v != m {
			err = errutil.New("test", cnt, str, "unexpected value", v)
		} else if x != r {
			err = errutil.New("test", cnt, str, "unexpected terminal", x)
		}
		cnt++
		return
	}
	assert := testify.New(t)
	x := true
	// empty endings
	x = x && assert.NoError(test("|", "", '|'))
	x = x && assert.NoError(test("}", "", '}'))
	x = x && assert.NoError(test("~}", "", '~'))
	// once more with some content
	x = x && assert.NoError(test("ab|", "ab", '|'))
	x = x && assert.NoError(test("ab}", "ab", '}'))
	x = x && assert.NoError(test("ab~}", "ab", '~'))
	// once more with internal traiing spaces
	x = x && assert.NoError(test("ab |", "ab", '|'))
	x = x && assert.NoError(test("ab }", "ab", '}'))
	x = x && assert.NoError(test("ab ~}", "ab", '~'))
	// trim alone is an error
	x = x && assert.Error(test("ab~ }", "", 0))
	// trim and brackets in a quote should be fine.
	x = x && assert.NoError(test("'~'}", "'~'", '}'))
	x = x && assert.NoError(test("'}'}", "'}'", '}'))
	x = x && assert.NoError(test("'~}'}", "'~}'", '}'))
	// doubled filter chars are doubled filter chars
	x = x && assert.NoError(test("ab||}", "ab||", '}'))
	// no terminal  errors
	x = x && assert.Error(test("ab", "", 0))
	canTrim = false
	x = x && assert.Error(test("ab~}", "ab", '~'))
}
