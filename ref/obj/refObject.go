package obj

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/enum"
	"github.com/ionous/iffy/ref/prop"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

type RefObject struct {
	id    ident.Id   // id of the object, blank if anonymous.
	value r.Value    // stores the concrete value. ex. Rock, not *Rock.
	run   rt.Runtime // necessary for translating ids to objects.
}

// MakeObject wraps the passed value as an anonymous object.
func MakeObject(id ident.Id, i interface{}, run rt.Runtime) rt.Object {
	rval, e := unique.ValuePtr(i)
	if e != nil {
		panic(e)
	}
	return RefObject{id, rval, run}
}

// Id returns the unique identifier for this Object.
// Blank for anonymous and temporary objects.
func (n RefObject) Id() ident.Id {
	return n.id
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
func (n RefObject) Property(name string) (ret Property, okay bool) {
	rtype := n.Type()
	if path, idx := enum.PropertyPath(rtype, name); len(path) > 0 {
		pf := prop.MakeField(rtype.FieldByIndex(path), n.value.FieldByIndex(path))
		if idx < 0 {
			ret, okay = pf, true
		} else {
			ret, okay = prop.MakeState(pf, idx), true
		}
	}
	return
}

// GetValue sets the value of the passed pointer to the value of the named property.
func (n RefObject) GetValue(prop string, pv interface{}) (err error) {
	if p, ok := n.Property(prop); !ok {
		err = errutil.New(n.propN(prop), "unknown property")
	} else {
		src := p.Value()
		dst := r.ValueOf(pv)
		rt.ScopeBlock(n.run, n, func() {
			if e := n.run.Pack(dst, src); e != nil {
				err = errutil.New(n.propN(prop), "cant unpack, because", e)
			}
		})
	}
	return
}

/// SetValue sets the named property to the passed value.
func (n RefObject) SetValue(prop string, v interface{}) (err error) {
	if v == nil {
		panic(errutil.New(n.propN(prop), "is nil"))
	}
	if p, ok := n.Property(prop); !ok {
		err = errutil.New(n.propN(prop), "unknown property")
	} else {
		dst := r.New(p.Type()) // create a new destination for the value.
		src := r.ValueOf(v)
		rt.ScopeBlock(n.run, n, func() {
			if e := n.run.Pack(dst, src); e != nil {
				err = errutil.New(n.propN(prop), "cant run", dst.Type(), "from", src.Type(), "because", e)
			} else {
				err = p.SetValue(dst.Elem())
			}
		})
	}
	return
}

func (n RefObject) propN(p string) string { return n.String() + "." + p }
