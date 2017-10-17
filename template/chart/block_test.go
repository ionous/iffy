package chart

import (
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestBlocks(t *testing.T) {
	var mock mockParser
	test := func(str string, match ...Block) (err error) {
		p := blockParser{newDirective: func() subBlockParser { return mock }}
		if end := parse(&p, str); end > 0 {
			err = errutil.New(str, endpointError(end))
		} else {
			res := p.GetBlocks()
			if matches := len(match); matches > 0 || len(res) != 0 {
				t.Log(str, pretty.Sprint(res), len(match), len(res))
				if diff := pretty.Diff(match, res); len(diff) > 0 {
					err = errutil.New(str, "mismatched results", strings.Join(diff, ";"))
				}
			}
		}
		return
	}
	assert := testify.New(t)
	dir := testBlock{}
	x := true &&
		assert.NoError(test("")) &&
		assert.NoError(test("abc", TextBlock{"abc"})) &&
		assert.NoError(test("{}", dir)) &&
		// mixed: front, end
		assert.NoError(test("abc{}", TextBlock{"abc"}, dir)) &&
		assert.NoError(test("{}abc", dir, TextBlock{"abc"})) &&
		assert.NoError(test("{}{}", dir, dir)) &&
		// long
		assert.NoError(test("abc{}d{}efg{}z", TextBlock{"abc"}, dir, TextBlock{"d"}, dir, TextBlock{"efg"}, dir, TextBlock{"z"}))
		// fake an error, for example an unclosed bracket.
	mock.err = errutil.New("error")
	x = x && assert.NoError(test("a{}b{}c", TextBlock{"a"}, ErrorBlock{mock.err}, TextBlock{"b"}, ErrorBlock{mock.err}, TextBlock{"c"}))
}

type mockParser struct{ err error }
type testBlock struct{}

func (testBlock) blockNode() {}

func (m mockParser) NewRune(r rune) (ret State) {
	if !isCloseBracket(r) {
		ret = m // loop...
	} else {
		ret = terminal // done, eat rune
	}
	return
}
func (m mockParser) GetBlock() (Block, error) {
	return testBlock{}, m.err
}
