package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// RefRelation describes a single relationship "archetype"
type RefRelation struct {
	id        string
	cls       *RefClass // FIX: can we get by with r.Type?
	props     [2]RefInfo
	table     index.Table
	relations *Relations
}

type RefInfo struct {
	id        string
	cls       *RefClass
	fieldPath []int
}

func (reg *Relations) newRelation(relid string, cls *RefClass) (ret *RefRelation, err error) {
	var t int
	var i int
	var props [2]RefInfo
	choices := []index.Type{
		index.OneToOne, index.ManyToOne,
		index.OneToMany, index.ManyToMany}
	//
	add := func(field unique.FieldInfo, flag int) (err error) {
		if !(i < 2) {
			err = errutil.New("too many fields", field.Name)
		} else if field.Type.Kind() != r.Ptr {
			err = errutil.New("expected a pointer", field.Name)
		} else {
			elem := field.Type.Elem()
			if cls, ok := reg.objectClasses.GetClass(elem.Name()); !ok {
				err = errutil.New("unknown class", field.Name, elem.Name())
			} else {
				p := &props[i]
				p.id = id.MakeId(field.Name)
				p.fieldPath = append(field.Path, field.Index)
				p.cls = cls.(*RefClass)
				t |= flag << uint(i)
				i++
			}
		}
		return
	}
	flags := map[string]int{
		"one":  0,
		"many": 1,
	}
OutOfLoop:
	for fw := unique.Fields(cls.rtype); fw.HasNext(); {
		field := fw.GetNext()
		tag := unique.Tag(field.Tag)
		if rel, ok := tag.Find("rel"); ok {
			if flag, ok := flags[rel]; !ok {
				err = errutil.New("unknown relation", rel)
				break OutOfLoop
			} else if e := add(field, flag); e != nil {
				err = e
				break OutOfLoop
			}
		}
	}
	if err == nil {
		ret = &RefRelation{relid, cls, props, index.MakeTable(choices[t]), reg}
	}
	return
}

// GetId returns the unique identifier for this types.
func (rel *RefRelation) GetId() string {
	return rel.id
}

// GetType of the relation: one-to-one to many-to-many.
func (rel *RefRelation) GetType() index.Type {
	return rel.table.Type()
}

// Relate defines a connection between two objects.
func (rel *RefRelation) Relate(src, dst rt.Object, data interface{}) (prev interface{}, err error) {
	var donotforgetdata, orprev interface{}
	if rec, e := rel.relate(src, dst); e != nil {
		err = e
	} else {
		ret = &RefRelative{rec, rel}
	}
	return
}

func reduce(obj rt.Object) (id string, okay bool) {
	if obj != nil {
		id = obj.GetId()
		okay = len(id) > 0
	}
	return
}

// returns a relation record
func (rel *RefRelation) relate(src, dst rt.Object) (ret *RefObject, err error) {
	if s, ok := reduce(src); !ok {
		err = errutil.Fmt("primary object is anonymous", src.GetClass())
	} else if d, ok := reduce(dst); !ok {
		err = errutil.Fmt("secondary object is anonymous", dst.GetClass())
	} else {
		onInsert := func(old interface{}) (ret interface{}, err error) {
			if ref, ok := old.(*RefObject); !ok {
				ret = rel.relations.objects.newObject(rel.cls)
			} else {
				ret = ref
			}
			return
		}
		if rec, e := rel.table.Relate(s, d, onInsert); e != nil {
			err = e
		} else {
			ret = rec.(*RefObject)
		}
	}
	return
}
