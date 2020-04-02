package qna

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/tables"
)

// ObjectValues populates itself from the database on demand.
type ObjectValues struct {
	pairs mapType
	db    *tables.Cache
}

type keyType struct {
	owner, member string
}

type mapType map[keyType]interface{}

type mapTarget struct {
	key   keyType
	pairs mapType
	value interface{}
}

func (k *mapTarget) Scan(v interface{}) (err error) {
	// bytes will need special processing ( copies )
	k.pairs[k.key], k.value = v, v
	return
}

func NewObjectValues(db *sql.DB) *ObjectValues {
	return &ObjectValues{make(mapType), tables.NewCache(db)}
}

// other possibilities as needed:
// get aspect,
// get trait, -- possibly we make these explicit iffy commands initially, then select the correct command based on type.
// get class,

// GetValue sets the value of the passed pointer to the value of the named property.
func (n *ObjectValues) GetObject(obj, field string, pv interface{}) (err error) {
	key := keyType{obj, field}
	if v, ok := n.pairs[key]; ok {
		err = Assign(pv, v)
	} else {
		var permissive bool
		tgt := mapTarget{key: key, pairs: n.pairs}
		var rows tables.RowScanner
		switch field {
		case object.Kind:
			rows = n.db.QueryRow("select kind from mdl_noun where noun=?",
				obj)
		case object.Kinds:
			// objects and kinds are distinct namespaces
			// so we can reuse the object property cache to cache kind info
			// alternatively, we could give each object its path...
			// and that might be a little bit nicer.
			rows = n.db.QueryRow("select path from mdl_kind where kind=?",
				obj)
		case object.Exists:
			rows = n.db.QueryRow("select count() from mdl_noun where noun=?",
				obj)
		default:
			// FIX? needs more work to determine if the field really exists
			// ex. possibly a union query of class field with a nil value
			permissive = true
			rows = n.db.QueryRow("select value from run_init where noun=? and field=? order by tier limit 1",
				obj, field)
		}
		if e := rows.Scan(&tgt); e == nil {
			err = Assign(pv, tgt.value)
			//
		} else if e == sql.ErrNoRows {
			if !permissive {
				err = errutil.New("field not found", obj, field)
			} else {
				n.pairs[key] = nil
				err = Assign(pv, nil)
			}
		} else {
			err = e
		}
	}

	return err
}

// Assign sets the named property to the passed value.
func (n *ObjectValues) SetObject(obj, field string, v interface{}) (err error) {
	key := keyType{obj, field}
	n.pairs[key] = v
	return
}

var notImplemented = errutil.New("not implemented")
