package core

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/object"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
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
		if e := testTrue(t, &run, &ObjectExists{this}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &IsNotTrue{&ObjectExists{nothing}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("kind_of", func(t *testing.T) {
		if cls, e := safe.GetText(&run, &KindOf{this}); e != nil {
			t.Fatal(e)
		} else if cls.String() != base.Text {
			t.Fatal("unexpected", cls)
		}
	})
	t.Run("is_kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &IsKindOf{this, base.Text}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &IsKindOf{that, base.Text}); e != nil {
			t.Fatal(e)
		}

		if e := testTrue(t, &run, &IsKindOf{that, derived.Text}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &IsNotTrue{&IsKindOf{this, derived.Text}}); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("is_exact_kind_of", func(t *testing.T) {
		if e := testTrue(t, &run, &CompareText{&KindOf{this}, &EqualTo{}, base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{&KindOf{that}, &NotEqualTo{}, base}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{&KindOf{that}, &EqualTo{}, derived}); e != nil {
			t.Fatal(e)
		}
		if e := testTrue(t, &run, &CompareText{&KindOf{this}, &NotEqualTo{}, derived}); e != nil {
			t.Fatal(e)
		}
	})
}

func named(n string) *Text {
	return &Text{n}
}

type modelTest struct {
	baseRuntime
	clsMap map[string]string
}

func (m *modelTest) GetField(target, field string) (ret g.Value, err error) {
	switch target {
	case object.Value:
		if _, ok := m.clsMap[field]; !ok {
			err = g.UnknownObject(field)
		} else {
			ret = &objValue{model: m, name: field}
		}
	default:
		err = g.UnknownField{target, field}
	}
	return
}

type objValue struct {
	g.PanicValue
	model *modelTest
	name  string
}

func (j *objValue) Affinity() affine.Affinity {
	return affine.Object
}

func (j *objValue) FieldByName(field string) (ret g.Value, err error) {
	switch m := j.model; field {
	case object.Kind:
		if cls, ok := m.clsMap[j.name]; !ok {
			err = g.UnknownField{j.name, field}
		} else {
			ret = g.StringOf(cls)
		}

	case object.Kinds:
		if cls, ok := m.clsMap[j.name]; !ok {
			err = g.UnknownField{j.name, field}
		} else if path, ok := m.clsMap[cls]; !ok {
			err = errutil.New("modelTest: unknown class", cls)
		} else {
			ret = g.StringOf(path)
		}

	default:
		err = g.UnknownField{j.name, field}
	}
	return
}
