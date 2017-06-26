package reflector

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref"
	r "reflect"
)

type RefModel struct {
	objects      map[string]*RefInst
	linearObject []*RefInst
	classes      map[string]*RefClass
	linearClass  []*RefClass
}

// Model provides the starting point for all game objects and classes.
func (m *RefModel) NumClass() int {
	return len(m.linearClass)
}

func (m *RefModel) ClassNum(i int) ref.Class {
	return m.linearClass[i]
}

func (m *RefModel) GetClass(name string) (ret ref.Class, okay bool) {
	id := id.MakeId(name)
	ret, okay = m.classes[id]
	return
}

func (m *RefModel) NumObject() int {
	return len(m.linearObject)
}

func (m *RefModel) ObjectNum(i int) (ret ref.Object) {
	return m.linearObject[i]
}

func (m *RefModel) GetObject(name string) (ret ref.Object, okay bool) {
	id := id.MakeId(name)
	ret, okay = m.objects[id]
	return
}

func (m *RefModel) NewObject(class string) (ret ref.Object, err error) {
	id := id.MakeId(class)
	if cls, ok := m.classes[id]; !ok {
		err = errutil.New("no such class", class)
	} else {
		inst := r.New(cls.rtype)
		ret = &RefInst{rval: inst.Elem(), cls: cls}
	}
	return
}
