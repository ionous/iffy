package index

import (
	"github.com/ionous/sliceOf"
	"github.com/stretchr/testify/suite"
	"sort"
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
	n := Index{}
	for _, v := range goblins {
		added := n.addInsert(v, "")
		assert.True(added, v)
	}
	if assert.Len(n.Lines, len(goblins)) {
		var sorted, lines []string
		for i, l := range n.Lines {
			lines = append(lines, l[Primary])
			sorted = append(sorted, goblins[i])
		}
		sort.Strings(sorted)
		assert.NotEqual(goblins, sorted)
		assert.Equal(sorted, lines)
	}
}

func getKeys(index Index, c Column) (ret []string) {
	for _, l := range index.Lines {
		ret = append(ret, l[c])
	}
	return
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
		added := n.addInsert("goblin", v)
		assert.True(added, v)
	}
	if assert.Len(n.Lines, len(goblins)) {
		var sorted, lines []string
		for i, l := range n.Lines {
			if assert.Equal("goblin", l[Primary]) {
				lines = append(lines, l[Secondary])
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
	changed = n.addInsert("claire", "rocky")
	assert.True(changed)
	assert.Len(n.Lines, 1)

	changed = n.addInsert("claire", "loofah")
	assert.True(changed)
	assert.Len(n.Lines, 2)

	//
	changed = n.addInsert("claire", "loofah")
	assert.False(changed)
	assert.Len(n.Lines, 2)
}

// TestAddReplace verifies a unique index should replace similar records,
// but not change exact records.
func (assert *IndexSuite) TestAddReplace() {
	n := Index{}
	//
	var was string
	was = n.addReplace("claire", "rocky")
	assert.Empty(was)
	assert.Len(n.Lines, 1)

	was = n.addReplace("claire", "loofah")
	assert.Equal(was, "rocky")
	assert.Len(n.Lines, 1)

	was = n.addReplace("claire", "loofah")
	assert.Empty(was)
	assert.Len(n.Lines, 1)
}

func makeLines(src ...string) []Row {
	r := make([]Row, len(src))
	for i, s := range src {
		r[i] = Row{s, s}
	}
	return r
}

// adapt existing tests of deletion to deletion cursor
func deleteKeys(n *Index, keys []string) (missing []string) {
	dc := NewDeletionCursor(n)
	var ok bool
	for _, k := range keys {
		if ok = dc.DeletePair(k, k); !ok {
			missing = append(missing, k)
		}
	}
	if ok {
		dc.Flush()
	}
	return
}

// TestDeleteOne verifies deleting a single entry.
func (assert *IndexSuite) TestDeleteOne() {
	src := sliceOf.String("a", "b", "c")
	for _, s := range src {
		n := Index{Lines: makeLines(src...)}
		assert.Len(n.Lines, 3)
		if missing := deleteKeys(&n, sliceOf.String(s)); assert.Nil(missing) {
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
		t.Log("deleting", keys, "from", getKeys(n, Primary))
		if missing := deleteKeys(&n, keys); assert.Nil(missing) {
			if !assert.EqualValues(match, getKeys(n, Primary)) {
				break
			}
		} else {
			break
		}
	}
}

//  TestDeleteJoin verifies we can delete multiple elements in a row
func (assert *IndexSuite) TestDeleteJoin() {
	n := Index{Lines: makeLines("a", "b", "c", "d")}
	if missing := deleteKeys(&n, sliceOf.String("b", "c")); assert.Len(missing, 0) {
		assert.EqualValues(sliceOf.String("a", "d"), getKeys(n, Primary))
	}
}

// TestDeleteLast verifies we can delete the last element, have more keys, and not die.
func (assert *IndexSuite) TestDeleteLast() {
	n := Index{Lines: makeLines("a", "b", "c", "d")}
	if missing := deleteKeys(&n, sliceOf.String("d", "e")); assert.Len(missing, 1) {
		assert.EqualValues(sliceOf.String("a", "b", "c"), getKeys(n, Primary))
	}
}

// TestDeleteMising verifies we can have missing keys in various positions
func (assert *IndexSuite) TestDeleteMising() {
	t := assert.T()
	src := sliceOf.String("a", "b", "c", "d")
	for i, _ := range src {
		cut := src[i]
		lines := append([]string{}, src[:i]...)
		lines = append(lines, src[i+1:]...)
		n := Index{Lines: makeLines(lines...)}
		t.Log("before", lines)
		if missing := deleteKeys(&n, sliceOf.String(cut)); assert.Len(missing, 1) {
			after := getKeys(n, Primary)
			assert.EqualValues(lines, after)
			t.Log("after", after)
		}
	}
}

func makeMinorLines(pk string, src ...string) []Row {
	r := make([]Row, len(src))
	for i, s := range src {
		r[i] = Row{pk, s}
	}
	return r
}

func (assert *IndexSuite) TestFind() {
	a := makeMinorLines("abernathy", "a", "b", "c")
	b := makeMinorLines("claire", "a", "b", "c")
	c := makeMinorLines("zog", "a", "b", "c")
	lines := append(a, append(b, c...)...)
	n := Index{Lines: lines}
	if i, pk := n.FindFirst(0, "claire"); assert.True(pk, "found claire") {
		assert.Equal(b[0], lines[i])
	}
	//
	_, nopk := n.FindFirst(0, "not claire")
	assert.False(nopk, "shouldnt have found not claire")
	//
	if i, pair := n.FindPair(0, "claire", "b"); assert.True(pair) {
		assert.Equal(b[1], lines[i])
	}
	//
	_, nopair := n.FindPair(0, "claire", "missing")
	assert.False(nopair)
	//
}

func (assert *IndexSuite) TestFindDelete() {
	n := Index{}
	major := sliceOf.String("alice", "bob", "claire")
	minor := sliceOf.String("a", "b", "c")
	for _, name := range major {
		m := makeMinorLines(name, minor...)
		n.Lines = append(n.Lines, m...)
	}
	t := assert.T()
	t.Log(n.Lines)
	if i, ok := n.FindFirst(0, "claire"); assert.True(ok, "found claire") {
		line := n.Lines[i]
		assert.EqualValues("claire", line[Primary])
		assert.EqualValues("a", line[Secondary])
		n.Delete(i)
		t.Log(n.Lines)
		assert.Len(n.Lines, (len(major)*len(minor) - 1))
		if i, ok := n.FindFirst(0, "claire"); assert.True(ok, "found claire again") {
			line := n.Lines[i]
			assert.EqualValues("claire", line[Primary])
			assert.EqualValues("b", line[Secondary])
		}
	}
}
