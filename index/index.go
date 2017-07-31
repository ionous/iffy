package index

import "sort"

// Index stores Rows sorted by two keys, a major and a minor
type Index struct {
	Unique bool
	Rows   []Row
}

func (index *Index) Walk(major string, visit func(other string) bool) (ret int) {
	if i, ok := index.FindFirst(0, major); ok {
		for a := index.Rows; i < len(a); i++ {
			if line := a[i]; line.Major != major {
				break
			} else {
				ret++
				if visit(line.Minor) {
					break
				}
			}
		}
	}
	return
}

// FindFirst major key in this index starting from row n.
// Returns the index of the key if true, or the insert point for the pair if false.
func (index *Index) FindFirst(n int, major string) (int, bool) {
	a := index.Rows
	i := n + sort.Search(len(a)-n, func(i int) bool {
		return a[i+n].Major >= major
	})
	exact := i < len(a) && a[i].Major == major
	return i, exact
}

// Find a unique pair of keys n this index starting from row n.
// Returns the index of the pair if true, or the insert point for the pair if false.
func (index *Index) FindPair(n int, major, minor string) (int, bool) {
	a := index.Rows
	i := n + sort.Search(len(a)-n, func(i int) (greatOrEqual bool) {
		if aMajor := a[i+n].Major; aMajor > major {
			greatOrEqual = true
		} else if aMajor == major {
			greatOrEqual = a[i+n].Minor >= minor
		}
		return
	})
	exact := i < len(a) && a[i].Major == major && a[i].Minor == minor
	return i, exact
}

// Delete the passed row.
func (index *Index) Delete(i int) {
	a := index.Rows
	index.Rows = a[:i+copy(a[i:], a[i+1:])]
}

// Update adds, inserts, or sets the passed key if needed.
// It returns the previously existing row data if any.
// Note: the value of data is ignored.
func (index *Index) UpdateRow(major, minor string) (old string) {
	if index.Unique {
		prev, _ := index.addReplace(major, minor)
		old = prev
	} else {
		index.addInsert(major, minor)
	}
	return
}

func (index *Index) AddRow(major, minor string) (old string) {
	if !index.Unique {
		index.addInsert(major, minor)
	} else if prev, at := index.addReplace(major, minor); len(prev) > 0 {
		index.Rows[at].Minor = prev // restore what replace just did
		old = prev
	}
	return
}

// old is empty if the element gets added; otherwise its the previous minor key
func (index *Index) addReplace(major, minor string) (old string, at int) {
	// opt: since unique means one major key, we dont have to look at the minor key.
	if i, ok := index.FindFirst(0, major); !ok {
		// if we didn't find the element, insert
		index.insert(i, major, minor)
		at = i
	} else if prev := index.Rows[i].Minor; prev != minor {
		// otherwise, if the minor key is different, replace it.
		index.Rows[i].Minor = minor
		old, at = prev, i
	}
	return
}

// true if the pair was inserted; false if the pair already existed.
func (index *Index) addInsert(major, minor string) (changed bool) {
	if i, ok := index.FindPair(0, major, minor); !ok {
		// if we didn't find the element, insert
		index.insert(i, major, minor)
		changed = true
	}
	return
}

// insert the passed line at i
func (index *Index) insert(i int, major, minor string) {
	a := append(index.Rows, Row{})
	copy(a[i+1:], a[i:])
	a[i] = Row{major, minor}
	index.Rows = a
}
