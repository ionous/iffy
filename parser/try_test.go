package parser_test

import (
	. "github.com/ionous/iffy/parser"
	testify "github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestWord(t *testing.T) {
	match := func(input, goal string) (Result, error) {
		match := &Word{goal}
		return match.Scan(nil, nil, Cursor{Words: strings.Fields(input)})
	}
	t.Run("match", func(t *testing.T) {
		if res, e := match("Beep", "beep"); testify.NoError(t, e) {
			testify.EqualValues(t, ResolvedWord{"Beep"}, res)
		}
	})
	t.Run("mismatch", func(t *testing.T) {
		if res, e := match("Boop", "beep"); testify.Error(t, e) {
			testify.Nil(t, res)
		}
	})
}

func TestVariousOf(t *testing.T) {
	words := func(goal ...string) (ret []Scanner) {
		for _, g := range goal {
			ret = append(ret, &Word{g})
		}
		return
	}

	match := func(input string, goal Scanner) (Result, error) {
		return goal.Scan(nil, nil, Cursor{Words: strings.Fields(input)})
	}
	t.Run("any", func(t *testing.T) {
		wordList := words("beep", "blox")
		t.Run("match", func(t *testing.T) {

			if res, e := match("Beep", &AnyOf{wordList}); testify.NoError(t, e) {
				testify.EqualValues(t,
					ResolvedWord{"Beep"}, res)
			}
			if res, e := match("Blox", &AnyOf{wordList}); testify.NoError(t, e) {
				testify.EqualValues(t,
					ResolvedWord{"Blox"}, res)
			}
		})
		t.Run("mismatch", func(t *testing.T) {
			if res, e := match("Boop", &AnyOf{wordList}); testify.Error(t, e) {
				testify.Nil(t, res)
			}
			if res, e := match("Beep", &AnyOf{}); testify.Error(t, e) {
				testify.Nil(t, res)
			}
			if res, e := match("", &AnyOf{}); testify.Error(t, e) {
				testify.Nil(t, res)
			}
		})
	})

	t.Run("all", func(t *testing.T) {
		wordList := words("beep", "blox")
		t.Run("match", func(t *testing.T) {
			assert := testify.New(t)
			if res, e := match("Beep BLOX", &AllOf{wordList}); assert.NoError(e) {
				if res, ok := res.(*ResultList); assert.True(ok) {
					matched, res := res.WordsMatched(), res.Results()
					assert.EqualValues(2, matched)
					assert.EqualValues([]Result{
						ResolvedWord{"Beep"},
						ResolvedWord{"BLOX"}},
						res)
				}
			}
		})
		t.Run("mismatch", func(t *testing.T) {
			if res, e := match("BLOX Beep", &AllOf{wordList}); testify.Error(t, e) {
				testify.Nil(t, res)
			}
			if res, e := match("Beep", &AllOf{wordList}); testify.Error(t, e) {
				testify.Nil(t, res)
			}
			if res, e := match("BLOX", &AllOf{wordList}); testify.Error(t, e) {
				testify.Nil(t, res)
			}
			if res, e := match("Boop", &AllOf{wordList}); testify.Error(t, e) {
				testify.Nil(t, res)
			}
			if res, e := match("", &AllOf{wordList}); testify.Error(t, e) {
				testify.Nil(t, res)
			}
			if res, e := match("", &AllOf{}); testify.Error(t, e) {
				testify.Nil(t, res)
			}
		})
	})
}
