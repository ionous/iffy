package assembly

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

func copyRules(db *sql.DB) error {
	_, e := db.Exec(
		`insert into mdl_rule select pattern, idProg from asm_rule`)
	return e
}

func checkRuleSetup(db *sql.DB) (err error) {
	var pat, patType, ruleType string

	// mismatched
	if e := tables.QueryAll(db,
		`select distinct pattern, pt, rt
		from asm_rule_match am
		where am.matched = 0
		order by pattern, pt, rt;`,
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
