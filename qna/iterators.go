package qna

import "github.com/ionous/errutil"

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

func (it *NumberList) Count() int {
	return len(it.list) - it.i
}

func (it *NumberList) HasNext() bool {
	return it.i < len(it.list)
}

func (it *NumberList) GetNumber() (ret float64, err error) {
	if !it.HasNext() {
		err = StreamExceeded
	} else if e := Assign(&ret, it.list[it.i]); e != nil {
		err = e
	} else {
		it.i++
	}
	return
}

func NewTextList(list []string) *TextList {
	return &TextList{list: list}
}

func (it *TextList) Count() int {
	return len(it.list) - it.i
}

func (it *TextList) HasNext() bool {
	return it.i < len(it.list)
}

func (it *TextList) GetText() (ret string, err error) {
	if !it.HasNext() {
		err = StreamExceeded
	} else if e := Assign(&ret, it.list[it.i]); e != nil {
		err = e
	} else {
		it.i++
	}
	return
}

const StreamExceeded errutil.Error = "stream exceeded"
