package index

import "github.com/ionous/errutil"

// Type describes which columns are unique.
type Type int

//go:generate stringer -type=Type
const (
	OneToOne Type = iota
	ManyToOne
	OneToMany
	ManyToMany
)

type Table struct {
	Primary, Secondary Index
	Data               map[Row]interface{}
}

func (t *Table) Type() (ret Type) {
	uni, que := t.Primary.Unique, t.Secondary.Unique
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

func NewTable(relation Type) *Table {
	var uni, que bool
	switch relation {
	case ManyToMany:
	case ManyToOne:
		uni = true
	case OneToMany:
		que = true
	case OneToOne:
		uni, que = true, true
	}

	return &Table{
		Index{uni, nil},
		Index{que, nil},
		make(map[Row]interface{}),
	}
}

type OnInsert func(oldData interface{}) (newData interface{}, err error)

func NoData(interface{}) (ret interface{}, err error) {
	return
}

func (r *Table) GetData(p, s string) (ret interface{}, okay bool) {
	ret, okay = r.Data[Row{p, s}]
	return
}
func (r *Table) RelatePair(primary, secondary string, onInsert OnInsert) (ret bool, err error) {
	if len(primary) > 0 || len(secondary) > 0 {
		if len(primary) == 0 {
			ret = r.remove(secondary, &r.Secondary, &r.Primary)
		} else if len(secondary) == 0 {
			ret = r.remove(primary, &r.Primary, &r.Secondary)
		} else {
			row := Row{primary, secondary}
			if newData, e := onInsert(r.Data[row]); e != nil {
				err = e
			} else {
				if old := r.Primary.UpdateRow(primary, secondary); len(old) > 0 {
					delete(r.Data, Row{primary, old})
					ret = true
				}
				if old := r.Secondary.UpdateRow(secondary, primary); len(old) > 0 {
					delete(r.Data, Row{old, secondary})
					ret = true
				}
				r.Data[row] = newData
			}
		}
	}
	return
}

func (r *Table) AddPair(primary, secondary string, onInsert OnInsert) (err error) {
	row := Row{primary, secondary}
	if newData, e := onInsert(r.Data[row]); e != nil {
		err = e
	} else if old := r.Primary.AddRow(primary, secondary); len(old) > 0 {
		err = errutil.New("paring already exists", primary, old)
	} else if old := r.Secondary.AddRow(secondary, primary); len(old) > 0 {
		err = errutil.New("paring already exists", secondary, old)
	} else {
		r.Data[row] = newData
	}
	return
}

func (r *Table) DeletePair(major, minor string) (changed bool) {
	if pr, ok := r.Primary.FindPair(0, major, minor); ok {
		if sr, ok := r.Secondary.FindPair(0, minor, major); !ok {
			panic(errutil.New("remove couldnt find reverse pair", minor, major))
		} else {
			r.Primary.Delete(pr)
			r.Secondary.Delete(sr)
			delete(r.Data, Row{major, minor})
			changed = true
		}
	}
	return
}

// major is the major key in near, a minor key in far
func (r *Table) remove(major string, near, far *Index) (changed bool) {
	rev := near == &r.Secondary
	if i, ok := near.FindFirst(0, major); ok {
		changed = true // otherwise we panic
		if near.Unique {
			// if we are unique, there's only one item pair with the matching major.
			minor := near.Rows[i].Minor // note: line.Major == major
			near.Delete(i)
			if i, ok := far.FindPair(0, minor, major); ok {
				far.Delete(i)
			} else {
				panic(errutil.New("remove couldnt find reverse pair", minor, major))
			}
			if rev {
				delete(r.Data, Row{minor, major})
			} else {
				delete(r.Data, Row{major, minor})
			}
		} else {
			// if we are not unique, then we may have a lot of pairs.
			// we will want to delete those pairs here in near, and there in far.
			a, n, dc := near.Rows, i, NewDeletionCursor(far)
			for {
				// note the minor majors are sorted, so they are ever increasing.
				if minor := a[n].Minor; !dc.DeletePair(minor, major) {
					panic(errutil.New("remove couldnt find reverse pairs", minor, major))
				} else if rev {
					delete(r.Data, Row{minor, major})
				} else {
					delete(r.Data, Row{major, minor})
				}
				if n = n + 1; n == len(a) || a[n].Major != major {
					break
				}
			}
			// FIX: add a test that checks near.Rows after a multi-remove.
			// cut all of the pairs with the requested major
			near.Rows = a[:i+copy(a[i:], a[n:])]
			// cut all (remaining) things in the far side
			dc.Flush()
		}
	}
	return
}
