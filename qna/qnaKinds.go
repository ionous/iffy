package qna

import (
	"database/sql"

	"github.com/ionous/iffy/affine"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/tables"
)

type qnaKinds struct {
	kinds     map[string]*g.Kind
	fieldsFor *sql.Stmt // selects field, type for a named kind
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
	if k, ok := km.kinds[name]; ok {
		ret = k
	} else {
		// creates the kind if it needs to.
		var fields []g.Field
		var field, fieldType string
		var affinity affine.Affinity
		if q, e := km.fieldsFor.Query(name); e != nil {
			err = e
		} else if e := tables.ScanAll(q, func() (err error) {
			// by default the type and the affinity are the same
			// ( like reflect, where type and kind are the same for primitive types )
			if len(affinity) == 0 {
				affinity = affine.Affinity(fieldType)
			}
			// note: package generic finds the aspect record to lookup traits on demand.
			// fix: should there be a way to pre-populate?
			if err == nil {
				fields = append(fields, g.Field{
					Name:     field,
					Affinity: affinity,
					Type:     fieldType,
				})
			}
			return
		}, &field, &fieldType, &affinity); e != nil {
			err = e
		} else {
			ret = km.addKind(name, g.NewKind(km, name, fields))
		}
	}
	return
}
