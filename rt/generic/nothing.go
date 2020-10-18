package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

type Nothing struct{}

func (n Nothing) GetBool(rt.Runtime) (_ bool, err error) {
	err = errutil.New("value is not a bool")
	return
}
func (n Nothing) GetNumber(rt.Runtime) (_ float64, err error) {
	err = errutil.New("value is not a number")
	return
}
func (n Nothing) GetText(rt.Runtime) (_ string, err error) {
	err = errutil.New("value is not a text")
	return
}
func (n Nothing) GetNumList(rt.Runtime) (_ []float64, err error) {
	err = errutil.New("value is not a number list")
	return
}
func (n Nothing) GetTextList(rt.Runtime) (_ []string, err error) {
	err = errutil.New("value is not a text list")
	return
}
func (n Nothing) GetLen(rt.Runtime) (_ int, err error) {
	err = errutil.New("value is not measurable")
	return
}
func (n Nothing) GetIndex(rt.Runtime, int) (_ rt.Value, err error) {
	err = errutil.New("value is not indexable")
	return
}
func (n Nothing) GetFieldByName(rt.Runtime, string) (_ rt.Value, err error) {
	err = errutil.New("value is not an object")
	return
}
