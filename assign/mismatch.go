package assign

import "github.com/ionous/errutil"

func Mismatch(name string, want, have interface{}) error {
	return mismatch{name, want, have}
}

type mismatch struct {
	Name       string
	Want, Have interface{}
}

func (m mismatch) Error() string {
	return errutil.Fmt("while %s expected %T, it has %T", m.Name, m.Want, m.Have).Error()
}
