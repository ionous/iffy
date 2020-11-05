package generic

import (
	"github.com/ionous/iffy/affine"
)

type Bool struct {
	Nothing
	val bool
}

type Float struct {
	Nothing
	val float64
}

type Int struct {
	Nothing
	val int
}

type String struct {
	Nothing
	val string
}

func NewBool(v bool) *Bool {
	return &Bool{val: v}
}
func NewFloat(v float64) *Float {
	return &Float{val: v}
}
func NewInt(v int) *Int {
	return &Int{val: v}
}
func NewString(v string) *String {
	return &String{val: v}
}

func (n *Bool) Affinity() affine.Affinity { return affine.Bool }
func (n *Bool) Type() string              { return "bool" }
func (n *Bool) GetBool() (ret bool, err error) {
	ret = n.val
	return
}

func (n *Float) Affinity() affine.Affinity { return affine.Number }
func (n *Float) Type() string              { return "float64" }
func (n *Float) GetNumber() (ret float64, err error) {
	ret = n.val
	return
}

func (n *Int) Affinity() affine.Affinity { return affine.Number }
func (n *Int) Type() string              { return "int" }
func (n *Int) GetNumber() (ret float64, err error) {
	ret = float64(n.val)
	return
}

func (n *String) Affinity() affine.Affinity { return affine.Text }
func (n *String) Type() string              { return "string" }
func (n *String) GetText() (ret string, err error) {
	ret = n.val
	return
}
