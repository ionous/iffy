package ephemera

import (
	"database/sql"

	"github.com/ionous/errutil"
)

// KidsOf the passed ancestor as specified in the Kinds table.
// note: this is raw user data, there may be ambiguities or conflicts in the pairings.
func KidsOf(db *sql.DB, ancestor string, cb func(kid string)) (err error) {
	// query the child names of the parent named p
	if siblings, e := db.Query(
		`select distinct named.name as kid 
					from named,( select idNamedKind
								from kind,( select rowid as parentId 
										    from named 
										    where named.name = ? )
								where idNamedParent = parentId )
					where named.rowid = idNamedKind`, ancestor); e != nil {
		err = e
	} else {
		for cnt := 0; siblings.Next(); cnt++ {
			var kid string
			if e := siblings.Scan(&kid); e != nil {
				err = e
				break
			} else {
				cb(kid)
			}
		}
		// tests if early exit
		if e := siblings.Err(); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
