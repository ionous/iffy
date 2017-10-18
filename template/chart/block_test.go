package chart

import (
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestDirectives(t *testing.T) {
	var mock mockBlockParser
	var mockBlockFactory blockFactory = func() subBlockParser {
		return mock
	}
	test := func(str string, match ...Block) (err error) {
		p := blockParser{newBlock: mockBlockFactory}
		if end := parse(&p, str); end > 0 {
			err = errutil.New(str, endpointError(end))
		} else {
			res := p.GetBlocks()
			if diff := pretty.Diff(match, res); len(diff) > 0 {
				err = errutil.New(str, "mismatched results", pretty.Sprint(diff))
			} else {
				t.Log(str, res)
			}
		}
		return
	}
	assert := testify.New(t)
	dir := TestDirective{}
	x := assert.NoError(test(""))
	x = x && assert.NoError(test("abc", TextBlock{"abc"}))
	x = x && assert.NoError(test("{}", dir))
	// mixed: front, end
	x = x && assert.NoError(test("abc{}", TextBlock{"abc"}, dir))
	x = x && assert.NoError(test("{}abc", dir, TextBlock{"abc"}))
	x = x && assert.NoError(test("{}{}", dir, dir))
	// long
	x = x && assert.NoError(test("abc{}d{}efg{}z", TextBlock{"abc"}, dir, TextBlock{"d"}, dir, TextBlock{"efg"}, dir, TextBlock{"z"}))
	// fake an error, for example an unclosed bracket.
	mock.err = errutil.New("error")
	err := ErrorBlock{mock.err}
	x = x && assert.NoError(test("{}", err))
	x = x && assert.NoError(test("a{}b{}c", TextBlock{"a"}, err, TextBlock{"b"}, err, TextBlock{"c"}))
}

// test block stands in for a directive
type TestDirective struct{}

func (TestDirective) blockNode() {}
func (TestDirective) argNode()   {}

// generates test blocks or an error
type mockBlockParser struct{ err error }

func (m mockBlockParser) NewRune(r rune) (ret State) {
	if !isCloseBracket(r) {
		ret = m // loop...
	} else {
		ret = terminal // done, eat rune
	}
	return
}
func (m mockBlockParser) GetBlock() (Block, error) {
	return TestDirective{}, m.err
}
