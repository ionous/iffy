package generic

import (
	"github.com/ionous/iffy/affine"
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

func (n *Bool) Affinity() affine.Affinity { return affine.Bool }
func (n *Bool) Type() string              { return "bool" }
func (n *Bool) GetBool() (ret bool, err error) {
	ret = n.Value
	return
}

func (n *Float) Affinity() affine.Affinity { return affine.Number }
func (n *Float) Type() string              { return "float64" }
func (n *Float) GetNumber() (ret float64, err error) {
	ret = n.Value
	return
}

func (n *Int) Affinity() affine.Affinity { return affine.Number }
func (n *Int) Type() string              { return "int" }
func (n *Int) GetNumber() (ret float64, err error) {
	ret = float64(n.Value)
	return
}

func (n *String) Affinity() affine.Affinity { return affine.Text }
func (n *String) Type() string              { return "string" }
func (n *String) GetText() (ret string, err error) {
	ret = n.Value
	return
}
