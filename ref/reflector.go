package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

type Classes map[string]*RefClass

func (cs Classes) GetClass(name string) (ret rt.Class, okay bool) {
	id := id.MakeId(name)
	ret, okay = cs[id]
	return
}

func (cs Classes) FindType(name string) (ret r.Type, okay bool) {
	id := id.MakeId(name)
	if a, ok := cs[id]; ok {
		ret, okay = a.rtype, true
	}
	return
}

func (cs Classes) RegisterType(rtype r.Type) (err error) {
	_, err = cs.addClass(rtype.Elem())
	return
}

func (cs Classes) addClass(rtype r.Type) (ret *RefClass, err error) {
	clsid := id.MakeId(rtype.Name())
	// does the class already exist?
	if cls, exists := cs[clsid]; exists {
		// does the id and class match?
		if cls.rtype != rtype {
			err = errutil.New("class name needs to be unique", cls.rtype.Name(), clsid)
		} else {
			ret = cls
		}
	} else {
		// make a new class:
		cls := &RefClass{id: clsid, rtype: rtype}
		cs[clsid] = cls

		// parse the properties
		if ptype, pidx, props, e := MakeProperties(rtype, &cls.meta); e != nil {
			err = e
		} else {
			cls.props = props
			if ptype == nil {
				ret = cls
			} else {
				if p, e := cs.addClass(ptype); e != nil {
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

// questions:
// how far do you want to take pointers?
// suck in their references if set? what about their class definitions?
func (cs Classes) MakeModel(instances []interface{}) (ret Objects, err error) {
	objects := make(Objects)

	// note: building copies up front because:
	// 1. error checking
	// 2. simplify coding
	// 3. basis for inception-style code generation
	// but it does have its down sides....
	// might be better to use the new iterator instead
	// with a simple isPod check/er.
	for i, inst := range instances {
		rval := r.ValueOf(inst).Elem()
		// create the class first:
		if cls, e := cs.addClass(rval.Type()); e != nil {
			err = e
			break
		} else if idField := cls.findId(); len(idField) < 0 {
			err = errutil.New("couldnt find id for", cls)
		} else {
			name := rval.FieldByName(idField)
			if !name.IsValid() || name.Kind() != r.String {
				err = errutil.New("instance needs an valid id", i, rval)
				break
			} else {
				objid := id.MakeId(name.String())
				// println("making", objid)
				if orig, ok := objects[objid]; ok {
					err = errutil.New("instance needs unique name", name, orig, objid)
					break
				}
				inst := &RefInst{id: objid, rval: rval, cls: cls}
				objects[objid] = inst
			}
		}
	}
	if err == nil {
		ret = objects
	}
	return
}

func MakeProperties(rtype r.Type, pdata *unique.Metadata) (parent r.Type, parentIdx int, props []rt.Property, err error) {
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
		} else {
			unique.MergeMetadata(field.StructField, pdata)
			//
			if field.IsParent() {
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
	}
	return
}
