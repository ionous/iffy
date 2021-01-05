package scope

import (
	g "github.com/ionous/iffy/rt/generic"
)

type TargetRecord struct {
	Target string
	Record *g.Record
}

func (k *TargetRecord) GetField(target, field string) (ret g.Value, err error) {
	if target != k.Target {
		err = g.UnknownTarget{target}
	} else {
		ret, err = k.Record.GetNamedField(field)
	}
	return
}

// SetField writes the value of 'v' into the value at 'name'.
func (k *TargetRecord) SetField(target, field string, v g.Value) (err error) {
	if target != k.Target {
		err = g.UnknownTarget{target}
	} else {
		err = k.Record.SetNamedField(field, v)
	}
	return
}
