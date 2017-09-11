package ops

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
)

type Commands struct {
	xform Transform
	els   []*Command
}

func (cs *Commands) AddElement(el spec.Spec) (err error) {
	if c, ok := el.(*Command); !ok {
		err = errutil.Fmt("unexpected element type %T", el)
	} else {
		cs.els = append(cs.els, c)
	}
	return
}
