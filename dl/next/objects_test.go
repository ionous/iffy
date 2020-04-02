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
		this.Text: base.Text,
		that.Text: derived.Text,
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
		testTrue(t, &run, &IsKindOf{that, base}) //

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

func (m *modelTest) GetObject(obj, field string, pv interface{}) (err error) {
	switch field {
	case object.Exists:
		_, ok := m.clsMap[obj]
		bptr := pv.(*bool)
		*bptr = ok

	case object.Kind:
		if cls, ok := m.clsMap[obj]; !ok {
			err = errutil.New("unknown object", cls)
		} else {
			sptr := pv.(*string)
			*sptr = cls
		}
	default:
		err = errutil.New("unknown field", field)
	}
	return
}

func (m *modelTest) IsCompatible(childKind, parentKind string) bool {
	return parentKind == childKind || parentKind == "base"
}
