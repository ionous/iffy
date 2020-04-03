package qna

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/tables"
)

// Fields populates itself from the database on demand.
type Fields struct {
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

func NewObjectValues(db *sql.DB) *Fields {
	return &Fields{make(mapType), tables.NewCache(db)}
}

// other possibilities as needed:
// get aspect,
// get trait, -- possibly we make these explicit iffy commands initially, then select the correct command based on type.
// get class,

// GetValue sets the value of the passed pointer to the value of the named property.
func (n *Fields) GetField(obj, field string, pv interface{}) (err error) {
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
			rows = n.db.QueryRow(
				`select kind || ( case path when '' then ('') else ("," || path) end ) as path
				from mdl_noun mn 
				join mdl_kind mk 
					using (kind)
				where noun=?`,
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
func (n *Fields) SetField(obj, field string, v interface{}) (err error) {
	key := keyType{obj, field}
	n.pairs[key] = v
	return
}

var notImplemented = errutil.New("not implemented")
