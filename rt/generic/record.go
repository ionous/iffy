package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/object"
)

type Record struct {
	kind   *Kind
	values []Value
}

func (d *Record) Kind() *Kind {
	return d.kind
}

func (d *Record) Type() string {
	return d.kind.name
}

// GetNamedField distinguishes itself from Value.FieldByName to help with find in files.
// Record doesnt directly implement generic.Value, nor any method other than "Type"
func (d *Record) GetNamedField(field string) (ret Value, err error) {
	switch k := d.kind; field {
	case object.Name:
		err = errutil.New("records don't have names")

	case object.Kind, object.Kinds:
		ret = StringOf(d.kind.name)

	default:
		if i := k.FieldIndex(field); i < 0 {
			err = UnknownField{k.name, field}
		} else if v, e := d.GetFieldByIndex(i); e != nil {
			err = e
		} else {
			ft := k.fields[i] // isTrait if we found aspect (a) while looking for field (t)
			if isTrait := ft.Type == "aspect" && ft.Name != field; !isTrait {
				ret = v
			} else {
				// we were looking for trait (t)
				trait := v.String()
				ret = BoolFrom(trait == field, "trait")
			}
		}
	}
	return
}

// GetFieldByIndex cant ask for traits, only their aspects.
func (d *Record) GetFieldByIndex(i int) (ret Value, err error) {
	if fv, ft := d.values[i], d.kind.fields[i]; fv != nil {
		ret = fv
	} else {
		if ft.Type == "aspect" {
			// if we're asking for an aspect, the default value will be the string of its first trait
			if k, e := d.kind.kinds.GetKindByName(ft.Name); e != nil {
				err = e
			} else {
				firstTrait := k.Field(0)                   // first trait is the default
				nv := StringFrom(firstTrait.Name, "trait") // better as "aspect", "trait", or something else?
				ret, d.values[i] = nv, nv
			}
		} else {
			if nv, e := NewDefaultValue(d.kind.kinds, ft.Affinity, ft.Type); e != nil {
				err = e
			} else {
				ret, d.values[i] = nv, nv
			}
		}
	}
	return
}

func (d *Record) SetNamedField(field string, val Value) (err error) {
	k := d.kind
	if i := k.FieldIndex(field); i < 0 {
		err = UnknownField{k.name, field}
	} else {
		ft := k.fields[i] // isTrait if we found aspect (a) while looking for field (t)
		if isTrait := ft.Type == "aspect" && ft.Name != field; !isTrait {
			err = d.SetFieldByIndex(i, val)
		} else if yes := val.Affinity() == affine.Bool && val.Bool(); !yes {
			err = errutil.Fmt("error setting trait: couldn't determine the meaning of %q %s %v", field, val.Affinity(), val)
		} else {
			// set the aspect to the value of the requested trait
			d.values[i] = StringFrom(field, "aspect")
		}
	}
	return
}

func (d *Record) SetFieldByIndex(i int, val Value) (err error) {
	ft := d.kind.fields[i]
	if a, t := val.Affinity(), val.Type(); !matchTypes(d.kind.kinds, ft.Affinity, ft.Type, a, t) {
		err = errutil.Fmt("%s of %s is not %s of %s ( setting field %q )", a, t, ft.Affinity, ft.Type, ft.Name)
	} else {
		d.values[i] = val
	}
	return
}

func matchTypes(ks Kinds, fa affine.Affinity, ft string, va affine.Affinity, vt string) (okay bool) {
	if fa == va {
		recordLike := fa == affine.Object || fa == affine.Record || fa == affine.RecordList
		if !recordLike {
			okay = true
		} else if vk, e := ks.GetKindByName(vt); e == nil {
			// a field takes: cats
			// my value is things, animals, cats, tigers.
			// so we should ask: does the value's path contains the field's kind.
			okay = vk.Implements(ft)
		}
	}
	return
}
