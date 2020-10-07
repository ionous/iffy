package rt

import "github.com/ionous/iffy/object"

// eliminates some boilerplate code when working with runtime variables
type Variable string

func (v Variable) GetBool(run Runtime) (ret bool, err error) {
	if val, e := run.GetField(object.Variables, string(v)); e != nil {
		err = e
	} else {
		ret, err = val.GetBool(run)
	}
	return
}

func (v Variable) GetNumber(run Runtime) (ret float64, err error) {
	if val, e := run.GetField(object.Variables, string(v)); e != nil {
		err = e
	} else {
		ret, err = val.GetNumber(run)
	}
	return
}

func (v Variable) GetText(run Runtime) (ret string, err error) {
	if val, e := run.GetField(object.Variables, string(v)); e != nil {
		err = e
	} else {
		ret, err = val.GetText(run)
	}
	return
}
