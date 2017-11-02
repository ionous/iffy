package chart

import (
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestBlocks(t *testing.T) {
	test := func(str string, match ...Block) (err error) {
		t.Log("test:", str)
		p := MakeBlockParser(MockFactory{})
		if e := parse(&p, str); e != nil {
			err = e
		} else if res, e := p.GetBlocks(); e != nil {
			err = e
		} else {
			t.Log("output:", res)
			if diff := pretty.Diff(match, res); len(diff) > 0 {
				t.Log("wanted:", match)
				err = errutil.New(str, "mismatched results", pretty.Sprint(diff))
			}
		}
		return
	}
	assert := testify.New(t)
	dir := &Directive{}
	abc := newText("abc")

	x := assert.NoError(test(""))
	x = x && assert.NoError(test("abc", abc))
	x = x && assert.NoError(test("{}", dir))
	// mixed: front, end
	x = x && assert.NoError(test("abc{}", abc, dir))
	x = x && assert.NoError(test("{}abc", dir, abc))
	x = x && assert.NoError(test("{}{}", dir, dir))
	// long
	x = x && assert.NoError(test("abc{}d{}efg{}z", abc, dir, newText("d"), dir, newText("efg"), dir, newText("z")))
}
