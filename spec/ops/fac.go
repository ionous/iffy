package ops

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
)

// Factory implements spec.SpecFactory.
type Factory struct {
	cmds  *Ops
	xform Transform
}

func NewFactory(cmds *Ops, xform Transform) *Factory {
	return &Factory{cmds, xform}
}

func (fac *Factory) CmdFromName(name string) (ret *Command, err error) {
	if rtype, ok := fac.cmds.FindType(name); ok {
		ret = &Command{
			xform:  fac.xform,
			target: NewTarget(rtype),
		}
	} else if rtype, ok := fac.cmds.ShadowTypes.FindType(name); ok {
		ret = &Command{
			xform:  fac.xform,
			target: Shadow(rtype),
		}
	} else {
		err = errutil.New("unknown command", name)
	}
	return
}

func (fac *Factory) CmdFromPointer(ptr interface{}) (*Command, error) {
	cmd := &Command{xform: fac.xform, target: InPlace(ptr)}
	return cmd, nil
}

// NewSpec implements spec.SpecFactory.
func (fac *Factory) NewSpec(name string) (spec.Spec, error) {
	cmd, err := fac.CmdFromName(name)
	return cmd, err
}

// NewSpecs implements spec.SpecFactory.
// the c algorithm creates NewSpecs, and then assigns it to a slot
// we need the slot to target the array properly, so we just wait,
func (fac *Factory) NewSpecs() (spec.Specs, error) {
	return &Commands{xform: fac.xform}, nil
}
