package qna

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/tables"
)

// Fields implements rt.Fields: key,field,value storage for nouns, kinds, and patterns.
// It reads its data from the play database and caches the results in memory.
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

func NewObjectValues(db *tables.Cache) *Fields {
	return &Fields{make(mapType), db}
}

// SetField to the passed value.

func (n *Fields) SetField(obj, field string, v interface{}) (err error) {
	if strings.HasPrefix(field, object.Prefix) {
		err = errutil.New("can't change internal field", field)
	} else {
		// check if the specified field is a trait
		if a, e := n.GetField(obj+"."+field, object.Aspect); e != nil {
			err = e
		} else {
			// no, just set the field normally.
			if aspect := a.(string); len(aspect) == 0 {
				// fix, future: verify type and existence?
				key := keyType{obj, field}
				n.pairs[key] = v
			} else {
				// yes, then we want to change the aspect not the trait
				if val, ok := v.(bool); !ok || !val {
					err = errutil.Fmt("%q.%q can only be set to true; have %T(%v)", obj, field, v, v)
				} else {
					// set
					err = n.SetField(obj, aspect, field)
				}
			}
		}
	}
	return
}

// GetField sets the value of the passed pointer to the value of the named property.
func (n *Fields) GetField(obj, field string) (ret interface{}, err error) {
	key := keyType{obj, field}
	if val, ok := n.pairs[key]; ok {
		ret = val
	} else {
		switch field {
		case object.Aspect:
			// noun.trait; we use "max" in order to always return a value.
			ret, err = n.getCachingQuery(key,
				`select ifnull(max(aspect),"") from mdl_noun_traits 
				where (noun||'.'||trait)=?`, obj)

		case object.Kind:
			ret, err = n.getCachingQuery(key,
				`select kind from mdl_noun where noun=?`, obj)

		case object.Kinds:
			ret, err = n.getCachingQuery(key,
				`select kind || ( case path when '' then ('') else (',' || path) end ) as path
				from mdl_noun mn 
				join mdl_kind mk 
					using (kind)
				where noun=?`, obj)

		case object.Exists:
			ret, err = n.getCachingQuery(key, `select count() from mdl_noun where noun=?`, obj)

		case object.BoolRule, object.NumberRule, object.TextRule,
			object.ExecuteRule, object.NumListRule, object.TextListRule:
			ret, err = n.getCachingRules(key, obj, field[1:])

		default:
			// see if the user is asking for the status of a trait
			if a, e := n.GetField(obj+"."+field, object.Aspect); e != nil {
				err = e
			} else {
				if aspect := a.(string); len(aspect) > 0 {
					ret, err = n.getCachingStatus(obj, aspect, field)
				} else {
					ret, err = n.getCachingField(key, obj, field)
				}
			}
		}
	}
	return
}

func (n *Fields) GetFieldByIndex(obj string, idx int) (ret string, err error) {
	if idx <= 0 {
		err = errutil.New("GetFieldByIndex out of range", idx)
	} else {
		// first, lookup the parameter name
		key := keyType{obj, "$" + strconv.Itoa(idx)}
		// we use the cache to keep $(idx) -> param name.
		val, ok := n.pairs[key]
		if !ok {
			val, err = n.getCachingQuery(key,
				`select param from mdl_pat where pattern=? and idx=?`,
				obj, idx)
		}
		if field, ok := val.(string); !ok {
			err = fieldNotFound{key.owner, key.member}
		} else {
			ret = field
		}
	}
	return
}

// return true if the object's aspect equals the specified trait.
func (n *Fields) getCachingStatus(obj, aspect, trait string) (ret bool, err error) {
	if val, e := n.GetField(obj, aspect); e != nil {
		err = e
	} else {
		ret = val == trait
	}
	return
}

func (n *Fields) getCachingField(key keyType, obj, field string) (ret interface{}, err error) {
	// FIX? needs more work to determine if the field really exists
	// ex. possibly a union query of class field with a nil value
	if v, e := n.getCachingQuery(key,
		`select value 
		from run_value 
		where noun=? and field=? 
		order by tier asc nulls last limit 1`,
		obj, field); e == nil {
		ret = v
	} else if _, ok := e.(fieldNotFound); !ok {
		err = e
	} else {
		n.pairs[key] = nil
		ret = nil
	}
	return
}

// getCachingQuery uses the rowscanner to write the results of a query into the cache
func (n *Fields) getCachingQuery(key keyType, q string, args ...interface{}) (ret interface{}, err error) {
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
