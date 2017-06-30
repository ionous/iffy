package index

import "sort"

// Index stores lines sorted by two keys, a major and a minor
type Index struct {
	Major  Column
	Unique bool
	Lines  []*Line
}

func (index *Index) Walk(key string, visit func(other, data string) bool) (ret int) {
	a, major, minor := index.Lines, index.Major, (index.Major+1)&1
	i := sort.Search(len(a), func(i int) bool {
		return a[i][major] >= key
	})
	for ; i < len(a); i++ {
		if line := a[i]; line[major] != key {
			break
		} else {
			if visit(line[minor], line[LineData]) {
				break
			}
			ret++
		}
	}
	return
}

func (index *Index) Add(newLine *Line) (changed bool) {
	if index.Unique {
		changed = index.addReplace(newLine)
	} else {
		changed = index.addInsert(newLine)
	}
	return
}

func (index *Index) addReplace(newLine *Line) (changed bool) {
	a, major := index.Lines, index.Major
	newMajor := newLine[major]
	// opt: since unique means one major key, we dont have to look at the minor key.
	i := sort.Search(len(a), func(i int) bool {
		return a[i][major] >= newMajor
	})
	// if we didn't find the element, insert
	if i == len(a) || !newLine.match(a[i], Primary) {
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
	newMajor, newMinor := newLine[major], newLine[minor]
	i := sort.Search(len(a), func(i int) bool {
		newLine := a[i]
		return newLine[major] >= newMajor && newLine[minor] >= newMinor
	})
	// if we didn't find the element ( first and second matching ) insert
	if i == len(a) || !newLine.match(a[i], Primary) || !newLine.match(a[i], Secondary) {
		changed = index.insert(i, newLine)
	} else if !newLine.match(a[i], LineData) {
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

// Remove all lines with the passed major key, returning the minor keys of those lines.
func (index *Index) Remove(key string) (removed []string) {
	// search for the start of key in a
	a, major, minor := index.Lines, index.Major, (index.Major+1)&1
	i := sort.Search(len(a), func(i int) bool {
		return a[i][major] >= key
	})
	// search for the end of key in a
	j := i
	for ; j < len(a) && a[j][major] == key; j++ {
		removed = append(removed, a[j][minor])
	}
	// its possibly that key doesn't exist, and j never moved.
	if j > i {
		// remove everything from start to end
		index.Lines = a[:i+copy(a[i:], a[j:])]
	}
	return
}

// DeleteKeys removes one or more keys, returning a list of keys not found.
func (index *Index) DeleteKeys(keys []string) (missing []string) {
	if cnt := len(keys); cnt > 0 {
		// keys are sorted. once we find one, we can clip our range.
		a, major, n := index.Lines, index.Major, 0
		for k := 0; k < cnt; {
			key := keys[k]
			// look for key starting from 'n'
			i := n + sort.Search(len(a)-n, func(i int) bool {
				return a[i+n][major] >= key
			})
			// we expect every key to exist.
			// we have to test > in case our reset ( below ) put us beyond the array.
			if i >= len(a) || a[i][major] != key {
				missing = append(missing, key)
				k++ // advance to next key, dont change n
			} else {
				// update n, checking to see how many entries match this key.
				for n = i + 1; n != len(a) && a[n][major] == key; n++ {
				}
				// update k, speculatively checking to see if entries match the following key/s.
				// ( if nothing matches, we are left at the next key and n is unchanged )
				for k++; k < len(keys); k++ {
					key, matched := keys[k], false
					for ; n != len(a) && a[n][major] == key; n++ {
						matched = true
					}
					if !matched {
						break
					}
				}
				// chop out everything from i to n.
				// ( copy returns the number of elements copied )
				a = a[:i+copy(a[i:], a[n:])]
				n = i + 1 // reset n for next loop, looking just after our cut point.
				// we know i itself doesnt have our next key, or we would have consolidated the cut above.
			}
		}
		index.Lines = a
	}
	return
}
