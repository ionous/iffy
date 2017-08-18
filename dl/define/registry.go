package define

import (
	// "github.com/ionous/iffy/dl/std"
	// "github.com/ionous/iffy/pat/patbuilder"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/spec/ops"
)

var globalRegistry Registry

type Registry struct {
	callbacks []func(*ops.Builder)
}

// Register definitions globally. Used mainly via go init()
func (r *Registry) Register(cb func(c *ops.Builder)) {
	r.callbacks = append(r.callbacks, cb)
}

// Register definitions globally. Used mainly via go init()
func Register(cb func(c *ops.Builder)) {
	globalRegistry.Register(cb)
}

// Define implements Statement by using all Register(ed) definitions.
func (r Registry) Define(f *Facts) (err error) {
	classes := ref.NewClasses()
	unique.RegisterBlocks(unique.PanicTypes(classes),
		(*Classes)(nil),
	)

	// objects := ref.NewObjects(classes)
	// unique.RegisterValues(unique.PanicValues(objects),
	// 	Thingaverse.objects(sliceOf.String("apple", "pen", "thing#1", "thing#2"))...)

	cmds := ops.NewOps()
	unique.RegisterBlocks(unique.PanicTypes(cmds),
		(*Commands)(nil),
	)

	unique.RegisterBlocks(unique.PanicTypes(cmds.ShadowTypes),
		(*Patterns)(nil),
	)

	// if patterns, e := patbuilder.NewPatternMaster(cmds, classes,
	// 	(*Patterns)(nil)).Build(std.StdPatterns); e != nil {
	// 	err = e
	// } else {

	var root struct{ Definitions }
	if c, ok := cmds.NewBuilder(&root); ok {
		if c.Cmds().Begin() {
			for _, v := range r.callbacks {
				v(c)
			}
			c.End()
		}
		if e := c.Build(); e != nil {
			err = e
		} else {
			err = root.Define(f)
		}
	}
	return
}
