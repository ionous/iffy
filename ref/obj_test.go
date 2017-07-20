package ref_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
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
	objects *ref.Objects
	first   *BaseClass
	second  *DerivedClass
}

func (assert *ObjectSuite) TearDownTest() {
	errutil.Panic = false
}

func (assert *ObjectSuite) SetupTest() {
	errutil.Panic = true
	// reset the registries every time:
	classes := ref.NewClasses()
	unique.RegisterTypes(unique.PanicTypes(classes),
		(*BaseClass)(nil),
		(*DerivedClass)(nil))

	objects := ref.NewObjects(classes)
	first := &BaseClass{Name: "first", State: Yes, Labeled: true}
	second := &DerivedClass{BaseClass{Name: "second", State: Maybe}}
	first.Object = &second.BaseClass
	second.Object = first
	assert.first = first
	assert.second = second
	unique.RegisterValues(unique.PanicValues(objects), first, second)
	assert.objects = objects.Build()
}

func (assert *ObjectSuite) TestDerivation() {
	if n, ok := assert.objects.GetObject("first"); assert.True(ok) {
		assert.Equal("$first", n.GetId())
		cls := n.GetClass()
		assert.NotNil(cls)
		assert.Equal("$baseClass", cls.GetId())
		parent, ok := cls.GetParent()
		assert.Nil(parent)
		assert.False(ok)
	}
	if d, ok := assert.objects.GetObject("second"); assert.True(ok) {
		assert.Equal("$second", d.GetId())
		cls := d.GetClass()
		assert.NotNil(cls)
		assert.Equal("$derivedClass", cls.GetId())
		if parent, ok := cls.GetParent(); assert.True(ok) {
			assert.Equal("$baseClass", parent.GetId())
		}
	}
}

func (assert *ObjectSuite) TestStateAccess() {
	test := func(obj, prop string, value bool) {
		if n, ok := assert.objects.GetObject(obj); assert.True(ok) {
			var res bool
			if e := n.GetValue(prop, &res); assert.NoError(e) {
				assert.Equal(value, res, strings.Join(sliceOf.String(obj, prop), " "))
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

func (assert *ObjectSuite) TestStateSet() {
	if n, ok := assert.objects.GetObject("first"); assert.True(ok) {
		var res bool
		// start with yes, it should be true
		n.GetValue("yes", &res)
		if assert.True(res) {
			// try to change the value to maybe
			n.SetValue("maybe", true)
			// yes should now be false.
			n.GetValue("yes", &res)
			if assert.False(res) {
				// and maybe should now be true
				n.GetValue("maybe", &res)
				assert.True(res)
				assert.Panics(func() {
					// try to change states in an illegal way:
					n.SetValue("maybe", false)
				})
				// add verify it didnt change:
				n.GetValue("maybe", &res)
				assert.True(res)
			}
		}
		//
		n.GetValue("yes", &res)
		assert.False(res)
		assert.NotPanics(func() {
			n.SetValue("state", "yes")
		})
		n.GetValue("yes", &res)
		assert.True(res)
	}
	// check, change, and check the labeled bool.
	toggle := func(name, prop string, goal bool) {
		if n, ok := assert.objects.GetObject(name); assert.True(ok) {
			var res bool
			n.GetValue(prop, &res)
			if assert.NotEqual(goal, res, "initial value") {
				n.SetValue(prop, goal)
				n.GetValue(prop, &res)
				assert.Equal(goal, res)
			}
		}
	}
	toggle("second", "labeled", true)
	toggle("second", "labeled", false)
}

// test that normal properties are accessible
func (assert *ObjectSuite) TestPropertyAccess() {
	var expected = []struct {
		name string
		pv   interface{}
	}{
		{"Name", new(string)},
		{"Num", new(float64)},
		{"Text", new(string)},
		{"Object", new(*BaseClass)},
		{"Nums", new([]float64)},
		{"Texts", new([]string)},
		{"Objects", new([]*BaseClass)},
	}
	test := func(n rt.Object) {
		for _, v := range expected {
			if e := n.GetValue(v.name, v.pv); assert.NoError(e) {
				//
			}
		}
	}
	if n, ok := assert.objects.GetObject("first"); assert.True(ok) {
		test(n)
		obj := *expected[3].pv.(**BaseClass)
		assert.Equal(&assert.second.BaseClass, obj)
	}
	if d, ok := assert.objects.GetObject("second"); assert.True(ok) {
		test(d)
		obj := *expected[3].pv.(**BaseClass)
		assert.Equal(assert.first, obj)
	}
}
