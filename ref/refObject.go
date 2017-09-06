package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/enum"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

type RefObject struct {
	id      ident.Id
	value   r.Value // stores the concrete value. ex. Rock, not *Rock.
	objects ObjectMap
}

// Id returns the unique identifier for this Object.
// Blank for anonymous and temporary objects.
func (n RefObject) Id() ident.Id {
	return n.id
}

// Value holding the data for the object.
func (n RefObject) Value() r.Value {
	return n.value
}

// String representation of the object.
func (n RefObject) String() (ret string) {
	if n.id.IsValid() {
		ret = n.id.Name
	} else {
		ret = n.value.Type().Name()
	}
	return
}

// Type implements rt.Object.
func (n RefObject) Type() r.Type {
	return n.value.Type()
}

// Property implements rt.Object.
func (n RefObject) Property(name string) (ret rt.Property, okay bool) {
	pid := ident.IdOf(name)
	if path, idx := enum.PropertyPath(n.value.Type(), pid.Name); len(path) > 0 {
		field := n.value.FieldByIndex(path)
		ret, okay = RefProp{pid, n, field, idx}, true
	}
	return
}

type RefProp struct {
	id      ident.Id
	obj     RefObject
	r.Value     // underlying value within the object
	index   int // index of choice if property was accessed by state
}

func (p RefProp) String() string {
	return p.obj.String() + "." + p.id.Name
}

func (p RefProp) Id() ident.Id {
	return p.id
}

// GetValue implements rt.Property.
func (p RefProp) GetValue(pv interface{}) (err error) {
	if e := p.getValue(pv); e != nil {
		err = errutil.New(p, e)
	}
	return
}

func (p RefProp) getValue(pv interface{}) (err error) {
	if pdst := r.ValueOf(pv); pdst.Kind() != r.Ptr {
		err = errutil.New("expected a pointer out value")
	} else {
		dst := pdst.Elem()
		// a bool is an indicator of state lookup
		if dst.Kind() != r.Bool {
			err = p.obj.objects.coerce(dst, p.Value)
		} else if p.Kind() == r.Bool {
			dst.SetBool(p.Bool())
		} else {
			c, idx := p.Int(), p.index
			match := c == int64(idx)
			dst.SetBool(match)
		}
	}
	return
}

func (p RefProp) SetValue(v interface{}) (err error) {
	if e := p.setValue(v); e != nil {
		err = errutil.New(p, e)
	}
	return
}
func (p RefProp) setValue(v interface{}) (err error) {
	if b, ok := v.(bool); ok {
		err = p.setBool(b)
	} else {
		v := r.ValueOf(v)
		enumish := p.Kind() == r.Int && v.Kind() == r.String
		if !enumish {
			err = p.obj.objects.coerce(p.Value, v)
		} else {
			if choices := enum.Enumerate(p.Type()); len(choices) == 0 {
				err = errutil.New("not an enumerated field")
			} else {
				choice := ident.IdOf(v.String())
				if i, ok := enum.ChoiceToIndex(choice, choices); !ok {
					err = errutil.New("set unknown choice", choice, choices)
				} else {
					err = coerceValue(p.Value, r.ValueOf(i))
				}
			}
		}
	}
	return
}

func (p RefProp) setBool(b bool) (err error) {
	if p.Kind() == r.Bool {
		err = CoerceValue(p.Value, b)
	} else if b {
		// if setting the choice to true, then we are setting the field to the choice.
		err = CoerceValue(p.Value, p.index)
	} else {
		// we have to try to generate an opposite value.
		if choices := enum.Enumerate(p.Type()); len(choices) == 0 {
			err = errutil.New("not an enumerated field")
		} else if cnt := len(choices); cnt > 2 {
			err = errutil.New("no opposite value. too many choices")
		} else {
			// idx= 0; 2-(0+1)=1
			// idx= 1; 2-(1+1)=0
			// ret can be out of range for 1 length enums
			idx := 2 - (p.index + 1)
			err = CoerceValue(p.Value, idx)
		}
	}
	return
}
