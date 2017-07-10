package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// Classes maps ids to RefClass.
// Compatible with unique.TypeRegistry
type Classes struct {
	ClassMap
}
type ClassMap map[string]*RefClass

// FIX: NewClasses should live in a builder, just like event and pat
func NewClasses() *Classes {
	return &Classes{make(ClassMap)}
}

// RegisterType and all parent types.
// Compatible with unique.TypeRegistry
func (reg *Classes) RegisterType(rtype r.Type) (err error) {
	_, err = reg.RegisterClass(rtype)
	return
}

// FindType returns the originally specified type, and not the wrapper class.
// Compatible with unique.TypeRegistry.
func (reg *Classes) FindType(name string) (ret r.Type, okay bool) {
	id := id.MakeId(name)
	if a, ok := reg.ClassMap[id]; ok {
		ret, okay = a.rtype, true
	}
	return
}

// GetClass compatible with rt.Runtime
func (reg *Classes) GetClass(name string) (ret rt.Class, okay bool) {
	id := id.MakeId(name)
	ret, okay = reg.ClassMap[id]
	return
}

// GetByType for cache usage
func (reg *Classes) GetByType(rtype r.Type) (ret *RefClass, err error) {
	name := rtype.Name()
	id := id.MakeId(name)
	if cls, ok := reg.ClassMap[id]; !ok {
		err = errutil.New("class not found", name)
	} else if cls.rtype != rtype {
		err = errutil.New("class conflict", name, cls, rtype)
	} else {
		ret = cls
	}
	return
}

func (reg *Classes) RegisterClass(rtype r.Type) (ret *RefClass, err error) {
	clsid := id.MakeId(rtype.Name())
	// does the class already exist?
	if cls, exists := reg.ClassMap[clsid]; exists {
		// does the id and class match?
		if cls.rtype != rtype {
			err = errutil.New("class name needs to be unique", cls.rtype.Name(), clsid)
		} else {
			ret = cls
		}
	} else {
		// make a new class:
		cls := &RefClass{id: clsid, rtype: rtype}
		reg.ClassMap[clsid] = cls

		// parse the properties
		if ptype, pidx, props, e := MakeProperties(rtype); e != nil {
			err = e
		} else {
			cls.props = props
			if ptype == nil {
				ret = cls
			} else {
				if p, e := reg.RegisterClass(ptype); e != nil {
					err = e
				} else {
					cls.parent = p
					cls.parentIdx = pidx
					ret = cls
				}
			}
		}
	}
	return
}

func MakeProperties(rtype r.Type) (parent r.Type, parentIdx int, props []rt.Property, err error) {
	ids := make(map[string]string)

	for fw := unique.Fields(rtype); fw.HasNext(); {
		field := fw.GetNext()
		if field.Target != rtype {
			break // weve advanced to the parent
		}
		//
		if !field.IsPublic() {
			err = errutil.New("expected only exportable fields", field.Name)
			break
		} else if field.IsParent() {
			parent = field.Type
			parentIdx = field.Index
		} else {
			id := id.MakeId(field.Name)
			if was := ids[id]; len(was) > 0 {
				err = errutil.New("duplicate property was:", was, "now:", field.Name)
				break
			} else if cat, e := Categorize(field.Type); e != nil {
				err = errutil.New("error categorizing", field.Name, e)
				break
			} else {
				var p rt.Property
				base := RefProp{id, field.Index, cat}
				if cat != rt.State {
					p = &base
				} else {
					if choices, e := EnumFromField(field.StructField); e != nil {
						err = errutil.New("error enumerating", field.Name, field.Type, e)
						break
					} else {
						p = &RefEnum{base, choices}
					}
				}
				//
				props = append(props, p)
			}
		}
	}
	return
}
