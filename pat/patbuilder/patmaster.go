package patbuilder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/pat/patspec"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/spec/ops"
)

type Pm struct {
	ops      *ops.Ops
	patterns *Patterns
	err      error
}

func NewPatternMaster(ops *ops.Ops, classes ref.Classes, block ...interface{}) Pm {
	patterns := NewPatterns(classes)
	err := unique.RegisterBlocks(patterns, block...)
	return Pm{ops, patterns, err}
}

// FIX: too many steps here. look at simplifying the process a bit.
func (pm Pm) Build(buildPatterns ...func(c *ops.Builder)) (ret *Patterns, err error) {
	if pm.err != nil {
		err = pm.err
	} else {
		// Accumulate patterns into root.
		var root struct {
			Patterns patspec.PatternSpecs
		}

		if c, ok := pm.ops.NewBuilder(&root); !ok {
			err = errutil.New("why does this return okay anyway?")
		} else if c.Param("patterns").Cmds().Begin() {
			for _, b := range buildPatterns {
				b(c)
			}
			c.End()
			// Execute the accumulated pattern definitions
			if e := c.Build(); e != nil {
				err = e
			} else if e := root.Patterns.Generate(pm.patterns); e != nil {
				err = e
			} else {
				ret = pm.patterns
			}
		}
	}
	return
}
