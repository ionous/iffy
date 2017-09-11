package rule

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/spec"
	"github.com/ionous/iffy/spec/ops"
)

func Master(cmds *ops.Ops, xform ops.Transform, pt *unique.Stack, buildPatterns ...func(c spec.Block)) (ret Rules, err error) {
	// Accumulate rules into root.
	var root struct{ Mandates }

	if c, ok := cmds.NewXBuilder(&root, xform); !ok {
		err = errutil.New("why does this return okay anyway?")
	} else if c.Param("mandates").Cmds().Begin() {
		for _, b := range buildPatterns {
			b(c)
		}
		c.End()
		// Execute the accumulated pattern definitions
		if e := c.Build(); e != nil {
			err = e
		} else {
			rules := MakeRules()
			if e := root.Mandate(pt.Types, rules); e != nil {
				err = e
			} else {
				ret = rules
			}
		}
	}
	return
}
