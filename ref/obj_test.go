package ref

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
	. "github.com/ionous/iffy/tests"
	"github.com/ionous/sliceOf"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
)

func TestObjectSuite(t *testing.T) {
	suite.Run(t, new(ObjectSuite))
}

type ObjectSuite struct {
	suite.Suite
	m Objects
	b *BaseClass
	d *DerivedClass
}

func (t *ObjectSuite) TearDownTest() {
	errutil.Panic = false
}
func (t *ObjectSuite) SetupTest() {
	errutil.Panic = true
	b := &BaseClass{Name: "first", State: Yes, Labeled: true}
	d := &DerivedClass{BaseClass{Name: "second", State: Maybe}}
	cs := make(Classes)
	if m, e := cs.MakeModel(sliceOf.Interface(b, d)); t.NoError(e) {
		t.b = b
		t.d = d
		t.m = m
	}
}

func (t *ObjectSuite) TestDerivation() {
	if n, ok := t.m.GetObject("first"); t.True(ok) {
		t.Equal("$first", n.GetId())
		cls := n.GetClass()
		t.NotNil(cls)
		t.Equal("$baseClass", cls.GetId())
		parent, ok := cls.GetParent()
		t.Nil(parent)
		t.False(ok)
	}
	if d, ok := t.m.GetObject("second"); t.True(ok) {
		t.Equal("$second", d.GetId())
		cls := d.GetClass()
		t.NotNil(cls)
		t.Equal("$derivedClass", cls.GetId())
		if parent, ok := cls.GetParent(); t.True(ok) {
			t.Equal("$baseClass", parent.GetId())
		}
	}
}

func (t *ObjectSuite) TestStateAccess() {
	test := func(obj, prop string, value bool) {
		if n, ok := t.m.GetObject(obj); t.True(ok) {
			var res bool
			if e := n.GetValue(prop, &res); t.NoError(e) {
				t.Equal(value, res, strings.Join(sliceOf.String(obj, prop), " "))
			}
		}
	}

	test("first", "yes", true)
	test("first", "no", false)
	test("first", "maybe", false)
	test("first", "labeled", true)
	//
	test("second", "yes", false)
	test("second", "no", false)
	test("second", "maybe", true)
	test("second", "labeled", false)
}

func (t *ObjectSuite) TestStateSet() {
	if n, ok := t.m.GetObject("first"); t.True(ok) {
		var res bool
		// start with yes, it should be true
		n.GetValue("yes", &res)
		if t.True(res) {
			// try to change the value to maybe
			n.SetValue("maybe", true)
			// yes should now be false.
			n.GetValue("yes", &res)
			if t.False(res) {
				// and maybe should now be true
				n.GetValue("maybe", &res)
				t.True(res)
				t.Panics(func() {
					// try to change states in an illegal way:
					n.SetValue("maybe", false)
				})
				// add verify it didnt change:
				n.GetValue("maybe", &res)
				t.True(res)
			}
		}
	}
	toggle := func(name, prop string, goal bool) {
		if n, ok := t.m.GetObject(name); t.True(ok) {
			var res bool
			n.GetValue(prop, &res)
			if t.NotEqual(goal, res, "initial value") {
				n.SetValue(prop, goal)
				n.GetValue(prop, &res)
				t.Equal(goal, res)
			}
		}
	}
	toggle("second", "labeled", true)
	toggle("second", "labeled", false)
}

// test that normal properties are accessible
func (t *ObjectSuite) xTestPropertyAccess() {
	var expected = []struct {
		name string
		pv   interface{}
	}{
	// {"Name", new(string)},
	// {"Num", new(float64)},
	// {"Text", new(string)},
	// {"Object", rt.Pointer},
	// {"Nums", new([]float64)},
	// {"Texts", rt.Text | rt.Array},
	// {"Objects", rt.Pointer | rt.Array},
	}
	test := func(n rt.Object) {
		for _, v := range expected {
			if e := n.GetValue(v.name, v.pv); t.NoError(e) {
				//
			}
		}
	}
	if n, ok := t.m.GetObject("first"); t.True(ok) {
		//
		test(n)
	}
	if d, ok := t.m.GetObject("second"); t.True(ok) {

		//
		test(d)
	}
}
