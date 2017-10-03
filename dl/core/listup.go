package core

import (
	"github.com/ahmetb/go-linq"
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

// really a pattern thing:
const Undecided errutil.Error = "undecided"

type Filter struct {
	List   rt.ObjListEval
	Accept rt.BoolEval
}

type Reverse struct {
	List rt.ObjListEval
}

// ListUp command generates a list of objects.
type ListUp struct {
	Source          rt.ObjectEval
	Next            rt.ObjectEval
	AllowDuplicates bool
	MaxObjects      int
}

func (f *Filter) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	if l, e := f.List.GetObjectStream(run); e != nil {
		err = e
	} else {
		q := linq.From(l).Where(func(i interface{}) bool {
			return i != nil
		}).Select(func(i interface{}) (ret interface{}) {
			el := i.(stream.ValueError)
			if el.Error != nil {
				ret = el
			} else {
				obj := el.Value.(rt.Object)
				rt.ScopeBlock(run, obj, func() {
					if ok, e := f.Accept.GetBool(run); e != nil {
						ret, _ = stream.Error(e)
					} else if ok {
						ret, _ = stream.Value(obj)
					} // we leave ret at nil
				})
			}
			return
		})
		ret = stream.NewObjectStream(q.Iterate())
	}
	return
}

func (rev *Reverse) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	if l, e := rev.List.GetObjectStream(run); e != nil {
		err = e
	} else {
		q := linq.From(l).Reverse()
		ret = stream.NewObjectStream(q.Iterate())
	}
	return
}

func (up *ListUp) GetObjectStream(run rt.Runtime) (ret rt.ObjectStream, err error) {
	if prev, e := up.Source.GetObject(run); e != nil {
		err = e
	} else if prev != nil {
		visit := makeVisited(up.AllowDuplicates)
		ret = stream.NewObjectStream(func() (ret interface{}, okay bool) {
			rt.ScopeBlock(run, prev, func() {
				if obj, e := up.Next.GetObject(run); e != nil {
					if e != Undecided {
						ret, okay = stream.Error(e)
					}
					prev = nil
				} else if !visit.inRange(up.MaxObjects) {
					e := errutil.New("too many objects in list up", visit.cnt)
					ret, okay = stream.Error(e)
					prev = nil
				} else if !visit.firstTime(obj) {
					e := errutil.New("duplicate object in list up", prev, "became", obj, "after", visit.cnt, "calls")
					ret, okay = stream.Error(e)
					prev = nil
				} else {
					ret, okay = stream.Value(obj)
					prev = obj
				}
			})
			return
		})
	}
	return
}

type visited struct {
	dupes map[rt.Object]bool
	cnt   int
}

func makeVisited(allowDupes bool) (ret visited) {
	if allowDupes {
		ret.dupes = make(map[rt.Object]bool)
	}
	return
}

func (v *visited) firstTime(obj rt.Object) (ret bool) {
	if v.dupes != nil {
		if v.dupes[obj] {
			ret = true
		} else {
			v.dupes[obj] = true
		}
	}
	v.cnt++
	return
}

func (v *visited) inRange(max int) bool {
	return max < 0 || v.cnt < max || (max == 0 && v.cnt < 25)
}
