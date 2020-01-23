package assembly

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	_ "github.com/mattn/go-sqlite3"
)

// goal: build table of mdl_field,default value
// uses: eph_named_default(kind, field, value): for the user's requested defaults.
//       mdl_kind(kind, path): for hierarchy.
//       mdl_field(kind, field, type)
// considerations:
// . property's actual kind ( default specified against a derived type )
// . contradiction in specified values
// . contradiction in specified value vs field type ( alt: implicit conversion )
// . missing properties ( kind, field pair doesn't exist in model )
// o certainties: usually, seldom, never, always.
// o misspellings, near spellings ( ex. for missing fields )
func DetermineDefaults(m *Modeler, db *sql.DB) (err error) {
	// grab out the *actual* kind,field pairings
	var fieldType string
	var curr defaultInfo
	var list []defaultInfo
	// collect mdl_field rowid, type, and value.
	if e := dbutil.QueryAll(db, `
	select ep.idModelField, mf.type, ed.value 
	from eph_modeled_default ep 
		join mdl_field mf 
		on mf.rowid= ep.idModelField
	left join eph_default ed
		on ed.rowid = ep.idEphDefault`,
		func() (err error) {
			list = append(list, curr)
			return
		},
		&curr.idModelField, &fieldType, &curr.value,
	); e != nil {
		err = e
	} else {
		for _, n := range list {
			if e := m.WriteDefault(n.idModelField, n.value); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}
