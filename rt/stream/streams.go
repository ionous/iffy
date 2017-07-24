package stream

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type Len interface {
	Len() int
}

func NewNumberStream(list []float64) rt.NumberStream {
	return &NumberIt{list: list}
}

type NumberIt struct {
	list []float64
	idx  int
}

func (it *NumberIt) HasNext() bool {
	return it.idx < len(it.list)
}

func (it *NumberIt) GetNext() (ret float64, err error) {
	if !it.HasNext() {
		err = rt.StreamExceeded
	} else {
		ret = it.list[it.idx]
		it.idx++
	}
	return
}

func NewTextStream(list []string) rt.TextStream {
	return &TextIt{list: list}
}

type TextIt struct {
	list []string
	idx  int
}

func (it *TextIt) HasNext() bool {
	return it.idx < len(it.list)
}

func (it *TextIt) GetNext() (ret string, err error) {
	if !it.HasNext() {
		err = rt.StreamExceeded
	} else {
		ret = it.list[it.idx]
		it.idx++
	}
	return
}

func NewObjectStream(list []rt.Object) rt.ObjectStream {
	return &ObjectIt{list: list}
}

type ObjectIt struct {
	list []rt.Object
	idx  int // FIX? can we just slice elements from list, and always use index 0?
}

func (it *ObjectIt) Len() int {
	return len(it.list)
}

func (it *ObjectIt) HasNext() bool {
	return it.idx < len(it.list)
}

func (it *ObjectIt) GetNext() (ret rt.Object, err error) {
	if !it.HasNext() {
		err = rt.StreamExceeded
	} else {
		ret = it.list[it.idx]
		it.idx++
	}
	return
}

func NewNameStream(run rt.Runtime, list []string) rt.ObjectStream {
	return &NameIt{run: run, list: list}
}

type NameIt struct {
	run  rt.Runtime
	list []string
	idx  int
}

func (it *NameIt) Len() int {
	return len(it.list)
}

func (it *NameIt) HasNext() bool {
	return it.idx < len(it.list)
}

func (it *NameIt) GetNext() (ret rt.Object, err error) {
	if !it.HasNext() {
		err = rt.StreamExceeded
	} else {
		ref := it.list[it.idx]
		if obj, ok := it.run.FindObject(ref); !ok {
			err = errutil.New("couldnt find object named", ref)
		} else {
			ret = obj
			it.idx++
		}
	}
	return
}
