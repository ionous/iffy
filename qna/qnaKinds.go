package qna

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/object"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/tables"
)

type qnaKinds struct {
	kinds                        map[string]*g.Kind
	typeOf, fieldsFor, traitsFor *sql.Stmt // selects field, type for a named kind
}

// aspects are a specific kind of record where every field is a boolean trait
func (km *qnaKinds) addKind(n string, k *g.Kind) *g.Kind {
	if km.kinds == nil {
		km.kinds = make(map[string]*g.Kind)
	}
	km.kinds[n] = k
	return k
}

func (km *qnaKinds) GetKindByName(name string) (ret *g.Kind, err error) {
	name = lang.Breakcase(name)
	if k, ok := km.kinds[name]; ok {
		ret = k
	} else {
		var role string
		if e := km.typeOf.QueryRow(name).Scan(&role); e != nil {
			err = e
		} else {
			if len(role) == 0 {
				err = errutil.Fmt("no such kind %q", name)
			} else if role == object.Aspect {
				if ts, e := km.queryTraits(name); e != nil {
					err = e
				} else {
					ret = km.addKind(name, g.NewKind(km, name, ts))
				}
			} else {
				if ts, e := km.queryFields(name); e != nil {
					err = e
				} else {
					ret = km.addKind(name, g.NewKind(km, role, ts))
				}
			}
		}
	}
	return
}

func (km *qnaKinds) queryFields(name string) (ret []g.Field, err error) {
	// creates the kind if it needs to.
	var field, fieldType string
	var affinity affine.Affinity
	if q, e := km.fieldsFor.Query(name); e != nil {
		err = e
	} else {
		err = tables.ScanAll(q, func() (err error) {
			// by default the type and the affinity are the same
			// ( like reflect, where type and kind are the same for primitive types )
			if len(affinity) == 0 {
				affinity = affine.Affinity(fieldType)
			}
			ret = append(ret, g.Field{
				Name:     field,
				Affinity: affinity,
				Type:     fieldType,
			})
			return
		}, &field, &fieldType, &affinity)
	}
	return
}

func (km *qnaKinds) queryTraits(name string) (ret []g.Field, err error) {
	var trait string
	if q, e := km.traitsFor.Query(name); e != nil {
		err = e
	} else {
		err = tables.ScanAll(q, func() (err error) {
			ret = append(ret, g.Field{
				Name:     trait,
				Affinity: affine.Bool,
				Type:     "trait",
			})
			return
		}, &trait)
	}
	return
}
