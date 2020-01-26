package assembly

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	"github.com/ionous/iffy/ephemera"
	_ "github.com/mattn/go-sqlite3"
)

// goal: build table of mdl_default(mdl_field,value) for kinds
// uses: asm_default(kind, field, value): for the user's requested defaults.
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
	var store defaultStore
	var curr, last defaultInfo
	if e := dbutil.QueryAll(db,
		`select at.target, mf.field, mf.type, ed.value 
	from asm_default_tree at 
		join mdl_field mf 
		on mf.rowid= at.idModelField
	left join eph_default ed
		on ed.rowid = at.idEphDefault
	order by at.target, mf.field`,
		func() (err error) {
			if nv, e := convertField(curr.fieldType, curr.value); e != nil {
				err = e
			} else if !last.isValid() {
				last, last.value = curr, nv
			} else if last.kind != curr.kind || last.field != curr.field {
				store.add(last)
				last, last.value = curr, nv
			} else if !reflect.DeepEqual(last.value, nv) {
				err = errutil.Fmt("conflicting defaults: %s != %s:%T", last.String(), curr.String(), nv)
			}
			return
		},
		&curr.kind, &curr.field, &curr.fieldType,
		&curr.value,
	); e != nil {
		err = e
	} else {
		store.add(last)
		err = store.write(m)
	}
	return
}

type defaultInfo struct {
	kind, field, fieldType string
	value                  interface{}
}

func (n *defaultInfo) isValid() bool {
	return len(n.field) > 0
}

func (n *defaultInfo) String() string {
	return n.kind + "." + n.field + ":" + n.fieldType + fmt.Sprintf("(%v:%T)", n.value, n.value)
}

type defaultStore struct {
	list []defaultInfo
}

func (store *defaultStore) add(n defaultInfo) {
	store.list = append(store.list, n)
}

func (store *defaultStore) write(m *Modeler) (err error) {
	for _, n := range store.list {
		if e := m.WriteDefault(n.kind, n.field, n.value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// out types are currently: int, float32, or string.
func convertField(fieldType string, value interface{}) (ret interface{}, err error) {
	switch v := reflect.ValueOf(value); fieldType {
	case ephemera.PRIM_DIGI:
		switch k := v.Kind(); k {
		case reflect.Float64:
			ret = float32(v.Float())
		case reflect.Int64:
			ret = int(v.Int())
		default:
			err = errutil.New("can't convert from", k, "to int")
		}
	case ephemera.PRIM_TEXT:
		switch k := v.Kind(); k {
		case reflect.String:
			ret = v.String()
		default:
			err = errutil.New("can't convert from", k, "to string")
		}
	default:
		err = errutil.New("unhandled field type", fieldType)
	}
	return
}
