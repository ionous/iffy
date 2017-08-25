package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref/unique"
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

func (cb *ClassBuilder) RegisterClass(rtype r.Type) (ret RefClass, err error) {
	clsid := id.MakeId(rtype.Name())
	// does the class already exist?
	if cls, exists := cb.ClassMap[clsid]; exists {
		// does the id and class match?
		if cls.Type != rtype {
			err = errutil.New("class name needs to be unique", cls.Type.Name(), clsid)
		} else {
			ret = cls
		}
	} else {
		// make a new class:
		cls := makeClass(rtype)
		cb.ClassMap[clsid] = cls

		// parse the properties, including parents....
		// look  for conflicting properties, etc.
		ids := make(map[string]string)
		if !unique.WalkProperties(rtype, func(field *r.StructField, path []int) (done bool) {
			// register parent classes
			if IsParentField(field) {
				if _, e := cb.RegisterClass(field.Type); e != nil {
					err, done = e, true
				}
			} else {
				id := id.MakeId(field.Name)
				if was, exists := ids[id]; !exists {
					ids[id] = field.Name
				} else {
					err = errutil.New("duplicate property was:", was, "now:", field.Name)
					done = true
				}
			}
			return
		}) {
			ret = cls
		}
	}
	return
}

func IsParentField(f *r.StructField) bool {
	_, ok := unique.Tag(f.Tag).Find("parent")
	return ok
}
