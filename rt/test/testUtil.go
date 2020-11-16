package test

import (
	"fmt"
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	g "github.com/ionous/iffy/rt/generic"
)

type Kinds struct {
	kinds  kindMap
	fields fieldMap
}

type kindMap map[string]*g.Kind
type fieldMap map[string][]g.Field

// register kinds from a struct using reflection
func (ks *Kinds) AddKinds(is ...interface{}) {
	for _, el := range is {
		ks.fields = kindsForType(ks.fields, r.TypeOf(el).Elem())
	}
}

func (ks *Kinds) Kind(name string) (ret *g.Kind) {
	if k, e := ks.GetKindByName(name); e != nil {
		panic(e)
	} else {
		ret = k
	}
	return
}

//
func (ks *Kinds) GetKindByName(name string) (ret *g.Kind, err error) {
	if k, ok := ks.kinds[name]; ok {
		ret = k // we created the kind already
	} else if fs, ok := ks.fields[name]; !ok {
		err = errutil.New("unknown kind", name)
	} else {
		if ks.kinds == nil {
			ks.kinds = make(kindMap)
		}
		// create the kind from the stored fields
		k := g.NewKind(ks, name, fs)
		ks.kinds[name] = k
		ret = k
	}
	return
}

// generate kinds from a struct using reflection
func kindsForType(kinds fieldMap, t r.Type) fieldMap {
	type stringer interface{ String() string }
	rstringer := r.TypeOf((*stringer)(nil)).Elem()
	if kinds == nil {
		kinds = make(fieldMap)
	}

	var fields []g.Field
	for i, cnt := 0, t.NumField(); i < cnt; i++ {
		f := t.Field(i)
		fieldType := f.Type
		var a affine.Affinity
		var t string
		switch k := fieldType.Kind(); k {
		default:
			panic(errutil.Sprint("unknown kind", k))
		case r.Bool:
			a, t = affine.Text, "aspect"
			// the name of the aspect is the name of the field
			kinds[f.Name] = []g.Field{
				// false first.
				{Name: "Not " + f.Name, Affinity: affine.Bool, Type: "trait"},
				{Name: "Is " + f.Name, Affinity: affine.Bool, Type: "trait"},
			}

		case r.String:
			a, t = affine.Text, k.String()
		case r.Struct:
			a, t = affine.Record, fieldType.Name()
			kinds = kindsForType(kinds, fieldType)

		case r.Slice:
			elType := fieldType.Elem()
			switch k := elType.Kind(); k {
			case r.String:
				a, t = affine.TextList, k.String()
			case r.Float64:
				a, t = affine.NumList, k.String()
			case r.Struct:
				a, t = affine.RecordList, elType.Name()
				kinds = kindsForType(kinds, elType)

			default:
				panic(errutil.Sprint("unknown slice", elType.String()))
			}

		case r.Float64:
			a, t = affine.Number, k.String()

		case r.Int:
			a, t = affine.Text, "aspect"
			if !fieldType.Implements(rstringer) {
				panic("unknown enum")
			}
			x := r.New(fieldType).Elem()
			var traits []g.Field
			for j := int64(0); j < 25; j++ {
				x.SetInt(j)
				trait := x.Interface().(stringer).String()
				end := fmt.Sprintf("%s(%d)", fieldType.Name(), j)
				if trait == end {
					break
				}
				traits = append(traits, g.Field{Name: trait, Affinity: affine.Bool, Type: "trait"})
			}
			aspect := fieldType.Name()
			kinds[aspect] = traits
		}
		fields = append(fields, g.Field{Name: f.Name, Affinity: a, Type: t})

	}
	kinds[t.Name()] = fields
	return kinds
}
