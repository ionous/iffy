package chart

import (
	"github.com/ionous/errutil"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestQuotes(t *testing.T) {
	test := func(str string) (err error, ret interface{}) {
		var p QuoteParser
		if e := parse(&p, str); e != nil {
			err = e
		} else if n, e := p.GetString(); e != nil {
			err = e
		} else if n != str {
			err = errutil.New("mismatched strings", str, n)
		}
		return err, str
	}
	assert := testify.New(t)
	x := true
	x = x && assert.NoError(test("'singles'"))
	x = x && assert.NoError(test(`"doubles"`))
	x = x && assert.NoError(test("'escape\"'"))
	x = x && assert.NoError(test(`"\\"`))
	x = x && assert.NoError(test(string([]rune{'"', '\\', 'a', '"'})))
	x = x && assert.Error(test(string([]rune{'"', '\\', 'g', '"'})))
}
