package rt

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt/stream"
)

type Value interface {
	BoolEval
	NumberEval
	TextEval
	NumListEval
	TextListEval
}

type NoValue struct{}

type BoolValue struct {
	NoValue
	Value bool
}

type NumberValue struct {
	NoValue
	Value float64
}

type TextValue struct {
	NoValue
	Value string
}

type NumListValue struct {
	NoValue
	Value []float64
}

type TextListValue struct {
	NoValue
	Value []string
}

func (n NoValue) GetBool(Runtime) (ret bool, err error) {
	err = errutil.New("value is not a bool")
	return
}
func (n NoValue) GetNumber(Runtime) (ret float64, err error) {
	err = errutil.New("value is not a number")
	return
}
func (n NoValue) GetText(Runtime) (ret string, err error) {
	err = errutil.New("value is not a text")
	return
}
func (n NoValue) GetNumberStream(Runtime) (ret Iterator, err error) {
	err = errutil.New("value is not a number stream")
	return
}
func (n NoValue) GetTextStream(Runtime) (ret Iterator, err error) {
	err = errutil.New("value is not a text stream")
	return
}

//
func (n *BoolValue) GetBool(Runtime) (ret bool, err error) {
	ret = n.Value
	return
}
func (n *NumberValue) GetNumber(Runtime) (ret float64, err error) {
	ret = n.Value
	return
}
func (n *TextValue) GetText(Runtime) (ret string, err error) {
	ret = n.Value
	return
}
func (n *NumListValue) GetNumberStream(Runtime) (ret Iterator, err error) {
	ret = stream.NewNumList(n.Value)
	return
}
func (n *TextListValue) GetTextStream(Runtime) (ret Iterator, err error) {
	ret = stream.NewTextList(n.Value)
	return
}
