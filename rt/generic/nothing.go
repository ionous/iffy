package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

type Nothing struct{}

func (n Nothing) Affinity() affine.Affinity {
	return ""
}
func (n Nothing) Type() string {
	return ""
}
func (n Nothing) GetBool() (_ bool, err error) {
	err = errutil.New("value is not a bool")
	return
}
func (n Nothing) GetNumber() (_ float64, err error) {
	err = errutil.New("value is not a number")
	return
}
func (n Nothing) GetText() (_ string, err error) {
	err = errutil.New("value is not a text")
	return
}
func (n Nothing) GetNumList() (_ []float64, err error) {
	err = errutil.New("value is not a number list")
	return
}
func (n Nothing) GetTextList() (_ []string, err error) {
	err = errutil.New("value is not a text list")
	return
}
func (n Nothing) GetLen() (_ int, err error) {
	err = errutil.New("value is not measurable")
	return
}
func (n Nothing) GetIndex(int) (_ rt.Value, err error) {
	err = errutil.New("value is not indexable")
	return
}
func (n Nothing) GetField(string) (_ rt.Value, err error) {
	err = errutil.New("value is not an object")
	return
}
func (n Nothing) SetField(string, rt.Value) (err error) {
	err = errutil.New("value is not writable")
	return
}
