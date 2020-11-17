package rt

import (
	"github.com/ionous/iffy/object"
	g "github.com/ionous/iffy/rt/generic"
)

// eliminates some boilerplate code when working with runtime variables
type Variable string

func (v Variable) GetValue(run Runtime) (g.Value, error) {
	return run.GetField(object.Variables, string(v))
}

func (v Variable) SetValue(run Runtime, val g.Value) error {
	return run.SetField(object.Variables, string(v), val)
}

func (v Variable) GetBool(run Runtime) (ret bool, err error) {
	if val, e := v.GetValue(run); e != nil {
		err = e
	} else {
		ret, err = val.GetBool()
	}
	return
}

func (v Variable) GetNumber(run Runtime) (ret float64, err error) {
	if val, e := v.GetValue(run); e != nil {
		err = e
	} else {
		ret, err = val.GetNumber()
	}
	return
}

func (v Variable) GetText(run Runtime) (ret string, err error) {
	if val, e := v.GetValue(run); e != nil {
		err = e
	} else {
		ret, err = val.GetText()
	}
	return
}

func (v Variable) GetRecord(run Runtime) (ret *g.Record, err error) {
	if val, e := v.GetValue(run); e != nil {
		err = e
	} else {
		ret, err = val.GetRecord()
	}
	return
}

func (v Variable) GetTextList(run Runtime) (ret []string, err error) {
	if val, e := v.GetValue(run); e != nil {
		err = e
	} else {
		ret, err = val.GetTextList()
	}
	return
}

func (v Variable) GetNumList(run Runtime) (ret []float64, err error) {
	if val, e := v.GetValue(run); e != nil {
		err = e
	} else {
		ret, err = val.GetNumList()
	}
	return
}

func (v Variable) GetRecordList(run Runtime) (retType string, retList []*g.Record, err error) {
	if val, e := v.GetValue(run); e != nil {
		err = e
	} else if vs, e := val.GetRecordList(); e != nil {
		err = e
	} else {
		retType, retList = val.Type(), vs
	}
	return
}

// returns the named object instead of the named variable
func (v Variable) GetObjectByName(run Runtime) (ret g.Value, err error) {
	switch val, e := run.GetField(object.Value, string(v)); e.(type) {
	case g.UnknownField:
		err = g.UnknownObject(string(v))
	default:
		ret, err = val, e
	}
	return
}

// first look for a variable named "name" in scope, unbox it (if need be) to return the object's id.
func (v Variable) GetObjectByVariable(run Runtime) (ret g.Value, err error) {
	switch val, e := v.GetValue(run); e.(type) {
	default:
		err = e
	// if there's no such variable, check if there's an object of that name.
	case g.UnknownTarget, g.UnknownField:
		ret, err = v.GetObjectByName(run)
	case nil:
		ret = val
	}
	return
}
