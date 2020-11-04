package qna

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/tables"
)

type qnaKinds map[string]*generic.Kind

//  interface for generic notary implementation
func (n *Runner) KindByName(name string) (ret *generic.Kind, err error) {
	if k, ok := n.kinds[name]; ok {
		ret = k
	} else {
		var aspects []*generic.Kind
		var fields []generic.Field
		var field, fieldType string
		// ex. number, text, aspect
		if q, e := n.fields.fieldsOf.Query(name); e != nil {
			err = e
		} else if e := tables.ScanAll(q, func() (err error) {
			var affinity affine.Affinity
			switch fieldType {
			default:
				// by default the type and the affinity are the same
				// ( just like in go where the type and the kind are the same for primitive types )
				affinity = affine.Affinity(fieldType)
			case "aspect":
				// aspects are stored as text in the runtime
				affinity = "text"
				// we need the aspect record to lookup related traits
				if aspect, e := n.KindByName(name); e != nil {
					err = errutil.Append(err, errutil.New("aspect not found", fieldType, e))
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
			k := generic.NewKind(name, fields, aspects)
			n.kinds[name] = k
			ret = k
		}
	}
	return
}
