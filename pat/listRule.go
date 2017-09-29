package pat

import (
	"github.com/ahmetb/go-linq"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/stream"
)

// expand a query which generates value/error(s)
func adaptText(run rt.Runtime, q linq.Query) linq.Query {
	return q.SelectMany(func(i interface{}) (ret linq.Query) {
		if n, e := i.(rt.TextListEval).GetTextStream(run); e != nil {
			e, _ := stream.Error(e)
			ret = linq.Repeat(e, 1)
		} else {
			ret = linq.From(n)
		}
		return
	})
}
func adaptNumbers(run rt.Runtime, q linq.Query) linq.Query {
	return q.SelectMany(func(i interface{}) (ret linq.Query) {
		if n, e := i.(rt.NumListEval).GetNumberStream(run); e != nil {
			e, _ := stream.Error(e)
			ret = linq.Repeat(e, 1)
		} else {
			ret = linq.From(n)
		}
		return
	})
}
func adaptObjects(run rt.Runtime, q linq.Query) linq.Query {
	return q.SelectMany(func(i interface{}) (ret linq.Query) {
		if n, e := i.(rt.ObjListEval).GetObjectStream(run); e != nil {
			e, _ := stream.Error(e)
			ret = linq.Repeat(e, 1)
		} else {
			ret = linq.From(n)
		}
		return
	})
}

func splitQuery(run rt.Runtime, ps interface{}) (ret linq.Query, err error) {
	if pre, post, e := split(run, ps); e != nil {
		err = e
	} else {
		ret = linq.From(pre).Concat(linq.From(post).Reverse())
	}
	return
}

type listRule interface {
	Applies(rt.Runtime) (Flags, error)
}

// walk over the list of ListRule, and split into 2 parts
func split(run rt.Runtime, ps interface{}) (pre, post []interface{}, err error) {
	// could also Select to get Bool, Flags, error; TakeWhile on flags and Error
	next := linq.From(ps).Reverse().Iterate()
	for item, ok := next(); ok; item, ok = next() {
		if flags, e := item.(listRule).Applies(run); e != nil {
			err = e
			break
		} else if flags >= 0 {
			if flags == Postfix {
				post = append(post, item)
			} else {
				pre = append(pre, item)
				// FIX: add Replace:
				if flags == Infix {
					break
				}
			}
		}
	}
	return
}
