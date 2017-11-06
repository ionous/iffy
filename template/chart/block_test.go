package chart

import (
	"github.com/ionous/errutil"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestBlocks(t *testing.T) {
	assert := testify.New(t)

	x := assert.NoError(testBlock(t, "", ""))
	x = x && assert.NoError(testBlock(t, "abc", "abc"))
	x = x && assert.NoError(testBlock(t, "{}", "{}"))
	// mixed: front, end
	x = x && assert.NoError(testBlock(t, "abc{}", "abc{}"))
	x = x && assert.NoError(testBlock(t, "{}abc", "{}abc"))
	x = x && assert.NoError(testBlock(t, "{}{}", "{}{}"))
	// long
	x = x && assert.NoError(testBlock(t, "abc{}d{}efg{}z", "abc{}d{}efg{}z"))
}

func TestTrim(t *testing.T) {
	assert := testify.New(t)

	x := true
	x = x && assert.NoError(testBlock(t, "{~~}", "{}"))
	x = x && assert.NoError(testBlock(t, "    {~~}    ", "{}"))

	x = x && assert.NoError(testBlock(t, "abc{~ }", "abc{}"))
	x = x && assert.NoError(testBlock(t, "abc   {~ }", "abc{}"))

	x = x && assert.NoError(testBlock(t, "{ ~}abc", "{}abc"))
	x = x && assert.NoError(testBlock(t, "{ ~}    abc", "{}abc"))

	x = x && assert.NoError(testBlock(t, "{ ~}{~ }", "{}{}"))
	x = x && assert.NoError(testBlock(t, "{ ~}  {~ }", "{}{}"))
	x = x && assert.NoError(testBlock(t, "abc {  }  d {   } efg  {  }z", "abc {}  d {} efg  {}z"))
	x = x && assert.NoError(testBlock(t, "abc {~ }  d {~ ~} efg  {~ }z", "abc{}  d{}efg{}z"))
}

func testBlock(t *testing.T, str string, want string) (err error) {
	t.Log("test:", str)
	p := MakeBlockParser(EmptyFactory{})
	if e := parse(&p, str); e != nil {
		err = e
	} else if res, e := p.GetBlocks(); e != nil {
		err = e
	} else {
		got := res.String()
		if got == want {
			t.Log("ok:", got)
		} else {
			err = errutil.New(want, "mismatched result:", got)
		}
	}
	return
}
func newText(t string) *TextBlock {
	return &TextBlock{t}
}
