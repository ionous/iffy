package generic

import (
	"sort"

	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
)

// we bake it down for faster, easier indexed access.
type Kind struct {
	name      string
	fields    []Field
	traits    []trait
	kinds     Kinds
	lastField int
}

type Field struct {
	Name     string
	Affinity affine.Affinity
	Type     string // ex. kind of record; "aspect", "trait", ...
}

func (k *Kind) Name() string {
	return k.name
}

func (k *Kind) NumField() int {
	return len(k.fields)
}

func (k *Kind) Field(i int) Field {
	k.lastField = i
	return k.fields[i]
}

// searches for the field which handles the passed field
// for traits, it returns the aspect.
// returns -1 if no matching field was found
func (k *Kind) FieldIndex(n string) (ret int) {
	if prev := k.lastField; prev >= 0 && k.fields[prev].Name == n {
		ret = prev
	} else {
		if aspect := findAspect(n, k.traits); len(aspect) > 0 {
			ret = k.fieldIndex(aspect)
		} else {
			ret = k.fieldIndex(n)
		}
		k.lastField = ret
	}
	return
}

func (k *Kind) fieldIndex(field string) (ret int) {
	ret = -1 // provisionally
	for i, f := range k.fields {
		if f.Name == field {
			ret = i
			break
		}
	}
	return
}

func (k *Kind) NewRecord() *Record {
	return &Record{kind: k, values: make([]rt.Value, len(k.fields))}
}

func (k *Kind) NewRecordSlice() *RecordSlice {
	return &RecordSlice{kind: k}
}

// aspects are a specific kind of record where every field is a boolean trait
func NewKind(name string, fields []Field, aspects []*Kind) *Kind {
	var allTraits []trait
	for _, aspect := range aspects {
		ts := makeTraits(aspect)
		allTraits = append(allTraits, ts...)
	}
	sort.Slice(allTraits, func(i, j int) bool {
		it, jt := allTraits[i], allTraits[j]
		return it.Trait < jt.Trait
	})

	return &Kind{name: name, fields: fields}
}
