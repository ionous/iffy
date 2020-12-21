package qna

import (
	"database/sql"
	"strings"

	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/object"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/tables"
)

// Fields implements rt.Fields: key,field,value storage for nouns, kinds, and patterns.
// It reads its data from the play database and caches the results in memory.
type Fields struct {
	activeDomains,
	activeNouns,
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
	isLike,
	relativesOf,
	reciprocalOf,
	relateTo,
	relativeKinds,
	updatePairs *sql.Stmt
}

func NewFields(db *sql.DB) (ret *Fields, err error) {
	var ps tables.Prep
	f := &Fields{
		activeDomains: ps.Prep(db,
			`select 1 from run_domain where active and domain=?`),
		activeNouns: ps.Prep(db,
			// instr(X,Y) finds the first occurrence of string Y in string X
			`select 1 from 
			mdl_noun mn 
			join run_domain rd 
			where rd.active and instr(mn.noun, '#' || rd.domain || '::') = 1
			and mn.noun=?`),

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
			`select kind || ( case path when '' then ('') else (',' || path) end ) as path
				from mdl_noun mn 
				join mdl_kind mk 
					using (kind)
				where noun=?`),

		kindOf: ps.Prep(db,
			`select kind
				from mdl_noun 
				where noun=?`),
		fieldsFor: ps.Prep(db,
			`select field, type from mdl_field
			where kind=?1
				union all
			select * from (
				select trait, 'trait' 
				from mdl_aspect
				where aspect = ?1
				order by rank 
			)`),
		traitsFor: ps.Prep(db,
			`select trait
				from mdl_aspect 
				where aspect=?
				order by rank`),
		// return the name of the aspect of the specified trait, or the empty string.
		aspectOf: ps.Prep(db,
			`select aspect
				from mdl_noun_traits 
				where (noun||'.'||trait)=?`),
		// given an id, find the name
		nameOf: ps.Prep(db,
			`select name
				from mdl_name
				join mdl_noun
					using (noun)
				where noun=?
				order by rank
				limit 1`),
		// given a name, find the id
		objOf: ps.Prep(db,
			`select noun
				from mdl_name
				join mdl_noun
					using (noun)
				where UPPER(name)=UPPER(?)
				order by rank
				limit 1`),
		// use the sqlite like function to match
		isLike: ps.Prep(db,
			`select ? like ?`),
		relativeKinds: ps.Prep(db,
			`select mr.kind, mr.otherKind, mr.cardinality
				from mdl_rel mr 
				where relation=?`),
		// instead of separately deleting old pairs and inserting new ones;
		// we insert and replace active ones.
		updatePairs: ps.Prep(db,
			`with next as (
			select noun, otherNoun, relation, cardinality 
			from mdl_pair 
			join mdl_rel mr 
				using (relation)
			where ?=ifnull(domain, 'entire_game')
			)
			insert or replace into run_pair
			select prev.noun, relation, prev.otherNoun, 0
				from next
				join run_pair prev 
					using (relation)
				where  ((prev.noun = next.noun and next.cardinality glob '*_one') or
						(prev.otherNoun = next.otherNoun and next.cardinality glob 'one_*')) 
			union all
			select next.noun, relation, next.otherNoun, 1 
			from next`),
		relativesOf: ps.Prep(db,
			`select otherNoun from run_pair where active and noun=?1 and relation=?2`),
		reciprocalOf: ps.Prep(db,
			`select noun from run_pair where active and otherNoun=?1 and relation=?2`),
		relateTo: ps.Prep(db,
			`with next as (
				select ?1 as noun, ?2 as otherNoun, ?3 as relation, ?4 as cardinality
			)
			insert or replace into run_pair
			select prev.noun, relation, prev.otherNoun, 0
			from next
			join run_pair prev 
				using (relation)
			where  ((prev.noun = next.noun and next.cardinality glob '*_one') or
					(prev.otherNoun = next.otherNoun and next.cardinality glob 'one_*')) 
			union all 
			select next.noun, relation, next.otherNoun, 1
			from next`),
	}
	if e := ps.Err(); e != nil {
		err = e
	} else {
		ret = f
	}
	return
}

func (n *Fields) UpdatePairs(domain string) (ret int, err error) {
	if res, e := n.updatePairs.Exec(domain); e != nil {
		err = e
	} else {
		ret = tables.RowsAffected(res)
	}
	return
}

func (n *Runner) IsLike(a, b string) (ret bool, err error) {
	err = n.fields.isLike.QueryRow(a, b).Scan(&ret)
	return
}

func (n *Runner) SetField(target, field string, val g.Value) (err error) {
	if len(target) == 0 {
		err = errutil.Fmt("no target specified for field %q", field)
	} else if writable := target[0] != object.Prefix ||
		target == object.Variables ||
		target == object.Counter; !writable {
		err = errutil.Fmt("can't change reserved field '%s.%s'", target, field)
	} else {
		switch e := n.ScopeStack.SetField(target, field, val); e.(type) {
		default:
			err = e
		case g.UnknownTarget, g.UnknownField:
			key := makeKey(target, field)
			err = n.setField(key, val)
		}
	}
	return
}

func (n *Runner) setField(key keyType, val g.Value) (err error) {
	// first, check if the specified field refers to a trait
	switch aspectOfTrait, e := n.GetField(object.Aspect, key.dot()); e.(type) {
	default:
		err = e // there was an unknown error
	case nil:
		if aspectName, b := aspectOfTrait.String(), val.Bool(); !b {
			// future: might maintain a table of opposite names ( similar to plurals )
			err = errutil.Fmt("error setting trait: couldn't determine the opposite of %q", key)
		} else {
			// recurse...
			targetAspect := keyType{key.target, aspectName}
			err = n.setField(targetAspect, g.StringOf(key.field))
		}

	case g.UnknownField:
		// didnt refer to a trait, so just set the field normally.
		// ( to set the field, we get the field to verify it exists, and to check affinity )
		if q, e := n.getOrCache(key.target, key.field, n.queryFieldValue); e != nil {
			err = e
		} else if a := q.Affinity(); a != val.Affinity() {
			err = errutil.New("value is not", a)
		} else if v, e := g.CopyValue(val); e != nil {
			err = e
		} else {
			n.pairs[key] = staticValue{a, v}
		}
	}
	return
}

// pv is a pointer to a pattern instance, and we copy its contents in.
func (n *Runner) GetEvalByName(name string, pv interface{}) (err error) {
	name = lang.Breakcase(name)
	outVal := r.ValueOf(pv).Elem() // outVal is a pattern instance who's fields get overwritten
	rtype := outVal.Type()
	// note: makeKey camelCases, while go types are PascalCase
	// this automatically keeps them from conflicting.
	key := makeKeyForEval(name, rtype.Name())
	if q, ok := n.pairs[key]; ok {
		eval := q.(patternValue).store
		rval := r.ValueOf(eval)
		outVal.Set(rval)
	} else {
		var val qnaValue
		switch e := n.fields.progBytes.QueryRow(key.target, key.field).Scan(&tables.GobScanner{outVal}); e {
		case nil:
			store := outVal.Interface()
			val = patternValue{store}
			// pretty.Println(store)
		case sql.ErrNoRows:
			err = key.unknown()
			val = errorValue{err}
		default:
			err = e
			val = errorValue{err}
		}
		// see notes: in theory GetEvalByName with
		n.pairs[key] = val
	}
	return
}

func (n *Runner) GetField(target, field string) (ret g.Value, err error) {
	switch target {
	case object.Aspect:
		// used internally: return the name of an aspect for a trait
		ret, err = n.getOrCache(object.Aspect, field, func(key keyType) (ret qnaValue, err error) {
			var val string
			if e := n.fields.aspectOf.QueryRow(field).Scan(&val); e != nil {
				err = e
			} else {
				ret = staticValue{affine.Text, val}
			}
			return
		})

	case object.Domain:
		// fix,once there's a domain hierarchy:
		// store the active path and test using find in path.
		var b bool
		if e := n.fields.activeDomains.QueryRow(field).Scan(&b); e != nil {
			err = e
		} else {
			ret = g.BoolOf(b)
		}

	case object.Kind:
		ret, err = n.getOrCache(object.Kind, field, func(key keyType) (ret qnaValue, err error) {
			var val string
			if e := n.fields.kindOf.QueryRow(field).Scan(&val); e != nil {
				err = e
			} else {
				ret = staticValue{affine.Text, val}
			}
			return
		})

	case object.Kinds:
		ret, err = n.getOrCache(object.Kinds, field, func(key keyType) (ret qnaValue, err error) {
			var val string
			if e := n.fields.ancestorsOf.QueryRow(field).Scan(&val); e != nil {
				err = e
			} else {
				ret = staticValue{affine.Text, val}
			}
			return
		})

	case object.Locale:
		// find the name of the parent, then return that cached object
		if parent, e := n.nounLocale.localeOf(field); e != nil {
			err = e
		} else if len(parent) == 0 {
			err = g.UnknownObject("") // fix: what's the right value for empty value?
		} else {
			ret, err = n.GetField(object.Value, parent)
		}

	case object.Name:
		// given an id, make sure the object should be available,
		// then return its author given name.
		if !n.activeNouns.isActive(field) {
			err = g.UnknownObject(field)
		} else {
			ret, err = n.getOrCache(object.Name, field, func(key keyType) (ret qnaValue, err error) {
				var val string
				if e := n.fields.nameOf.QueryRow(field).Scan(&val); e != nil {
					err = e
				} else {
					ret = staticValue{affine.Text, val}
				}
				return
			})
		}

	case object.Value:
		// fix: internal object handling needs some love; i dont much like the # test.
		if strings.HasPrefix(field, "#") {
			if !n.activeNouns.isActive(field) {
				// fix: differentiate b/t unknown and unavailable?
				err = g.UnknownObject(field)
			} else {
				ret, err = n.getOrCache(object.Value, field, func(key keyType) (ret qnaValue, err error) {
					ret = &qnaObject{n: n, id: field}
					return
				})
			}
		} else {
			// given a name, find an object (id) and make sure it should be available
			ret, err = n.getOrCache(object.Value, field, func(key keyType) (ret qnaValue, err error) {
				var id string
				if e := n.fields.objOf.QueryRow(field).Scan(&id); e != nil {
					err = e
				} else {
					if !n.activeNouns.isActive(id) {
						err = g.UnknownObject(id)
					} else {
						ret = &qnaObject{n: n, id: id}
					}
				}
				return
			})
		}

	default:
		switch v, e := n.ScopeStack.GetField(target, field); e.(type) {
		default:
			err = e
		case nil:
			ret = v
		case g.UnknownTarget, g.UnknownField:
			key := makeKey(target, field)
			if q, ok := n.pairs[key]; ok {
				ret, err = q.Snapshot(n)
			} else {
				// first: loop. ask if we are trying to find the value of a trait. ( noun.trait )
				switch aspectOfTrait, e := n.GetField(object.Aspect, key.dot()); e.(type) {
				default:
					err = e
				case nil:
					// we found the aspect name from the trait
					// now we need to ask for the current value of the aspect
					aspectName := aspectOfTrait.String()
					if q, e := n.getOrCache(key.target, aspectName, n.queryFieldValue); e != nil {
						err = e
					} else {
						// return whether the object's aspect equals the specified trait.
						// ( we dont cache this value because multiple things can change it )
						ret = g.BoolOf(key.field == q.String())
					}
				case g.UnknownField:
					// it wasnt a trait, so query the field value
					// fix: b/c its more common, should we do this first?
					ret, err = n.getOrCache(key.target, key.field, n.queryFieldValue)
				}
				return
			}
		}
	}
	return
}

// check the cache before asking the database for info
func (n *Runner) getOrCache(target, field string, cache func(key keyType) (ret qnaValue, err error)) (ret g.Value, err error) {
	key := makeKey(target, field)
	if q, ok := n.pairs[key]; ok {
		ret, err = q.Snapshot(n)
	} else {
		switch val, e := cache(key); e {
		case nil:
			ret, err = n.store(key, val)

		case sql.ErrNoRows:
			ret, err = n.store(key, errorValue{key.unknown()})

		default:
			err = errutil.New("runtime error:", e)
		}
	}
	return
}

// query the db for the value of an noun's field
func (n *Runner) queryFieldValue(key keyType) (ret qnaValue, err error) {
	var i interface{}
	var a affine.Affinity
	if e := n.fields.valueOf.QueryRow(key.target, key.field).Scan(&i, &a); e != nil {
		err = e
	} else {
		switch v := i.(type) {
		default:
			ret = staticValue{a, v}
		case []byte:
			if p, e := newEval(a, v); e != nil {
				err = e
			} else {
				ret = p
			}
		}
	}
	return
}

// store the passed value generator, and return the latest snapshot of it
func (n *Runner) store(key keyType, val qnaValue) (ret g.Value, err error) {
	n.pairs[key] = val
	return val.Snapshot(n)
}
