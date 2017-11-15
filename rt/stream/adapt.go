package stream

import (
	"github.com/ahmetb/go-linq"
	"github.com/ionous/iffy/rt"
)

// Count determines the number of elements in a stream in a generic way.
type Count interface {
	Count() int
}

// iterator implements iffy's streams -- TextStream, NumberStream, ObjectStream -- in a generic way.
type iterator struct {
	value interface{}
	okay  bool
	next  func() (interface{}, bool)
}

// Iterate satisfys linq.Iterable to work with linq.
func (n *iterator) Iterate() linq.Iterator {
	return func() (ret interface{}, okay bool) {
		if ok := n.okay; ok {
			ret, okay = n.value, true
			n.advance()
		}
		return
	}
}

func (n *iterator) Count() int {
	// FIX? as an optimization, perhaps FromList could return a custom type
	// which New*Stream detects to return an alternative implementation of Count which returns the length of the original list.
	return linq.FromIterable(n).Count()
}

func (n *iterator) HasNext() bool {
	return n.okay
}
func (n *iterator) GetNext() (ret interface{}, err error) {
	if !n.okay {
		err = rt.StreamExceeded
	} else {
		curr := n.value.(ValueError)
		if v, e := curr.Value, curr.Error; e != nil {
			err = e
		} else {
			ret = v
		}
		n.advance()
	}
	return
}

func (n *iterator) advance() {
	n.value, n.okay = n.next()
}

func (n *iterator) GetText() (ret string, err error) {
	if v, e := n.GetNext(); e != nil {
		err = e
	} else {
		ret = v.(string)
	}
	return
}
func (n *iterator) GetNumber() (ret float64, err error) {
	if v, e := n.GetNext(); e != nil {
		err = e
	} else {
		ret = v.(float64)
	}
	return
}
func (n *iterator) GetObject() (ret rt.Object, err error) {
	if v, e := n.GetNext(); e != nil {
		err = e
	} else {
		ret = v.(rt.Object)
	}
	return
}
