package ephemera

import (
	"database/sql"

	"github.com/ionous/iffy/dbutil"
)

// KidsOf the passed ancestor as specified in the Kinds table.
// note: this is raw user data, there may be ambiguities or conflicts in the pairings.
func KidsOf(db *sql.DB, ancestor string, cb func(kid string)) (err error) {
	// query the child names of the parent named p
	if siblings, e := db.Query(
		`select distinct n.name as kid 
					from eph_named n,( select idNamedKind
								from eph_kind,( select rowid as parentId 
										    from eph_named n
										    where n.name = ? )
								where idNamedParent = parentId )
					where n.rowid = idNamedKind`, ancestor); e != nil {
		err = e
	} else {
		var kid string
		err = dbutil.ScanAll(siblings, func() (err error) {
			cb(kid)
			return
		}, &kid)
	}
	return
}
