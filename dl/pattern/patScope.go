package pattern

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/object"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
	"github.com/ionous/iffy/rt/scope"
)

type patScope struct {
	ret string // name of return variable within the target record
	scope.TargetRecord
}

func newScope(ret string, rec *g.Record) *patScope {
	return &patScope{ret, scope.TargetRecord{object.Variables, rec}}
}

func (v *patScope) GetValue(aff affine.Affinity) (ret g.Value, err error) {
	if n := v.ret; len(n) > 0 {
		if res, e := v.Record.GetNamedField(n); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else if e := safe.Check(res, aff); e != nil {
			err = errutil.New("error trying to get return value", e)
		} else {
			ret = res
		}
	} else if len(aff) != 0 {
		err = errutil.New("caller expected", aff, "returned nothing")
	}
	return
}
