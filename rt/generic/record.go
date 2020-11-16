package generic

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/object"
)

type Record struct {
	kind   *Kind
	values []interface{}
}

func (r *Record) Type() string {
	return r.kind.name
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
		} else if v, e := r.GetFieldByIndex(i); e != nil {
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
func (r *Record) GetFieldByIndex(i int) (ret Value, err error) {
	if fv, ft := r.values[i], r.kind.fields[i]; fv != nil {
		ret, err = ValueFrom(fv, ft.Affinity, ft.Type)
	} else if nv, e := DefaultFrom(r.kind.kinds, ft.Affinity, ft.Type); e != nil {
		err = e
	} else if rv, ok := nv.(refValue); !ok {
		err = errutil.New("unable to determine default value from %T", nv)
	} else {
		// right now we generate records on demand ( so that we dont have to expand recursive records )
		// fix: assembly should probably throw those types out.
		r.values[i] = rv.v.Interface()
		ret = rv
	}
	return
}

func (r *Record) SetNamedField(field string, val Value) (err error) {
	k := r.kind
	if i := k.FieldIndex(field); i < 0 {
		err = UnknownField{k.name, field}
	} else {
		ft := k.fields[i]
		if isTrait := ft.Type == "aspect" && ft.Name != field; !isTrait {
			err = r.SetFieldByIndex(i, val)
		} else if b, e := val.GetBool(); e != nil {
			err = errutil.New("error setting trait:", e)
		} else if !b {
			err = errutil.Fmt("error setting trait: couldn't determine the opposite of %q", field)
		} else {
			// set the aspect to the value of the requested trait
			r.values[i] = field
		}
	}
	return
}

func (r *Record) SetFieldByIndex(i int, val Value) (err error) {
	ft := r.kind.fields[i]
	if val.Affinity() != ft.Affinity {
		err = errutil.New("value is not", ft.Affinity)
	} else if v, e := CopyValue(val); e != nil {
		err = e
	} else {
		r.values[i] = v
	}
	return
}
