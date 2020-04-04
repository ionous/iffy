package stream

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/assign"
)

const Exceeded errutil.Error = "stream exceeded"

type NumberList struct {
	i    int
	list []float64
}

type TextList struct {
	i    int
	list []string
}

func NewNumberList(list []float64) *NumberList {
	return &NumberList{list: list}
}

func (it *NumberList) Remaining() int {
	return len(it.list) - it.i
}

func (it *NumberList) HasNext() bool {
	return it.i < len(it.list)
}

func (it *NumberList) GetNext(pv interface{}) (err error) {
	if !it.HasNext() {
		err = Exceeded
	} else if e := assign.ToFloat(pv, it.list[it.i]); e != nil {
		err = e
	} else {
		it.i++
	}
	return
}

func NewTextList(list []string) *TextList {
	return &TextList{list: list}
}

func (it *TextList) Remaining() int {
	return len(it.list) - it.i
}

func (it *TextList) HasNext() bool {
	return it.i < len(it.list)
}

func (it *TextList) GetNext(pv interface{}) (err error) {
	if !it.HasNext() {
		err = Exceeded
	} else if e := assign.ToString(pv, it.list[it.i]); e != nil {
		err = e
	} else {
		it.i++
	}
	return
}
