package index

import "sort"

// Index stores lines sorted by two keys, a major and a minor
type Index struct {
	Major  Column
	Unique bool
	Lines  []*Line
}

func (index *Index) Walk(majorKey string, visit func(other, data string) bool) (ret int) {
	a, major, minor := index.Lines, index.Major, (index.Major+1)&1
	if i, ok := index.FindFirst(0, majorKey); ok {
		for ; i < len(a); i++ {
			if line := a[i]; line[major] != majorKey {
				break
			} else {
				ret++
				if visit(line[minor], line[LineData]) {
					break
				}
			}
		}
	}
	return
}

func (index *Index) FindFirst(n int, majorKey string) (int, bool) {
	a, major := index.Lines, index.Major
	i := n + sort.Search(len(a)-n, func(i int) bool {
		return a[i+n][major] >= majorKey
	})
	exact := i < len(a) && a[i][major] == majorKey
	return i, exact
}

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

func (index *Index) Delete(i int) {
	a := index.Lines
	index.Lines = a[:i+copy(a[i:], a[i+1:])]
}

func (index *Index) Update(newLine *Line) (changed bool) {
	if index.Unique {
		changed = index.addReplace(newLine)
	} else {
		changed = index.addInsert(newLine)
	}
	return
}

func (index *Index) addReplace(newLine *Line) (changed bool) {
	a, major := index.Lines, index.Major
	majorKey := newLine[major]
	// opt: since unique means one major key, we dont have to look at the minor key.
	if i, ok := index.FindFirst(0, majorKey); !ok {
		// if we didn't find the element, insert
		changed = index.insert(i, newLine)
	} else if !newLine.match(a[i], Secondary) || !newLine.match(a[i], LineData) {
		// otherwise, if any of the data is different, replace.
		a[i] = newLine
		changed = true
	}
	return
}

func (index *Index) addInsert(newLine *Line) (changed bool) {
	a, major, minor := index.Lines, index.Major, (index.Major+1)&1
	majorKey, minorKey := newLine[major], newLine[minor]
	if i, ok := index.FindPair(0, majorKey, minorKey); !ok {
		// if we didn't find the element, insert
		changed = index.insert(i, newLine)
	} else if a[i][LineData] != newLine[LineData] {
		// otherwise, if the line data is different, replace.
		a[i] = newLine
		changed = true
	}
	return
}

// insert the passed line at i
func (index *Index) insert(i int, newLine *Line) bool {
	a := append(index.Lines, nil)
	copy(a[i+1:], a[i:])
	a[i] = newLine
	index.Lines = a
	return true
}
