package next

import (
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
)

func TestObjects(t *testing.T) {
	this, that, nothing := &Text{"this"}, &Text{"that"}, &Text{"nothing"}
	base, derived := &Text{"base"}, &Text{"derived"}

	run := modelTest{clsMap: map[string]string{
		// objects:
		this.Text: base.Text,
		that.Text: derived.Text,
		// hierarchy:
		base.Text:    base.Text,
		derived.Text: derived.Text + "," + base.Text,
	}}

	t.Run("exists", func(t *testing.T) {
		testTrue(t, &run, &Exists{this})
		testTrue(t, &run, &IsNot{&Exists{nothing}})
	})
	t.Run("kind of", func(t *testing.T) {
		if cls, e := rt.GetText(&run, &KindOf{this}); e != nil {
			t.Fatal(e)
		} else if cls != base.Text {
			t.Fatal("unexpected", cls)
		}
	})
	t.Run("is kind of", func(t *testing.T) {
		testTrue(t, &run, &IsKindOf{this, base})
		testTrue(t, &run, &IsKindOf{that, base})

		testTrue(t, &run, &IsKindOf{that, derived})
		testTrue(t, &run, &IsNot{&IsKindOf{this, derived}})
	})
	t.Run("is exact kind of", func(t *testing.T) {
		testTrue(t, &run, &IsExactKindOf{this, base})
		testTrue(t, &run, &IsNot{&IsExactKindOf{that, base}})
		testTrue(t, &run, &IsExactKindOf{that, derived})
		testTrue(t, &run, &IsNot{&IsExactKindOf{this, derived}})
	})

}

type modelTest struct {
	baseRuntime
	clsMap map[string]string
}

func (m *modelTest) GetField(name, field string) (ret interface{}, err error) {
	switch field {
	case object.Exists:
		_, ok := m.clsMap[name]
		ret = ok

	case object.Kind:
		if cls, ok := m.clsMap[name]; !ok {
			err = errutil.New("unknown", name)
		} else {
			ret = cls
		}

	case object.Kinds:
		if cls, ok := m.clsMap[name]; !ok {
			err = errutil.New("unknown", name)
		} else if path, ok := m.clsMap[cls]; !ok {
			err = errutil.New("unknown class", cls)
		} else {
			ret = path
		}

	default:
		err = errutil.New("unknown field", field)
	}
	return
}
