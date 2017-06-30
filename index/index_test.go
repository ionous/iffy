package index

import (
	"github.com/ionous/sliceOf"
	"github.com/stretchr/testify/suite"
	"sort"
	// "strings"
	"testing"
)

func TestIndex(t *testing.T) {
	suite.Run(t, new(IndexSuite))
}

type IndexSuite struct {
	suite.Suite
}

// TestMajorOrder verfies the index sorts by major order.
func (assert *IndexSuite) TestMajorOrder() {
	goblins := sliceOf.String(
		"marja",
		"claire",
		"hiro",
		"grace",
		"sven",
	)
	n := Index{0, false, nil}
	for _, v := range goblins {
		added := n.Add(MakeLine(v, "", ""))
		assert.True(added, v)
	}
	if assert.Len(n.Lines, len(goblins)) {
		var sorted, lines []string
		for i, l := range n.Lines {
			lines = append(lines, l.Primary())
			sorted = append(sorted, goblins[i])
		}
		sort.Strings(sorted)
		assert.NotEqual(goblins, sorted)
		assert.Equal(sorted, lines)
	}
}

func GetLines(index Index, c Column) (ret []string) {
	for _, l := range index.Lines {
		ret = append(ret, l[c])
	}
	return ret
}

// TestMinorOrder verfies (a non-unique) index sorts by minor order.
func (assert *IndexSuite) TestMinorOrder() {
	goblins := sliceOf.String(
		"marja",
		"claire",
		"hiro",
		"grace",
		"sven",
	)
	n := Index{}
	for _, v := range goblins {
		added := n.Add(MakeLine("goblin", v, ""))
		assert.True(added, v)
	}
	if assert.Len(n.Lines, len(goblins)) {
		var sorted, lines []string
		for i, l := range n.Lines {
			if assert.Equal("goblin", l.Primary()) {
				lines = append(lines, l.Secondary())
				sorted = append(sorted, goblins[i])
			}
		}
		sort.Strings(sorted)
		assert.NotEqual(goblins, sorted)
		assert.Equal(sorted, lines)
	}
}

// TestAddInsert verifies a non-unique index should insert similar records,
// but not insert exact records.
func (assert *IndexSuite) TestAddInsert() {
	n := Index{}
	//
	var changed bool
	changed = n.addInsert(MakeLine("claire", "rocky", "rocko"))
	assert.True(changed)
	assert.Len(n.Lines, 1)

	changed = n.addInsert(MakeLine("claire", "loofah", "loofy"))
	assert.True(changed)
	assert.Len(n.Lines, 2)

	//
	changed = n.addInsert(MakeLine("claire", "loofah", "loafer"))
	assert.True(changed)
	assert.Len(n.Lines, 2)

	changed = n.addInsert(MakeLine("claire", "loofah", "loafer"))
	assert.False(changed)
	assert.Len(n.Lines, 2)
	assert.Equal(sliceOf.String("loofah", "rocky"), GetLines(n, Secondary))
	assert.Equal(sliceOf.String("loafer", "rocko"), GetLines(n, LineData))
}

// TestAddReplace verifies a unique index should replace similar records,
// but not change exact records.
func (assert *IndexSuite) TestAddReplace() {
	n := Index{}
	//
	var changed bool
	changed = n.addReplace(MakeLine("claire", "rocky", "rocko"))
	assert.True(changed)
	assert.Len(n.Lines, 1)

	changed = n.addReplace(MakeLine("claire", "loofah", "loofy"))
	assert.True(changed)
	assert.Len(n.Lines, 1)

	changed = n.addReplace(MakeLine("claire", "loofah", "loafer"))
	assert.True(changed)
	assert.Len(n.Lines, 1)

	changed = n.addReplace(MakeLine("claire", "loofah", "loafer"))
	assert.False(changed)
	assert.Len(n.Lines, 1)
}

func makeLines(src ...string) []*Line {
	r := make([]*Line, len(src))
	for i, s := range src {
		r[i] = MakeLine(s, "", "")
	}
	return r
}

// TestDeleteOne verifies deleting a single entry.
func (assert *IndexSuite) TestDeleteOne() {
	src := sliceOf.String("a", "b", "c")
	for _, s := range src {
		n := Index{Lines: makeLines(src...)}
		assert.Len(n.Lines, 3)
		if missing := n.DeleteKeys(sliceOf.String(s)); assert.Nil(missing) {
			assert.Len(n.Lines, 2)
		}
	}
}

// TestDeleteSplit verifies deleting a split range within multiple keys.
func (assert *IndexSuite) TestDeleteSplit() {
	t := assert.T()
	src := sliceOf.String("a", "b", "c", "d", "e")
	for i, _ := range src {
		n := Index{Lines: makeLines(src...)}
		assert.Len(n.Lines, 5)
		var keys, match []string
		alt := false
		for l, _ := range src {
			val := src[(l+i)%5]
			if alt {
				match = append(match, val)
			} else {
				keys = append(keys, val)
			}
			alt = !alt
		}
		sort.Strings(keys) // the code expects the keys in order
		sort.Strings(match)
		t.Log("deleting", keys, "from", GetLines(n, Primary))
		if missing := n.DeleteKeys(keys); assert.Nil(missing) {
			assert.EqualValues(match, GetLines(n, Primary))
		}
	}
}

//  TestDeleteJoin verifies can delete multiple elements in a row
func (assert *IndexSuite) TestDeleteJoin() {
	n := Index{Lines: makeLines("a", "b", "c", "d")}
	if missing := n.DeleteKeys(sliceOf.String("b", "c")); assert.Len(missing, 0) {
		assert.EqualValues(sliceOf.String("a", "d"), GetLines(n, Primary))
	}
}

// TestDeleteMulti verifies can delete multiple entries with the same key
func (assert *IndexSuite) TestDeleteMulti() {
	t := assert.T()
	src := sliceOf.String("a", "b", "c", "d")
	for i, _ := range src {
		var lines []string
		lines = append(lines, src[:i]...)
		for j := 0; j < 3; j++ {
			lines = append(lines, src[i])
		}
		lines = append(lines, src[i:]...)
		t.Log("before", lines)

		n := Index{Lines: makeLines(lines...)}
		if missing := n.DeleteKeys(sliceOf.String(src[i])); assert.Len(missing, 0) {
			assert.Len(n.Lines, 3)
			t.Log("after", GetLines(n, Primary))
		}
	}
}

// TestDeleteLast verifies we can delete the last element, have more keys, and not die.
func (assert *IndexSuite) TestDeleteLast() {
	n := Index{Lines: makeLines("a", "b", "c", "d")}
	if missing := n.DeleteKeys(sliceOf.String("d", "e")); assert.Len(missing, 1) {
		assert.EqualValues(sliceOf.String("a", "b", "c"), GetLines(n, Primary))
	}
}

// TestDeleteMising verifies can have missing keys in various positions
func (assert *IndexSuite) TestDeleteMising() {
	t := assert.T()
	src := sliceOf.String("a", "b", "c", "d")
	for i, _ := range src {
		cut := src[i]
		lines := append([]string{}, src[:i]...)
		lines = append(lines, src[i+1:]...)
		n := Index{Lines: makeLines(lines...)}
		t.Log("before", lines)
		if missing := n.DeleteKeys(sliceOf.String(cut)); assert.Len(missing, 1) {
			after := GetLines(n, Primary)
			assert.EqualValues(lines, after)
			t.Log("after", after)
		}
	}
}

func makeMinorLines(pk string, src ...string) []*Line {
	r := make([]*Line, len(src))
	for i, s := range src {
		r[i] = MakeLine(pk, s, "")
	}
	return r
}

func (assert *IndexSuite) TestRemoveNone() {
	n := Index{Lines: makeMinorLines("claire", "a", "e", "g", "x")}
	removed := n.Remove("not claire")
	assert.Len(removed, 0)
}

func (assert *IndexSuite) TestRemoveMany() {
	n := Index{}
	major := sliceOf.String("alice", "bob", "claire", "max")
	minor := sliceOf.String("a", "e", "g", "x")
	for _, name := range major {
		m := makeMinorLines(name, minor...)
		n.Lines = append(n.Lines, m...)
	}
	removed := n.Remove("claire")
	assert.EqualValues(minor, removed)
	assert.Len(n.Lines, (len(major)-1)*len(minor))
}
