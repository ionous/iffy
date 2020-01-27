package assembly

import (
	"database/sql"
	"reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	_ "github.com/mattn/go-sqlite3"
)

// goal: build table of mdl_value(noun, field, value) for instances.
// considerations:
// . property's actual kind ( default specified against a derived type )
// . contradiction in specified values
// . contradiction in specified value vs field type ( alt: implicit conversion )
// . missing properties ( kind, field pair doesn't exist in model )
// o certainties: usually, seldom, never, always.
// o misspellings, near spellings ( ex. for missing fields )
func DetermineValues(m *Modeler, db *sql.DB) (err error) {
	var store defaultValueStore
	var curr, last defaultValue
	if e := dbutil.QueryAll(db,
		`select at.target, mf.field, mf.type, ev.value
	from asm_value_tree at
		join mdl_field mf
		on mf.rowid = at.idModelField
	left join eph_value ev
		on ev.rowid = at.idEphValue
	order by at.target, mf.field`,
		func() (err error) {
			// if the modelField is the same, so is kind, field, type.
			if nv, e := convertField(curr.fieldType, curr.value); e != nil {
				err = e
			} else if !last.isValid() {
				last, last.value = curr, nv
			} else if last.target != curr.target || last.field != curr.field {
				store.add(last)
				last, last.value = curr, nv
			} else if !reflect.DeepEqual(last.value, nv) {
				err = errutil.Fmt("conflicting values: %s != %s:%T", last.String(), curr.String(), nv)
			}
			return
		},
		&curr.target, &curr.field, &curr.fieldType,
		&curr.value,
	); e != nil {
		err = e
	} else {
		store.add(last)
		err = store.writeValues(m)
	}
	return
}
