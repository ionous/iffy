package index

import "github.com/ionous/errutil"

type Table struct {
	Index [Columns]Index
	Data  map[Row]interface{}
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
		make(map[Row]interface{}),
	}
}

type OnInsert func(oldData interface{}) (newData interface{}, err error)

func NoData(interface{}) (ret interface{}, err error) {
	return
}

func (r *Table) Relate(primary, secondary string, onInsert OnInsert) (ret bool, err error) {
	if len(primary) > 0 || len(secondary) > 0 {
		if len(primary) == 0 {
			ret = r.remove(Secondary, secondary)
		} else if len(secondary) == 0 {
			ret = r.remove(Primary, primary)
		} else {
			row := Row{primary, secondary}
			if newData, e := onInsert(r.Data[row]); e != nil {
				err = e
			} else {
				near, far := &r.Index[Primary], &r.Index[Secondary]
				if old := near.UpdateRow(primary, secondary); len(old) > 0 {
					delete(r.Data, Row{primary, old})
				}
				if old := far.UpdateRow(primary, secondary); len(old) > 0 {
					delete(r.Data, Row{old, secondary})
				}
				r.Data[row] = newData
				ret = true
			}
		}
	}
	return
}

func (r *Table) DeletePair(majorKey, minorKey string) (changed bool) {
	near, far := &r.Index[Primary], &r.Index[Secondary]
	if n, ok := near.FindPair(0, majorKey, minorKey); ok {
		if f, ok := far.FindPair(0, minorKey, majorKey); !ok {
			err := errutil.New("remove couldnt find reverse pair", minorKey, majorKey)
			panic(err)
		} else {
			near.Delete(n)
			far.Delete(f)
			changed = true
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
