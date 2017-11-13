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
	p := BlockParser{factory: EmptyFactory{}}
	if e := Parse(&p, str); e != nil {
		err = e
	} else if res, e := p.GetDirectives(); e != nil {
		err = e
	} else if want != ignoreResult {
		got := Format(res)
		if got == want {
			t.Log("ok:", got)
		} else {
			err = mismatched(want, got)
		}
	}
	return
}

func mismatched(want, got string) error {
	return errutil.Fmt("want(%d): %s; != got(%d): %s.", len(want), want, len(got), got)
}
