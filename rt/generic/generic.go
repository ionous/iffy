package generic

import (
	"github.com/ionous/iffy/rt"
)

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
