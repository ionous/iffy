package assembly

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	"github.com/ionous/iffy/ephemera"
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
	var store defaultStore
	var curr, last defaultInfo
	if e := dbutil.QueryAll(db, `
	select ep.idModelField, mf.kind, mf.field, mf.type, ed.value 
	from eph_modeled_default ep 
		join mdl_field mf 
		on mf.rowid= ep.idModelField
	left join eph_default ed
		on ed.rowid = ep.idEphDefault
	order by ep.idModelField`,
		func() (err error) {
			// if the modelField is the same, so is kind, field, type.
			if last.idModelField != curr.idModelField {
				if e := store.add(&curr); e != nil {
					err = e
				} else {
					last = curr
				}
			} else if newValue, e := convertField(last.fieldType, curr.value); e != nil {
				err = e
			} else if !reflect.DeepEqual(last.value, newValue) {
				err = errutil.New("conflicting defaults", last, "!=", newValue, reflect.TypeOf(last.value), "!=", reflect.TypeOf(newValue))
			}
			return
		},
		&curr.idModelField,
		&curr.kind, &curr.field, &curr.fieldType,
		&curr.value,
	); e != nil {
		err = e
	} else if e := store.add(&last); e != nil {
		err = e
	} else {
		err = store.write(m)
	}
	return
}

type defaultInfo struct {
	idModelField           int64
	kind, field, fieldType string
	value                  interface{}
}

func (n *defaultInfo) String() string {
	return strings.Join([]string{n.kind, n.field, n.fieldType, fmt.Sprint(n.value)}, " ")
}

type defaultStore struct {
	list []defaultInfo
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

func (store *defaultStore) add(n *defaultInfo) (err error) {
	if n.idModelField > 0 {
		if v, e := convertField(n.fieldType, n.value); e != nil {
			err = e
		} else {
			n.value = v
			store.list = append(store.list, *n)
		}
	}
	return
}

func (store *defaultStore) write(m *Modeler) (err error) {
	for _, n := range store.list {
		if e := m.WriteDefault(n.idModelField, n.value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
