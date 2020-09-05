package assembly

import (
	"github.com/ionous/iffy/check"
	"github.com/ionous/iffy/dl/core"
)

func AssembleTests(asm *Assembler) (err error) {
	// todo: doesn't build mdl check
	// doesnt check for conflicts or errors in test definitions
	var name, expect string
	var prog []byte
	rule := BuildRule{
		Query: `select name, prog, expect
		from asm_check ek
		join asm_expect ex
			using (name)
		join eph_prog ep
			on (ek.idProg = ep.rowid)
		order by name, progType, idProg`,
		NewContainer: func(name string) interface{} {
			return &check.CheckOutput{
				Name:   name,
				Expect: expect,
				Test:   &core.Activity{},
			}
		},
		NewEl: func(c interface{}) interface{} {
			curr := c.(*check.CheckOutput)
			curr.Test.Exe = append(curr.Test.Exe, nil)
			return &curr.Test.Exe[len(curr.Test.Exe)-1]
		},
	}
	return rule.buildFromRule(asm, &name, &prog, &expect)
}
