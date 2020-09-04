package assembly

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/pattern"
	"github.com/ionous/iffy/tables"
)

func buildRules(asm *Assembler) (err error) {
	var name, lastName, typeName string
	var rule pattern.TextRule
	//
	var pats []*pattern.TextPattern
	gs := tables.NewGobScanner(&rule)
	if e := tables.QueryAll(asm.cache.DB(),
		`select pattern, type, prog
		from asm_rule
		where type='text_rule'`,
		func() (err error) {
			// new pattern
			if name != lastName {
				pats = append(pats, &pattern.TextPattern{
					Name: name,
				})
				lastName = name
			}
			copy := rule             // new rule
			pat := pats[len(pats)-1] // last
			pat.Rules = append(pat.Rules, &copy)
			return
		},
		&name, &typeName, gs); e != nil {
		err = errutil.New("buildRules", e)
	}
	for _, pat := range pats {
		if _, e := asm.WriteGob(pat.Name, pat); e != nil {
			errutil.Append(err, e)
		}
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
