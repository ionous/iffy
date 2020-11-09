package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/object"
)

type Record struct {
	Nothing
	kind   *Kind
	values []Value
}

var _ Value = (*Record)(nil) // ensure compatibility

func (r *Record) Affinity() affine.Affinity {
	return affine.Record
}

func (r *Record) Type() string {
	return r.kind.name
}

func (r *Record) GetRecord() (ret *Record, err error) {
	ret = r
	return
}

func (r *Record) GetNamedField(field string) (ret Value, err error) {
	switch k := r.kind; field {
	case object.Name:
		err = errutil.New("records don't have names")

	case object.Kind, object.Kinds:
		ret = StringOf(r.kind.name)

	default:
		if i := k.FieldIndex(field); i < 0 {
			err = UnknownField{k.name, field}
		} else {
			fv, ft := r.values[i], k.fields[i]
			if isTrait := ft.Type == "aspect" && ft.Name != field; isTrait {
				ret, err = r.getTraitValue(fv, field)
			} else {
				ret, err = r.getFieldValue(fv, ft)
			}
		}
	}
	return
}

func (r *Record) getTraitValue(fv Value, field string) (ret Value, err error) {
	if fv == nil {
		ret = False
	} else if trait, e := fv.GetText(); e != nil {
		err = e
	} else {
		// if the field is an aspect, and the caller was asking for a trait...
		// return the state of the trait
		ret, err = newBoolValue(trait == field, "")
	}
	return
}

func (r *Record) getFieldValue(fv Value, ft Field) (ret Value, err error) {
	if fv == nil {
		ret, err = DefaultFor(r.kind.kinds, ft.Affinity, ft.Type)
	} else {
		ret = fv
	}
	return
}

func (r *Record) SetNamedField(field string, val Value) (err error) {
	k := r.kind
	if i := k.FieldIndex(field); i < 0 {
		err = UnknownField{k.name, field}
	} else {
		ft := k.fields[i]
		if isTrait := ft.Type == "aspect" && ft.Name != field; isTrait {
			if b, e := val.GetBool(); e != nil {
				err = errutil.New("error setting trait:", e)
			} else if !b {
				err = errutil.Fmt("error setting trait: couldn't determine the opposite of %q", field)
			} else {
				// set the aspect to the value of the requested trait
				r.values[i] = StringOf(field)
			}
		} else if val.Affinity() != ft.Affinity {
			err = errutil.New("value is not", ft.Affinity)
		} else {
			r.values[i] = val
		}
	}
	return
}
