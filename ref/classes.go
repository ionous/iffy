package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

type ClassBuilder struct {
	ClassMap
}

func NewClasses() *ClassBuilder {
	return &ClassBuilder{make(ClassMap)}
}

// RegisterType and all parent types.
// Compatible with unique.TypeRegistry
func (cb *ClassBuilder) RegisterType(rtype r.Type) (err error) {
	_, err = cb.RegisterClass(rtype)
	return
}

func (cb *ClassBuilder) RegisterClass(rtype r.Type) (ret *RefClass, err error) {
	clsid := id.MakeId(rtype.Name())
	// does the class already exist?
	if cls, exists := cb.ClassMap[clsid]; exists {
		// does the id and class match?
		if cls.rtype != rtype {
			err = errutil.New("class name needs to be unique", cls.rtype.Name(), clsid)
		} else {
			ret = cls
		}
	} else {
		// make a new class:
		cls := &RefClass{id: clsid, rtype: rtype}
		cb.ClassMap[clsid] = cls

		// parse the properties
		if ptype, pidx, props, e := MakeProperties(rtype); e != nil {
			err = errutil.New(rtype, e)
		} else {
			cls.props = props
			if ptype == nil {
				ret = cls
			} else {
				if p, e := cb.RegisterClass(ptype); e != nil {
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

func MakeProperties(rtype r.Type) (parent r.Type, parentIdx []int, props []rt.Property, err error) {
	p := _Properties{}
	if e := p.MakeProperties(rtype, nil); e != nil {
		err = e
	} else {
		parent = p.parent
		parentIdx = p.parentIdx
		props = p.props
	}
	return
}

type _Properties struct {
	parent    r.Type
	parentIdx []int
	props     []rt.Property
}

func (k *_Properties) MakeProperties(rtype r.Type, base []int) (err error) {
	ids := make(map[string]string)
	for i := 0; i < rtype.NumField(); i++ {
		field := rtype.Field(i)
		//
		if !unique.IsPublic(field) {
			err = errutil.New("expected only exportable fields", field.Name)
			break
		} else if k.parent == nil && IsParentField(&field) {
			k.parent = field.Type
			k.parentIdx = append(base, i)
		} else {
			id := id.MakeId(field.Name)
			if was := ids[id]; len(was) > 0 {
				err = errutil.New("duplicate property was:", was, "now:", field.Name)
				break
			} else if unique.IsEmbedded(field) {
				if e := k.MakeProperties(field.Type, append(base, i)); e != nil {
					err = e
					break
				}
			} else if cat, e := Categorize(field.Type); e != nil {
				err = errutil.New("error categorizing", field.Name, e)
				break
			} else {
				var p rt.Property
				base := RefProp{id, append(base, i), cat}
				if cat != rt.State {
					p = &base
				} else {
					if choices, e := EnumFromField(&field); e != nil {
						err = errutil.New("error enumerating", field.Name, field.Type, e)
						break
					} else {
						p = &RefEnum{base, choices}
					}
				}
				//
				k.props = append(k.props, p)
			}

		}
	}
	return
}

func IsParentField(f *r.StructField) bool {
	_, ok := unique.Tag(f.Tag).Find("parent")
	return ok
}
