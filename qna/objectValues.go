package qna

import (
	"database/sql"
	"strconv"

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

// SetField to the passed value.
// fix, future: verify type?
func (n *Fields) SetField(obj, field string, v interface{}) (err error) {
	key := keyType{obj, field}
	n.pairs[key] = v
	return
}

// GetField sets the value of the passed pointer to the value of the named property.
func (n *Fields) GetField(obj, field string) (ret interface{}, err error) {
	key := keyType{obj, field}
	if val, ok := n.pairs[key]; ok {
		ret = val
	} else {
		switch field {
		case object.Kind:
			ret, err = n.cacheField(key, `select kind from mdl_noun where noun=?`, obj)

		case object.Kinds:
			ret, err = n.cacheField(key,
				`select kind || ( case path when '' then ('') else ("," || path) end ) as path
				from mdl_noun mn 
				join mdl_kind mk 
					using (kind)
				where noun=?`, obj)

		case object.Exists:
			ret, err = n.cacheField(key, `select count() from mdl_noun where noun=?`, obj)

		case object.BoolRule, object.NumberRule, object.TextRule,
			object.ExecuteRule, object.NumListRule, object.TextListRule:
			ret, err = n.cacheRules(key, obj, field[1:])

		default:
			// FIX? needs more work to determine if the field really exists
			// ex. possibly a union query of class field with a nil value
			if v, e := n.cacheField(key, `select value 
				from run_init 
				where noun=? and field=? 
				order by tier limit 1`,
				obj, field); e == nil {
				ret = v
			} else if _, ok := e.(fieldNotFound); !ok {
				err = e
			} else {
				n.pairs[key] = nil
				ret = nil
			}
		}
	}
	return
}

func (n *Fields) GetFieldByIndex(obj string, idx int) (ret interface{}, err error) {
	if idx <= 0 {
		err = errutil.New("GetFieldByIndex out of range", idx)
	} else {
		// first, lookup the parameter name
		key := keyType{obj, "$" + strconv.Itoa(idx)}
		// we use the cache to keep $(idx) -> param name.
		val, ok := n.pairs[key]
		if !ok {
			val, err = n.cacheField(key,
				`select param from mdl_pat where pattern=? and idx=?`,
				obj, idx)
		}
		if err == nil {
			if field, ok := val.(string); !ok {
				err = fieldNotFound{key.owner, key.member}
			} else {
				ret, err = n.GetField(obj, field)
			}
		}
	}
	return
}

func (n *Fields) cacheField(key keyType, q string, args ...interface{}) (ret interface{}, err error) {
	tgt := mapTarget{key: key, pairs: n.pairs}
	switch e := n.db.QueryRow(q, args...).Scan(&tgt); e {
	case nil:
		ret = tgt.value
	case sql.ErrNoRows:
		err = fieldNotFound{key.owner, key.member}
	default:
		err = e
	}
	return
}

var notImplemented = errutil.New("not implemented")
