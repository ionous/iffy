package stream

import (
	"github.com/ahmetb/go-linq"
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

func NewNumberStream(next linq.Iterator) rt.NumberStream {
	v, ok := next()
	return &iterator{v, ok, next}
}
func NewTextStream(next linq.Iterator) rt.TextStream {
	v, ok := next()
	return &iterator{v, ok, next}
}
func NewObjectStream(next linq.Iterator) rt.ObjectStream {
	v, ok := next()
	return &iterator{v, ok, next}
}
func NewNameStream(run rt.Runtime, list []string) rt.ObjectStream {
	names := linq.From(list).Iterate()
	return NewObjectStream(func() (ret interface{}, okay bool) {
		if n, ok := names(); ok {
			ref := n.(string)
			if obj, ok := run.GetObject(ref); !ok {
				ret, okay = Error(errutil.New("couldnt find object named", ref))
			} else {
				ret, okay = Value(obj)
			}
		}
		return
	})
}
