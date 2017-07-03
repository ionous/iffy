package index

import "github.com/ionous/errutil"

type Table struct {
	Index [Columns]Index
}

func (t *Table) Type() (ret Type) {
	uni, que := t.Index[Primary].Unique, t.Index[Secondary].Unique
	if uni {
		if que {
			ret = OneToOne
		} else {
			ret = ManyToOne
		}
	} else {
		if que {
			ret = OneToMany
		} else {
			ret = ManyToMany
		}
	}
	return
}

// Type describes which columns are unique.
type Type int

//go:generate stringer -type=Type
const (
	OneToOne Type = iota
	ManyToOne
	OneToMany
	ManyToMany
)

func MakeTable(rt Type) Table {
	var uni, que bool
	switch rt {
	case ManyToMany:
	case ManyToOne:
		uni = true
	case OneToMany:
		que = true
	case OneToOne:
		uni, que = true, true
	}

	return Table{
		[2]Index{
			Index{Primary, uni, nil},
			Index{Secondary, que, nil},
		},
	}
}

func (r *Table) Relate(primary, secondary string) (ret *KeyData, changed bool) {
	if len(primary) > 0 || len(secondary) > 0 {
		if len(primary) == 0 {
			changed = r.remove(Secondary, secondary)
		} else if len(secondary) == 0 {
			changed = r.remove(Primary, primary)
		} else {
			l := MakeKey(primary, secondary)
			// see for this, i dont really care about "changed" -- i care about "new"
			if _, ok := r.Index[Primary].Update(l); ok {
				if slot, ok := r.Index[Secondary].Update(l); !ok {
					err := errutil.New("couldnt update reverse pair", secondary, primary)
					panic(err)
				} else {
					ret, changed = slot, true
				}
			}
		}
	}
	return
}

func (r *Table) remove(c Column, majorKey string) (changed bool) {
	major, minor := c, (c+1)&1
	near, far := &r.Index[major], &r.Index[minor]
	if i, ok := near.FindFirst(0, majorKey); ok {
		changed = true // otherwise we panic
		if near.Unique {
			// if we are unique, there's only one item pair with the matching key.
			minorKey := near.Lines[i].Key[minor] // note: line.Key[major] == key
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
				if minorKey := a[n].Key[minor]; !dc.DeletePair(minorKey, majorKey) {
					err := errutil.New("remove couldnt find reverse pairs", minorKey, majorKey)
					panic(err)
				}
				if n = n + 1; n == len(a) || a[n].Key[major] != majorKey {
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
