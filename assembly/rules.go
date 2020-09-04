package assembly

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

func buildRules(asm *Assembler) (err error) {
	for _, rule := range rules {
		if e := buildFromRule(asm, &rule); e != nil {
			err = e
			break
		}
	}
	return
}

func buildFromRule(asm *Assembler, b *BuildRule) (err error) {
	list := make(map[string]interface{})
	var name, last string
	var prog []byte
	var curr interface{}
	if e := tables.QueryAll(asm.cache.DB(), b.Query,
		func() (err error) {
			if name != last || curr == nil {
				curr = b.NewContainer(name)
				list[name] = curr
				last = name
			}
			el := b.NewEl(curr)
			dec := gob.NewDecoder(bytes.NewBuffer(prog))
			return dec.Decode(el)
		}, &name, &prog); e != nil {
		err = errutil.New("buildRules", e)
	} else {
		err = asm.WriteGobs(list)
	}
	return
}

func checkRuleSetup(db *sql.DB) (err error) {
	var pat, patType, ruleType string
	// mismatched
	if e := tables.QueryAll(db,
		`select distinct pattern, pt, rt
		from asm_rule_match am
		where am.matched = 0
		order by pattern, pt, rt`,
		func() error {
			e := errutil.Fmt("pattern %q has type %q but rule %q",
				pat, patType, ruleType)
			err = errutil.Append(err, e)
			return nil // accumulate
		},
		&pat, &patType, &ruleType); e != nil {
		err = e
	} else {
		// missing
		if e := tables.QueryAll(db,
			`select distinct pattern
		from asm_pattern
		where decl = 1
		and pattern not in 
		(select pattern from asm_rule_match am
			where am.matched=1)`,
			func() error {
				e := errutil.Fmt("pattern %q has no valid rules", pat)
				err = errutil.Append(err, e)
				return nil // accumulate
			},
			&pat); e != nil {
			err = e
		}
	}
	return
}
