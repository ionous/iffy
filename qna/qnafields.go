package qna

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

// Values populates itself from the database on demand.
type Values struct {
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

func NewValues(db *sql.DB) *Values {
	return &Values{make(mapType), tables.NewCache(db)}
}

// other possibilities as needed:
// get aspect,
// get trait, -- possibly we make these explicit iffy commands initially, then select the correct command based on type.
// get class,

// GetValue sets the value of the passed pointer to the value of the named property.
func (n *Values) GetField(obj, field string, pv interface{}) (err error) {
	key := keyType{obj, field}
	if v, ok := n.pairs[key]; ok {
		err = Assign(pv, v)
	} else {
		tgt := mapTarget{key: key, pairs: n.pairs}
		if e := n.db.QueryRow("select value from run_init where noun=? and field=? order by tier limit 1",
			obj, field).Scan(&tgt); e == nil {
			err = Assign(pv, tgt.value)
			//
		} else if e == sql.ErrNoRows {
			n.pairs[key] = nil
			err = Assign(pv, nil)
			// option 1: generate the zero value in sql, somehow
			// option 2: reflect -- mattn already uses reflect
			// option 3: scan to an interface? and have a generic/switch unpack
			// option 4: assume pv is zero already.
			// option 5: store nil and use Assign below to create/set zero
		} else {
			err = e
		}
	}
	return err
}

// Assign sets the named property to the passed value.
func (n *Values) SetField(obj, field string, v interface{}) (err error) {
	key := keyType{obj, field}
	n.pairs[key] = v
	return
}

var notImplemented = errutil.New("not implemented")
