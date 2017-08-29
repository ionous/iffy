package rule

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/spec/ops"
)

func Master(cmds *ops.Ops, pt *unique.Stack, buildPatterns ...func(c *ops.Builder)) (ret Rules, err error) {
	// Accumulate rules into root.
	var root struct{ Mandates }

	if c, ok := cmds.NewBuilder(&root); !ok {
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
				rules.Sort()
				ret = rules
			}
		}
	}
	return
}
