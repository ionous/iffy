package assembly

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// goal: build table of mdl_default(kind,field,value) for archetypes.
// considerations:
// . property's actual kind ( default specified against a derived type )
// . contradiction in specified values
// . contradiction in specified value vs field type ( alt: implicit conversion )
// . missing properties ( kind, field pair doesn't exist in model )
// o certainties: usually, seldom, never, always.
// o misspellings, near spellings ( ex. for missing fields )
func DetermineDefaults(m *Modeler, db *sql.DB) (err error) {
	if e := determineDefaultFields(m, db); e != nil {
		err = e
	} else if e := determineDefaultTraits(m, db); e != nil {
		err = e
	}
	return
}
