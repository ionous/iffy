package chart

import (
	"github.com/ionous/errutil"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestBlocks(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testBlock(t, "", ""))
	x = x && assert.NoError(testBlock(t, "abc", "abc"))
	x = x && assert.NoError(testBlock(t, "{}", "{}"))
	// mixed: front, end
	x = x && assert.NoError(testBlock(t, "abc{}", "abc{}"))
	x = x && assert.NoError(testBlock(t, "{}abc", "{}abc"))
	x = x && assert.NoError(testBlock(t, "{}{}", "{}{}"))
	// long
	x = x && assert.NoError(testBlock(t, "abc{}d{}efg{}z", "abc{}d{}efg{}z"))
}

func TestKeys(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testBlock(t, "{key}", "{key:}"))
	x = x && assert.NoError(testBlock(t, "{ key }", "{key:}"))
	x = x && assert.Error(testBlock(t, "{1}", ignoreResult))
	x = x && assert.Error(testBlock(t, "{key1}", ignoreResult))
}

func TestTrim(t *testing.T) {
	assert, x := testify.New(t), true
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

func TestKeyTrim(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testBlock(t, "  {~key}", "{key:}"))
	x = x && assert.NoError(testBlock(t, "  { key}", "  {key:}"))
	x = x && assert.NoError(testBlock(t, "  {~key~}  ", "{key:}"))
	x = x && assert.NoError(testBlock(t, "  {key~}  ", "  {key:}"))
	x = x && assert.NoError(testBlock(t, "  {key}  ", "  {key:}  "))
}

func testBlock(t *testing.T, str string, want string) (err error) {
	t.Log("test:", str)
	p := MakeBlockParser(EmptyFactory{})
	if e := parse(&p, str); e != nil {
		err = e
	} else if res, e := p.GetBlocks(); e != nil {
		err = e
	} else if want != ignoreResult {
		got := Blocks{res}.String()
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
