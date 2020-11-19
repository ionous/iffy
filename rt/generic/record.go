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
			} else if trait, e := v.GetText(); e != nil {
				err = e
			} else {
				// if the field is an aspect, and the caller was asking for a trait...
				// return the state of the trait
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
			if nv, e := DefaultFrom(d.kind.kinds, ft.Affinity, ft.Type); e != nil {
				err = e
			} else {
				ret, err = d.cache(i, nv)
			}
		}
	}
	return
}

// this is a little ugly.
// fix? if DefaultFrom returned a raw value maybe this could be cleaned up.
// -or- if record.values stored refValues ( though that is random extra storage )
func (d *Record) cache(i int, nv Value) (ret Value, err error) {
	if el, ok := nv.(refValue); !ok {
		err = errutil.New("unexpected error creating default values")
	} else {
		d.values[i] = el.v.Interface()
		ret = el
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
		} else if yes, e := val.GetBool(); e != nil {
			err = errutil.New("error setting trait:", e)
		} else if !yes {
			err = errutil.Fmt("error setting trait: couldn't determine the opposite of %q", field)
		} else {
			// set the aspect to the value of the requested trait
			d.values[i] = field
		}
	}
	return
}

func (d *Record) SetFieldByIndex(i int, val Value) (err error) {
	ft := d.kind.fields[i]
	if !affine.MatchTypes(ft.Affinity, ft.Type, val.Affinity(), val.Type()) {
		err = errutil.New("value is not", ft.Affinity, ft.Type)
	} else if el, e := CopyValue(val); e != nil {
		err = e
	} else {
		d.values[i] = el
	}
	return
}
