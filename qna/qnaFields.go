package qna

import (
	"database/sql"

	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
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
	fieldsFor,
	traitsFor,
	aspectOf,
	nameOf,
	objOf,
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
		// countOf: ps.Prep(db,
		// 	`select count(), 'bool' from run_noun where noun=?`),
		ancestorsOf: ps.Prep(db,
			`select kind || ( case path when '' then ('') else (',' || path) end ) as path, 'text'
				from mdl_noun mn 
				join mdl_kind mk 
					using (kind)
				where noun=?`),
		kindOf: ps.Prep(db,
			`select kind, 'text' 
				from mdl_noun 
				where noun=?`),
		fieldsFor: ps.Prep(db,
			`select * from mdl_field
				union all
			select * from (
				select aspect as kind, trait as field, 'trait' from mdl_aspect
				order by rank 
			)
			where kind=?`),
		traitsFor: ps.Prep(db,
			`select trait
				from mdl_aspect 
				where aspect=?
				order by rank`),
		// return the name of the aspect of the specified trait, or the empty string.
		aspectOf: ps.Prep(db,
			`select aspect, 'text' 
				from mdl_noun_traits 
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
		objOf: ps.Prep(db,
			`select noun, 'object'
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
	if len(target) == 0 {
		err = errutil.Fmt("no target specified for field %q", field)
	} else if writable := target[0] != object.Prefix ||
		target == object.Variables ||
		target == object.Counter; !writable {
		err = errutil.Fmt("can't change reserved field '%s.%s'", target, field)
	} else {
		// fix? implement a proper move
		if val == nil {
			if x, e := n.GetField(target, field); e != nil {
				err = e
			} else {
				val, err = generic.MakeDefault(&n.kinds, x.Affinity(), x.Type())
			}
		}
		if err == nil {
			switch e := n.ScopeStack.SetField(target, field, val); e.(type) {
			default:
				err = e
			case rt.UnknownTarget, rt.UnknownField:
				key := makeKey(target, field)
				err = n.setField(key, val)
			}
		}
	}
	return
}

func (n *Runner) setField(key keyType, val rt.Value) (err error) {
	// first, check if the specified field refers to a trait
	switch q, e := n.getField(key.dot(object.Aspect)); e.(type) {
	default:
		err = e // there was an unknown error
	case nil:
		// get the name of the aspect
		if aspect, e := q.GetText(); e != nil {
			err = e
		} else {
			// we want to change the aspect not the trait...
			if b, e := val.GetBool(); e != nil {
				err = errutil.New("error setting trait:", e)
			} else if !b {
				// future: might maintain a table of opposite names ( similar to plurals )
				err = errutil.Fmt("error setting trait: couldn't determine the opposite of %q", key)
			} else {
				// recurse...
				targetAspect := keyType{key.target, aspect}
				err = n.setField(targetAspect, generic.NewString(key.field))
			}
		}
	case rt.UnknownField:
		// didnt refer to a trait, so just set the field normally.
		if q, e := n.cacheField(key); e != nil {
			err = e
		} else if a := q.Affinity(); a != val.Affinity() {
			err = errutil.New("value is not", a)
		} else {
			n.pairs[key] = val
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
	if q, ok := n.pairs[key]; ok {
		eval := q.(*evalValue).eval
		store := r.ValueOf(eval)
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
		// see notes: in theory GetEvalByName with
		n.pairs[key] = &evalValue{run: n, eval: store}
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
		if q, e := n.getField(makeKey(target, field)); e != nil {
			err = e
		} else {
			ret = q
		}
	}
	return
}

// check the cache before asking the database for info
func (n *Runner) getField(key keyType) (ret rt.Value, err error) {
	if q, ok := n.pairs[key]; !ok {
		ret, err = n.cacheField(key)
	} else if q == nil {
		err = key.unknown()
	} else {
		ret = q
	}
	return
}

// when we know that the field is not a reserved field, and we just want to check the value.
// ie. for aspects
func (n *Runner) getValue(key keyType) (ret rt.Value, err error) {
	if q, ok := n.pairs[key]; !ok {
		ret, err = n.cacheQuery(key, n.fields.valueOf, key.target, key.field)
	} else if q == nil {
		err = key.unknown()
	} else {
		ret = q
	}
	return
}

func (n *Runner) cacheField(key keyType) (ret rt.Value, err error) {
	switch target, field := key.target, key.field; target {
	case object.Name:
		// search for the object name using the object's id
		ret, err = n.cacheQuery(key, n.fields.nameOf, field)

	case object.Value:
		switch v, e := n.cacheQuery(key, n.fields.objOf, field); e.(type) {
		default:
			ret, err = v, e
		case rt.UnknownField:
			err = rt.UnknownObject(field)
		}
		//

	case object.Aspect:
		// used internally: return the name of an aspect for a trait
		ret, err = n.cacheQuery(key, n.fields.aspectOf, field)

	case object.Kind:
		ret, err = n.cacheQuery(key, n.fields.kindOf, field)

	case object.Kinds:
		ret, err = n.cacheQuery(key, n.fields.ancestorsOf, field)

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
			if aspectName, e := aspectOfTrait.GetText(); e != nil {
				err = e
			} else {
				aspectOfTarget := keyType{target, aspectName}
				if q, e := n.getValue(aspectOfTarget); e != nil {
					err = e
				} else if trait, e := q.GetText(); e != nil {
					err = errutil.Fmt("unexpected value in aspect '%v.%v' %v", target, aspectName, e)
				} else {
					// return whether the object's aspect equals the specified trait.
					// ( we dont cache this value because multiple things can change it )
					ret = generic.NewBool(trait == field)
				}
			}
		}
	}
	return
}

// query the db and store the returned value in the cache.
// note: all of the queries are expected to return two parts: the value and the typeName.
func (n *Runner) cacheQuery(key keyType, q *sql.Stmt, args ...interface{}) (ret rt.Value, err error) {
	var raw interface{}
	var a affine.Affinity
	switch e := q.QueryRow(args...).Scan(&raw, &a); e {
	case nil:
		if v, e := newValue(n, a, raw); e != nil {
			err = e
		} else {
			n.pairs[key] = v
			ret = v
		}
	case sql.ErrNoRows:
		n.pairs[key] = nil
		err = key.unknown()
	default:
		err = errutil.New("runtime error:", e)
	}
	return
}
