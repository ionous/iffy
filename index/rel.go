package index

import "github.com/ionous/errutil"

type Relation struct {
	Index [LineData]Index
}

// RelationType describes which columns are unique.
type RelationType int

//go:generate stringer -type=RelationType
const (
	ManyToMany RelationType = iota
	OneToMany
	ManyToOne
	OneToOne
)

// ManyToMany (0 << 1) | (0 << 0) false, false
// OneToMany  (0 << 1) | (1 << 0) false, true
// ManyToOne  (1 << 1) | (0 << 0) true, false
// OneToOne   (1 << 1) | (1 << 0) true, true

func MakeRelation(rt RelationType) Relation {
	return Relation{
		[LineData]Index{
			Index{Primary, (1<<1)&rt != 0, nil},
			Index{Secondary, (1<<0)&rt != 0, nil},
		},
	}
}

func (r *Relation) Relate(primary, secondary, data string) (changed bool) {
	if len(primary) > 0 || len(secondary) > 0 {
		if len(primary) == 0 {
			changed = r.remove(Secondary, secondary)
		} else if len(secondary) == 0 {
			changed = r.remove(Primary, primary)
		} else {
			l := MakeLine(primary, secondary, data)
			// relating secondary and primary using data.
			near, far := &r.Index[Primary], &r.Index[Secondary]
			if near.Update(l) {
				far.Update(l)
				changed = true
			}
		}
	}
	return
}

func (r *Relation) remove(c Column, majorKey string) (changed bool) {
	major, minor := c, (c+1)&1
	near, far := &r.Index[major], &r.Index[minor]
	if i, ok := near.FindFirst(0, majorKey); ok {
		changed = true // otherwise we panic
		if near.Unique {
			// if we are unique, there's only one item pair with the matching key.
			minorKey := near.Lines[i][minor] // note: line[major] == key
			near.Delete(i)
			if i, ok := far.FindPair(0, minorKey, majorKey); ok {
				far.Delete(i)
			} else {
				err := errutil.New("remove couldnt find reverse pair", minorKey, majorKey)
				panic(err)
			}
		} else {
			// if we are not unique, then we may have a lot of pairs.
			// we will want to delete those pairs here in near, and there in far.
			a, n, dc := near.Lines, i, NewDeletionCursor(far)
			for {
				// note the minor keys are sorted, so they are ever increasing.
				if minorKey := a[n][minor]; !dc.DeletePair(minorKey, majorKey) {
					err := errutil.New("remove couldnt find reverse pairs", minorKey, majorKey)
					panic(err)
				}
				if n = n + 1; n == len(a) || a[n][major] != majorKey {
					break
				}
			}
			// cut all of the pairs with the requested key
			a = a[:i+copy(a[i:], a[n:])]
			// cut all (remaining) things in the far side
			dc.Flush()
		}
	}
	return
}
