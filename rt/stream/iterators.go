package stream

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/assign"
)

const Exceeded errutil.Error = "stream exceeded"

type NumList struct {
	i    int
	list []float64
}

type TextList struct {
	i    int
	list []string
}

// implements rt.Iterator for a slice of float64
func NewNumList(list []float64) *NumList {
	return &NumList{list: list}
}

func (it *NumList) Remaining() int {
	return len(it.list) - it.i
}

func (it *NumList) HasNext() bool {
	return it.i < len(it.list)
}

func (it *NumList) GetNext(pv interface{}) (err error) {
	if !it.HasNext() {
		err = Exceeded
	} else if e := assign.FloatPtr(pv, it.list[it.i]); e != nil {
		err = e
	} else {
		it.i++
	}
	return
}

// implements rt.Iterator for a slice of string.
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
	} else if e := assign.StringPtr(pv, it.list[it.i]); e != nil {
		err = e
	} else {
		it.i++
	}
	return
}
