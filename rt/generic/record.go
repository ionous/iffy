package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/object"
)

type Record struct {
	kind   *Kind
	values []interface{}
}

func (d *Record) Kind() *Kind {
	return d.kind
}

func (d *Record) Type() string {
	return d.kind.name
}

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
			ft := k.fields[i]
			if isTrait := ft.Type == "aspect" && ft.Name != field; !isTrait {
				ret = v
			} else {
				// if the field is an aspect, and the caller was asking for a trait...
				// return the state of the trait
				trait := v.String()
				ret, err = newBoolValue(trait == field, "trait")
			}
		}
	}
	return
}

// GetFieldByIndex cant ask for traits, only their aspects.
func (d *Record) GetFieldByIndex(i int) (ret Value, err error) {
	if fv, ft := d.values[i], d.kind.fields[i]; fv != nil {
		ret, err = ValueFrom(fv, ft.Affinity, ft.Type)
	} else {
		if ft.Type == "aspect" {
			if k, e := d.kind.kinds.GetKindByName(ft.Name); e != nil {
				err = e
			} else {
				firstTrait := k.Field(0) // first trait is the default
				if nv, e := ValueFrom(firstTrait.Name, ft.Affinity, ft.Type); e != nil {
					err = e
				} else {
					ret, err = d.cache(i, nv)
				}
			}
		} else {
			if nv, e := NewDefaultValue(d.kind.kinds, ft.Affinity, ft.Type); e != nil {
				err = e
			} else {
				ret, err = d.cache(i, nv)
			}
		}
	}
	return
}

// this is a little ugly.
// fix? if NewDefaultValue returned a raw value maybe this could be cleaned up.
// -or- if record.values stored refValues ( though that is random extra storage )
func (d *Record) cache(i int, nv Value) (ret Value, err error) {
	if el, ok := nv.(refValue); !ok {
		err = errutil.New("unexpected error storing %v(%T)", nv, nv)
	} else {
		ret, d.values[i] = el, el.i
	}
	return
}

func (d *Record) SetNamedField(field string, val Value) (err error) {
	k := d.kind
	if i := k.FieldIndex(field); i < 0 {
		err = UnknownField{k.name, field}
	} else {
		ft := k.fields[i]
		if isTrait := ft.Type == "aspect" && ft.Name != field; !isTrait {
			err = d.SetFieldByIndex(i, val)
		} else if yes := val.Affinity() == affine.Bool && val.Bool(); !yes {
			err = errutil.Fmt("error setting trait: couldn't determine the meaning of %q %s %v", field, val.Affinity(), val)
		} else {
			// set the aspect to the value of the requested trait
			d.values[i] = field
		}
	}
	return
}

func (d *Record) SetFieldByIndex(i int, val Value) (err error) {
	ft := d.kind.fields[i]
	if !matchTypes(ft.Affinity, ft.Type, val.Affinity(), val.Type()) {
		err = errutil.New("value is not", ft.Affinity, ft.Type)
	} else {
		_, err = d.cache(i, val)
	}
	return
}

func matchTypes(a affine.Affinity, at string, b affine.Affinity, bt string) bool {
	return a == b && ((a != affine.Record && a != affine.RecordList) || (at == bt))
}
