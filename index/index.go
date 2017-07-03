package index

import "sort"

// Index stores lines sorted by two keys, a major and a minor
type Index struct {
	Major  Column
	Unique bool
	Lines  []*KeyData
}

func (index *Index) Walk(majorKey string, visit func(other string, data interface{}) bool) (ret int) {
	a, major, minor := index.Lines, index.Major, (index.Major+1)&1
	if i, ok := index.FindFirst(0, majorKey); ok {
		for ; i < len(a); i++ {
			if line := a[i]; line.Key[major] != majorKey {
				break
			} else {
				ret++
				if visit(line.Key[minor], line.Data) {
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
		return a[i+n].Key[major] >= majorKey
	})
	exact := i < len(a) && a[i].Key[major] == majorKey
	return i, exact
}

// Find a unique pair of keys n this index starting from row n.
// Returns the index of the pair if true, or the insert point for the pair if false.
func (index *Index) FindPair(n int, majorKey, minorKey string) (int, bool) {
	a, major, minor := index.Lines, index.Major, (index.Major+1)&1
	i := n + sort.Search(len(a)-n, func(i int) (greatOrEqual bool) {
		if aMajor := a[i+n].Key[major]; aMajor > majorKey {
			greatOrEqual = true
		} else if aMajor == majorKey {
			greatOrEqual = a[i+n].Key[minor] >= minorKey
		}
		return
	})
	exact := i < len(a) && a[i].Key[major] == majorKey && a[i].Key[minor] == minorKey
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
func (index *Index) Update(kd *KeyData) (ret *KeyData, changed bool) {
	if index.Unique {
		ret, changed = index.addReplace(kd)
	} else {
		ret, changed = index.addInsert(kd)
	}
	return
}

func (index *Index) addReplace(kd *KeyData) (ret *KeyData, changed bool) {
	a, major, minor := index.Lines, index.Major, (index.Major+1)&1
	// opt: since unique means one major key, we dont have to look at the minor key.
	if i, ok := index.FindFirst(0, kd.Key[major]); !ok {
		// if we didn't find the element, insert
		index.insert(i, kd)
		ret, changed = kd, true
	} else if al := a[i]; kd.Key[minor] != al.Key[minor] {
		// otherwise, if the secondary key is different, replace it.
		// return existing key/data so "views" of this are similarly updated.
		al.Key[minor] = kd.Key[minor]
		ret, changed = al, true
	} else {
		ret = al
	}
	return
}

func (index *Index) addInsert(kd *KeyData) (ret *KeyData, changed bool) {
	a, major, minor := index.Lines, index.Major, (index.Major+1)&1
	if i, ok := index.FindPair(0, kd.Key[major], kd.Key[minor]); !ok {
		// if we didn't find the element, insert
		index.insert(i, kd)
		ret, changed = kd, true
	} else {

		// otherwise, by definition our major and minor keys are equal
		// return existing key/data so "views" of this are similarly updated.
		ret = a[i]
	}
	return
}

// insert the passed line at i
func (index *Index) insert(i int, kd *KeyData) {
	a := append(index.Lines, nil)
	copy(a[i+1:], a[i:])
	a[i] = kd
	index.Lines = a
}
