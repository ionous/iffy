package reflector

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref"
	r "reflect"
)

// questions:
// how far do you want to take pointers?
// suck in their references if set? what about their class definitions?
func MakeModel(instances ...interface{}) (ret ref.Model, err error) {
	// tasks: walk the instances to extract some classes.
	var linearObject []*RefInst
	objects := make(map[string]*RefInst)
	cs := MakeClassSet()

	// note: building copies up front because:
	// 1. error checking
	// 2. simplify coding
	// 3. basis for inception-style code generation
	for i, inst := range instances {
		rval := r.ValueOf(inst)
		// create the class first:
		if cls, e := cs.AddClass(rval.Type()); e != nil {
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
				objid := MakeId(name.String())
				if orig, ok := objects[objid]; ok {
					err = errutil.New("instance needs unique name", name, orig, objid)
					break
				}
				inst := &RefInst{id: objid, rval: rval, cls: cls}
				objects[objid] = inst
				linearObject = append(linearObject, inst)
			}
		}
	}
	if err == nil {
		ret = &RefModel{
			objects:      objects,
			linearObject: linearObject,
			classes:      cs.classes,
			linearClass:  cs.linear,
		}
	}
	return
}

type ClassSet struct {
	linear  []*RefClass
	classes map[string]*RefClass
}

func MakeClassSet() ClassSet {
	return ClassSet{classes: make(map[string]*RefClass)}
}

func (cs *ClassSet) AddClass(rtype r.Type) (ret *RefClass, err error) {
	clsid := MakeId(rtype.Name())
	// does the class already exist?
	if cls, exists := cs.classes[clsid]; exists {
		// does the id and class match?
		if cls.rtype != rtype {
			err = errutil.New("class name needs to be unique", cls.rtype.Name(), clsid)
		}
		ret = cls
	} else {
		// make a new class:
		cls := &RefClass{id: clsid, rtype: rtype}
		cs.classes[clsid] = cls
		cs.linear = append(cs.linear, cls)

		// parse the properties
		if ptype, props, e := MakeProperties(rtype, &cls.meta); e != nil {
			err = e
		} else {
			cls.props = props
			if ptype == nil {
				ret = cls
			} else {
				if p, e := cs.AddClass(ptype); e != nil {
					err = e
				} else {
					cls.parent = p
					ret = cls
				}
			}
		}
	}
	return
}

func MakeProperties(rtype r.Type, pdata *Metadata) (parent r.Type, props []ref.Property, err error) {
	ids := make(map[string]string)
	//
	for i, cnt := 0, rtype.NumField(); i < cnt; i++ {
		field := rtype.Field(i)
		//
		if len(field.PkgPath) > 0 {
			err = errutil.New("expected only exportable fields", field.Name)
			break
		} else {
			MergeMetadata(field, pdata)
			//
			if field.Anonymous {
				if parent != nil {
					err = errutil.New("multiple parents specified", parent, field.Type)
					break
				} else {
					parent = field.Type
				}
			} else {
				id := MakeId(field.Name)
				if was := ids[id]; len(was) > 0 {
					err = errutil.New("duplicate property was:", was, "now:", field.Name)
					break
				} else if cat, e := Categorize(field.Type); e != nil {
					err = errutil.New("error categorizing", field.Name, e)
					break
				} else {
					var p ref.Property
					base := RefProp{id, i, cat}
					if cat != ref.State {
						p = &base
					} else {
						if choices, e := MakeEnum(field.Type); e != nil {
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
