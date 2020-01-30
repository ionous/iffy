package assembly

import (
	"database/sql"
)

// goal: build table of start_val(noun, field, value) for instances.
// considerations:
// . property's actual kind ( default specified against a derived type )
// . contradiction in specified values
// . contradiction in specified value vs field type ( alt: implicit conversion )
// . missing properties ( kind, field pair doesn't exist in model )
// o certainties: usually, seldom, never, always.
// o misspellings, near spellings ( ex. for missing fields )
func DetermineValues(m *Modeler, db *sql.DB) (err error) {
	if e := determineInitialFields(m, db); e != nil {
		err = e
	} else if e := determineInitialTraits(m, db); e != nil {
		err = e
	}
	return
}
