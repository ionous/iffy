package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type Nothing struct{}

type Bool struct {
	Nothing
	Value bool
}

type Int struct {
	Nothing
	Value int
}

type Float struct {
	Nothing
	Value float64
}

type String struct {
	Nothing
	Value string
}

type FloatList struct {
	Nothing
	Value []float64
}

type StringList struct {
	Nothing
	Value []string
}

func (n Nothing) GetBool(rt.Runtime) (ret bool, err error) {
	err = errutil.New("value is not a bool")
	return
}
func (n Nothing) GetNumber(rt.Runtime) (ret float64, err error) {
	err = errutil.New("value is not a number")
	return
}
func (n Nothing) GetText(rt.Runtime) (ret string, err error) {
	err = errutil.New("value is not a text")
	return
}
func (n Nothing) GetNumberStream(rt.Runtime) (ret rt.Iterator, err error) {
	err = errutil.New("value is not a number stream")
	return
}
func (n Nothing) GetTextStream(rt.Runtime) (ret rt.Iterator, err error) {
	err = errutil.New("value is not a text stream")
	return
}

//
func (n *Bool) GetBool(rt.Runtime) (ret bool, err error) {
	ret = n.Value
	return
}
func (n *Float) GetNumber(rt.Runtime) (ret float64, err error) {
	ret = n.Value
	return
}
func (n *Int) GetNumber(rt.Runtime) (ret float64, err error) {
	ret = float64(n.Value)
	return
}
func (n *String) GetText(rt.Runtime) (ret string, err error) {
	ret = n.Value
	return
}
func (n *FloatList) GetNumberStream(rt.Runtime) (ret rt.Iterator, err error) {
	ret = SliceFloats(n.Value)
	return
}
func (n *StringList) GetTextStream(rt.Runtime) (ret rt.Iterator, err error) {
	ret = SliceStrings(n.Value)
	return
}
