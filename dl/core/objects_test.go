package core

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

// test some simple functionality of the object commands using a mock runtime
func TestObjects(t *testing.T) {
	this, that, nothing := named("this"), named("that"), named("nothing")
	base, derived := &Text{"base"}, &Text{"derived"}

	run := modelTest{clsMap: map[string]string{
		// objects:
		"this": base.Text,
		"that": derived.Text,
		// hierarchy:
		"base":    base.Text,
		"derived": derived.Text + "," + base.Text,
	}}

	t.Run("exists", func(t *testing.T) {
		testTrue(t, &run, &ObjectExists{this})
		testTrue(t, &run, &IsNotTrue{&ObjectExists{nothing}})
	})
	t.Run("kind_of", func(t *testing.T) {
		if cls, e := rt.GetText(&run, &KindOf{this}); e != nil {
			t.Fatal(e)
		} else if cls != base.Text {
			t.Fatal("unexpected", cls)
		}
	})
	t.Run("is_kind_of", func(t *testing.T) {
		testTrue(t, &run, &IsKindOf{this, base})
		testTrue(t, &run, &IsKindOf{that, base})

		testTrue(t, &run, &IsKindOf{that, derived})
		testTrue(t, &run, &IsNotTrue{&IsKindOf{this, derived}})
	})
	t.Run("is_exact_kind_of", func(t *testing.T) {
		testTrue(t, &run, &IsExactKindOf{this, base})
		testTrue(t, &run, &IsNotTrue{&IsExactKindOf{that, base}})
		testTrue(t, &run, &IsExactKindOf{that, derived})
		testTrue(t, &run, &IsNotTrue{&IsExactKindOf{this, derived}})
	})
}

func named(n string) *ObjectName {
	return &ObjectName{Name: &Text{n}}
}

type modelTest struct {
	baseRuntime
	clsMap map[string]string
}

func (m *modelTest) GetField(target, field string) (ret rt.Value, err error) {
	switch target {
	case object.Value:
		if _, ok := m.clsMap[field]; !ok {
			err = rt.UnknownObject(field)
		} else {
			ret = &objTest{model: m, name: field}
		}
	default:
		err = rt.UnknownField{target, field}
	}
	return
}

type objTest struct {
	generic.Nothing
	model *modelTest
	name  string
}

func (j *objTest) GetField(field string) (ret rt.Value, err error) {
	switch m := j.model; field {
	case object.Kind:
		if cls, ok := m.clsMap[j.name]; !ok {
			err = rt.UnknownField{j.name, field}
		} else {
			ret = generic.NewString(cls)
		}

	case object.Kinds:
		if cls, ok := m.clsMap[j.name]; !ok {
			err = rt.UnknownField{j.name, field}
		} else if path, ok := m.clsMap[cls]; !ok {
			err = errutil.New("modelTest: unknown class", cls)
		} else {
			ret = generic.NewString(path)
		}

	default:
		err = rt.UnknownField{j.name, field}
	}
	return
}
