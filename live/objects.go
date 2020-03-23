package live

import (
	"database/sql"
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

type Values struct {
	pairs pairs
	db    *tables.Cache
}

type pair struct {
	target, member string
}

type pairs map[pair]interface{}

func NewValues(db *sql.DB) *Values {
	return &Values{make(pairs), tables.NewCache(db)}
}

// other possibilities as needed:
// get aspect,
// get trait
// get class,

// GetValue sets the value of the passed pointer to the value of the named property.
func (n *Values) GetField(noun, field string, pv interface{}) (err error) {
	pair := pair{noun, field}
	if v, ok := n.pairs[pair]; ok {
		err = setValue(pv, v)
	} else if e := n.db.QueryRow("select value from mdl_start where noun=? and field=?",
		noun, field).Scan(pv); e != nil {
		err = e
	}
	return err
}

// SetValue sets the named property to the passed value.
func (n *Values) SetField(obj, field string, v interface{}) (err error) {
	return notImplemented
}

func setValue(pv interface{}, value interface{}) (err error) {
	dst := r.ValueOf(pv).Elem()
	src := r.ValueOf(value)
	dst.Set(src)
	return nil
}

var notImplemented = errutil.New("not implemented")
