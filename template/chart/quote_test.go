package chart

import (
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestQuotes(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testQ(t, "'singles'", "singles"))
	x = x && assert.NoError(testQ(t, `"doubles"`, "doubles"))
	x = x && assert.NoError(testQ(t, "'escape\"'", "escape\""))
	x = x && assert.NoError(testQ(t, `"\\"`, "\\"))
	x = x && assert.NoError(testQ(t, string([]rune{'"', '\\', 'a', '"'}),
		"\a"))
	x = x && assert.Error(testQ(t, string([]rune{'"', '\\', 'g', '"'}),
		ignoreResult))
}

func testQ(t *testing.T, str, want string) (err error, ret interface{}) {
	t.Log("test:", str)
	var p QuoteParser
	if e := parse(&p, str); e != nil {
		err = e
	} else if got, e := p.GetString(); e != nil {
		err = e
	} else if want != ignoreResult {
		if got != want {
			err = mismatched(want, got)
		} else {
			t.Log("ok", got)
		}
	}
	return err, str
}
