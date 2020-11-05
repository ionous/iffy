package qna

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/tables"
)

type qnaKinds struct {
	kinds     map[string]*generic.Kind
	fieldsFor *sql.Stmt // selects field, type for a named kind
}

// aspects are a specific kind of record where every field is a boolean trait
func (km *qnaKinds) addKind(n string, k *generic.Kind) *generic.Kind {
	if km.kinds == nil {
		km.kinds = make(map[string]*generic.Kind)
	}
	km.kinds[n] = k
	return k
}

func (km *qnaKinds) KindByName(name string) (ret *generic.Kind, err error) {
	if k, ok := km.kinds[name]; ok {
		ret = k
	} else {
		// creates the kind if it needs to.
		var aspects []*generic.Kind
		var fields []generic.Field
		var field, fieldType string
		// ex. number, text, aspect
		if q, e := km.fieldsFor.Query(name); e != nil {
			err = e
		} else if e := tables.ScanAll(q, func() (err error) {
			var affinity affine.Affinity
			switch fieldType {
			default:
				// by default the type and the affinity are the same
				// ( just like in go where the type and the kind are the same for primitive types )
				affinity = affine.Affinity(fieldType)
			case "trait":
				affinity = affine.Bool
			case "aspect":
				// aspects are stored as text in the runtime
				affinity = affine.Text
				// we need the aspect record to lookup related traits
				if aspect, e := km.KindByName(name); e != nil {
					err = errutil.New("aspect not found", fieldType, e)
				} else {
					aspects = append(aspects, aspect)
				}
			}
			if err == nil {
				fields = append(fields, generic.Field{
					Name:     field,
					Affinity: affinity,
					Type:     fieldType,
				})
			}
			return
		}, &field, &fieldType); e != nil {
			err = e
		} else {
			ret = km.addKind(name, generic.NewKind(km, name, fields))
		}
	}
	return
}
