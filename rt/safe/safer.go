package safe

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

func Check(v g.Value, want affine.Affinity) (err error) {
	if va := v.Affinity(); len(want) > 0 && want != va {
		err = errutil.Fmt("have %q, wanted %q", va, want)
	}
	return
}

// resolve a requested variable name into a value of the desired affinity.
func CheckVariable(run rt.Runtime, n string, aff affine.Affinity) (ret g.Value, err error) {
	if v, e := run.GetField(object.Variables, n); e != nil {
		err = e
	} else if e := Check(v, aff); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}

func Unpack(src g.Value, field string, aff affine.Affinity) (ret g.Value, err error) {
	if !affine.HasFields(src.Affinity()) {
		err = errutil.New("Value", src, "doesn't have fields")
	} else if v, e := src.FieldByName(field); e != nil {
		err = e
	} else if e := Check(v, aff); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}
