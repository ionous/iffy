package ops

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	r "reflect"
)

func NewTarget(rtype r.Type) Linear {
	rval := r.New(rtype).Elem()
	return inplace(rval)
}

func InPlace(root interface{}) Linear {
	target := r.ValueOf(root).Elem()
	return inplace(target)
}

// Linear flattens a struct and its embedded fields into a nice long list.
// FIX: what i want, actually, is an iterator that walks the fields in order.
// and, when we discover a field has been asked for out of order, then we reset.
type Linear struct {
	r.Value
	fields []FieldIndex
}

type FieldIndex []int

func (c Linear) NumField() int {
	return len(c.fields)
}

// FieldNum returns the value of the requsted field. To maintain compatibility with reflect.Value: on error, Field returns an invalid Value.
func (c Linear) Field(n int) (ret r.Value) {
	if n < len(c.fields) {
		ret = c.FieldByIndex(c.fields[n])
	}
	return
}

func (c Linear) FieldByName(n string) (ret r.Value) {
	// FIX: searching every assign is annoying.
	k := ident.IdOf(n)
	unique.WalkProperties(c.Type(), func(f *r.StructField, idx []int) (done bool) {
		if k == ident.IdOf(f.Name) {
			ret, done = c.FieldByIndex(idx), true
		}
		return
	})
	return
}

func inplace(rval r.Value) Linear {
	var path []FieldIndex
	flatten(rval.Type(), nil, &path)
	return Linear{rval, path}
}

func flatten(rtype r.Type, base []int, plist *[]FieldIndex) {
	for i := 0; i < rtype.NumField(); i++ {
		f := rtype.Field(i)
		IsPublic := len(f.PkgPath) == 0
		if IsPublic {
			IsEmbedded := f.Anonymous && f.Type.Kind() == r.Struct
			idx := append(base, i)
			if !IsEmbedded {
				*plist = append(*plist, idx)
			} else {
				flatten(f.Type, idx, plist)
			}
		}
	}
}
