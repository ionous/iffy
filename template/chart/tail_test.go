package chart

import (
	"github.com/ionous/errutil"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestTail(t *testing.T) {
	canTrim, cnt := true, 0
	test := func(str, m string, r rune) (err error) {
		p := tailParser{canTrim: canTrim}
		cnt++
		if end := parse(&p, str); end > 0 {
			err = errutil.New("test", cnt, str, endpointError(end))
		} else if v, x, e := p.GetTail(); e != nil {
			err = errutil.New("test", cnt, str, e)
		} else if v != m {
			err = errutil.New("test", cnt, str, "unexpected value", v)
		} else if x != r {
			err = errutil.New("test", cnt, str, "unexpected terminal", x)
		}
		return
	}
	assert := testify.New(t)
	x :=
		// empty endings
		assert.NoError(test("|", "", '|')) &&
			assert.NoError(test("}", "", '}')) &&
			assert.NoError(test("~}", "", '~')) &&
			// once more with some content
			assert.NoError(test("ab|", "ab", '|')) &&
			assert.NoError(test("ab}", "ab", '}')) &&
			assert.NoError(test("ab~}", "ab", '~')) &&
			// once more with internal traiing spaces
			assert.NoError(test("ab |", "ab", '|')) &&
			assert.NoError(test("ab }", "ab", '}')) &&
			assert.NoError(test("ab ~}", "ab", '~')) &&
			// trim alone is an error
			assert.Error(test("ab~ }", "", 0)) &&
			// trim and brackets in a quote should be fine.
			assert.NoError(test("'~'}", "'~'", '}')) &&
			assert.NoError(test("'}'}", "'}'", '}')) &&
			assert.NoError(test("'~}'}", "'~}'", '}')) &&
			// doubled filter chars are doubled filter chars
			assert.NoError(test("ab||}", "ab||", '}')) &&
			// no terminal  errors
			assert.Error(test("ab", "", 0))
	// FIX: test can trim false
	canTrim = false
	_ = x && assert.Error(test("ab~}", "ab", '~'))
}
