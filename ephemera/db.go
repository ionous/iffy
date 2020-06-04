package ephemera

import (
	"database/sql"

	"github.com/ionous/iffy/tables"
)

// KidsOf the passed ancestor as specified in the Kinds table.
// note: this is raw user data, there may be ambiguities or conflicts in the pairings.
func KidsOf(db *sql.DB, pluralAncestor string, cb func(kid string)) (err error) {
	// ex. from cats are a kind of animal, extract "cats" when looking for "animals"
	if siblings, e := db.Query(
		`select distinct kid from asm_ancestry where parent = ?`, pluralAncestor); e != nil {
		err = e
	} else {
		var kid string
		err = tables.ScanAll(siblings, func() (err error) {
			cb(kid)
			return
		}, &kid)
	}
	return
}
