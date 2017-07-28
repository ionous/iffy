package index

import "sort"

// Index stores lines sorted by two keys, a major and a minor
type Index struct {
	Major  Column
	Unique bool
	Lines  []Row
}

func (index *Index) Walk(majorKey string, visit func(other string) bool) (ret int) {
	a, major, minor := index.Lines, index.Major, (index.Major+1)&1
	if i, ok := index.FindFirst(0, majorKey); ok {
		for ; i < len(a); i++ {
			if line := a[i]; line[major] != majorKey {
				break
			} else {
				ret++
				if visit(line[minor]) {
					break
				}
			}
		}
	}
	return
}

// FindFirst major key in this index starting from row n.
// Returns the index of the key if true, or the insert point for the pair if false.
func (index *Index) FindFirst(n int, majorKey string) (int, bool) {
	a, major := index.Lines, index.Major
	i := n + sort.Search(len(a)-n, func(i int) bool {
		return a[i+n][major] >= majorKey
	})
	exact := i < len(a) && a[i][major] == majorKey
	return i, exact
}

// Find a unique pair of keys n this index starting from row n.
// Returns the index of the pair if true, or the insert point for the pair if false.
func (index *Index) FindPair(n int, majorKey, minorKey string) (int, bool) {
	a, major, minor := index.Lines, index.Major, (index.Major+1)&1
	i := n + sort.Search(len(a)-n, func(i int) (greatOrEqual bool) {
		if aMajor := a[i+n][major]; aMajor > majorKey {
			greatOrEqual = true
		} else if aMajor == majorKey {
			greatOrEqual = a[i+n][minor] >= minorKey
		}
		return
	})
	exact := i < len(a) && a[i][major] == majorKey && a[i][minor] == minorKey
	return i, exact
}

// Delete the passed row.
func (index *Index) Delete(i int) {
	a := index.Lines
	index.Lines = a[:i+copy(a[i:], a[i+1:])]
}

// Update adds, inserts, or sets the passed key if needed.
// It will never remove a key, and it always returns avalid pointer.
// It returns the previously existing row data if any.
// Note: the value of data is ignored.
func (index *Index) UpdateRow(p, s string) (old string) {
	if index.Unique {
		old = index.addReplace(p, s)
	} else {
		index.addInsert(p, s)
	}
	return
}

// old is empty if the element gets added; otherwise its the previous pairing
func (index *Index) addReplace(p, s string) (old string) {
	row := Row{p, s}
	a, major, minor := index.Lines, index.Major, (index.Major+1)&1
	// opt: since unique means one major key, we dont have to look at the minor key.
	if i, ok := index.FindFirst(0, row[major]); !ok {
		// if we didn't find the element, insert
		index.insert(i, row)
	} else if prev, next := a[i][minor], row[minor]; prev != next {
		// otherwise, if the secondary key is different, replace it.
		old, a[i][minor] = prev, next
	}
	return
}

// true if the pair was inserted; false if the pair already existed.
func (index *Index) addInsert(p, s string) (changed bool) {
	row := Row{p, s}
	major, minor := index.Major, (index.Major+1)&1
	if i, ok := index.FindPair(0, row[major], row[minor]); !ok {
		// if we didn't find the element, insert
		index.insert(i, row)
		changed = true
	}
	return
}

// insert the passed line at i
func (index *Index) insert(i int, row Row) {
	a := append(index.Lines, Row{})
	copy(a[i+1:], a[i:])
	a[i] = row
	index.Lines = a
}
