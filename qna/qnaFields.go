package qna

import (
	"database/sql"

	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/tables"
)

// Fields implements rt.Fields: key,field,value storage for nouns, kinds, and patterns.
// It reads its data from the play database and caches the results in memory.
type Fields struct {
	valueOf,
	progBytes,
	countOf,
	ancestorsOf,
	kindOf,
	aspectOf,
	nameOf,
	idOf,
	isLike *sql.Stmt
}

func NewFields(db *sql.DB) (ret *Fields, err error) {
	var ps tables.Prep
	f := &Fields{
		valueOf: ps.Prep(db,
			`select value, type
				from run_value 
				where noun=? and field=? 
				order by tier asc nulls last limit 1`),
		progBytes: ps.Prep(db,
			// performs case preferred matching
			`select bytes 
				from mdl_prog
				where UPPER(name) = UPPER(?1)
				and type = ?2
				order by (name != ?1)
				limit 1`),
		countOf: ps.Prep(db,
			`select count(), 'bool' from run_noun where noun=?`),
		ancestorsOf: ps.Prep(db,
			`select kind || ( case path when '' then ('') else (',' || path) end ) as path, 'text'
				from mdl_noun mn 
				join mdl_kind mk 
					using (kind)
				where noun=?`),
		kindOf: ps.Prep(db,
			`select kind, 'text' from mdl_noun where noun=?`),
		// return the name of the aspect of the specified trait, or the empty string.
		aspectOf: ps.Prep(db,
			`select aspect, 'text' from mdl_noun_traits 
				where (noun||'.'||trait)=?`),
		// given an id, find the name
		nameOf: ps.Prep(db,
			`select name, 'text' 
				from mdl_name
				join run_noun
					using (noun)
				where noun=?
				order by rank
				limit 1`),
		// given a name, find the id
		idOf: ps.Prep(db,
			`select noun, 'text'
				from mdl_name
				join run_noun
					using (noun)
				where UPPER(name)=UPPER(?)
				order by rank
				limit 1`),
		// use the sqlite like function to match
		isLike: ps.Prep(db,
			`select ? like ?`),
	}
	if e := ps.Err(); e != nil {
		err = e
	} else {
		ret = f
	}
	return
}

func (n *Runner) IsLike(a, b string) (ret bool, err error) {
	err = n.fields.isLike.QueryRow(a, b).Scan(&ret)
	return
}

func (n *Runner) SetField(target, field string, val rt.Value) (err error) {
	if len(target) == 0 || target == object.Name || (target[0] == object.Prefix && target != object.Variables) {
		err = errutil.Fmt("can't change reserved field '%s.%s'", target, field)
	} else {
		switch e := n.ScopeStack.SetField(target, field, val); e.(type) {
		default:
			err = e
		case rt.UnknownTarget, rt.UnknownField:
			key := makeKey(target, field)
			err = n.setField(key, val)
		}
	}
	return
}

func (n *Runner) setField(key keyType, val rt.Value) (err error) {
	// first, check if the specified field refers to a trait
	switch p, e := n.getField(key.dot(object.Aspect)); e.(type) {
	default:
		err = e // there was an unknown error
	case nil:
		// get the name of the aspect
		if aspect, e := p.GetText(n); e != nil {
			err = e
		} else {
			// we want to change the aspect not the trait...
			if val, e := val.GetBool(n); e != nil {
				err = errutil.Fmt("error setting trait; have %v %v %s", key, val, e)
			} else if !val {
				// future: might maintain a table of opposite names ( similar to plurals )
				err = errutil.Fmt("error setting trait; %q can only be set to true, have %v", key, val)
			} else {
				// recurse...
				targetAspect := keyType{key.target, aspect}
				err = n.setField(targetAspect, &generic.String{Value: key.field})
			}
		}
	case rt.UnknownField:
		// didnt refer to a trait, so just set the field normally.
		if p, e := n.getField(key); e != nil {
			err = e
		} else {
			// note: we dont replace the generic value in the cache
			// we poke into that "box" and set its internal value.
			// in the process it validates the incoming data type.
			err = p.SetValue(n, val)
		}
	}
	return
}

// pv is a pointer to a pattern instance, and we copy its contents in.
func (n *Runner) GetEvalByName(name string, pv interface{}) (err error) {
	outVal := r.ValueOf(pv).Elem() // outVal is a pattern instance who's fields get overwritten
	rtype := outVal.Type()
	// note: makeKey camelCases, while go types are PascalCase
	// this automatically keeps them from conflicting.
	key := makeKeyForEval(name, rtype.Name())
	if val, ok := n.pairs[key]; ok {
		store := r.ValueOf(val.Interface())
		outVal.Set(store)
	} else {
		var store interface{}
		switch e := n.fields.progBytes.QueryRow(key.target, key.field).Scan(&tables.GobScanner{outVal}); e {
		case nil:
			store = outVal.Interface()
		case sql.ErrNoRows:
			err = key.unknown()
		default:
			err = e
		}
		n.pairs[key] = generic.NewValue("", store)
	}
	return
}

func (n *Runner) GetField(target, field string) (ret rt.Value, err error) {
	switch v, e := n.ScopeStack.GetField(target, field); e.(type) {
	default:
		err = e
	case nil:
		ret = v
	case rt.UnknownTarget, rt.UnknownField:
		ret, err = n.getField(makeKey(target, field))
	}
	return
}

// check the cache before asking the database for info
func (n *Runner) getField(key keyType) (ret *generic.Value, err error) {
	if val, ok := n.pairs[key]; !ok {
		ret, err = n.cacheField(key)
	} else if val == nil {
		err = key.unknown()
	} else {
		ret = val
	}
	return
}

// when we know that the field is not a reserved field, and we just want to check the value.
// ie. for aspects
func (n *Runner) getValue(key keyType) (ret *generic.Value, err error) {
	if val, ok := n.pairs[key]; !ok {
		ret, err = n.cacheQuery(key, n.fields.valueOf, key.target, key.field)
	} else if val == nil {
		err = key.unknown()
	} else {
		ret = val
	}
	return
}

var depth = 0

func (n *Runner) cacheField(key keyType) (ret *generic.Value, err error) {
	switch target, field := key.target, key.field; target {
	case object.Name:
		// search for the object name using the object's id
		ret, err = n.cacheQuery(key, n.fields.nameOf, field)

	case object.Id:
		// search for the object id by a partial object name
		ret, err = n.cacheQuery(key, n.fields.idOf, field)

	case object.Aspect:
		// return the name of an aspect for a trait
		ret, err = n.cacheQuery(key, n.fields.aspectOf, field)

	case object.Kind:
		ret, err = n.cacheQuery(key, n.fields.kindOf, field)

	case object.Kinds:
		ret, err = n.cacheQuery(key, n.fields.ancestorsOf, field)

	case object.Exists:
		// searches for an id match; never returns UnknownField
		ret, err = n.cacheQuery(key, n.fields.countOf, field)

	default:
		// see if the user is asking for the status of a trait
		switch aspectOfTrait, e := n.getField(key.dot(object.Aspect)); e.(type) {
		default:
			err = e
		case rt.UnknownField:
			ret, err = n.cacheQuery(key, n.fields.valueOf, target, field)
		case nil:
			// we found the aspect name from the trait
			// now we need to ask for the current value of the aspect
			if aspectName, e := aspectOfTrait.GetText(n); e != nil {
				err = e
			} else {
				aspectOfTarget := keyType{target, aspectName}
				if aspectValue, e := n.getValue(aspectOfTarget); e != nil {
					err = e
				} else if trait, e := aspectValue.GetText(n); e != nil {
					err = errutil.Fmt("unexpected value in aspect '%v.%v' %v", target, aspectName, e)
				} else {
					// return true if the object's aspect equals the specified trait.
					ret = generic.NewValue("bool", trait == field)
				}
			}
		}
	}
	return
}

// query the db and store the returned value in the cache.
// note: all of the queries are expected to return two parts: the value and the typeName.
func (n *Runner) cacheQuery(key keyType, q *sql.Stmt, args ...interface{}) (ret *generic.Value, err error) {
	var v interface{}
	var t string
	switch e := q.QueryRow(args...).Scan(&v, &t); e {
	case nil:
		q := generic.NewValue(t, v)
		n.pairs[key] = q
		ret = q
	case sql.ErrNoRows:
		n.pairs[key] = nil
		err = key.unknown()
	default:
		err = errutil.New("runtime error:", e)
	}
	return
}
