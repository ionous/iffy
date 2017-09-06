package ref_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
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
	objects ref.ObjectMap
	first   *BaseClass
	second  *DerivedClass
}

func (assert *ObjectSuite) SetupTest() {
	// reset the registries every time:
	classes := make(unique.Types)
	unique.RegisterTypes(unique.PanicTypes(classes),
		(*BaseClass)(nil),
		(*DerivedClass)(nil))

	objects := ref.NewObjects()
	first := &BaseClass{Name: "first", State: Yes, Labeled: true}
	second := &DerivedClass{BaseClass{Name: "second", State: Maybe}}
	first.Object = &second.BaseClass
	second.Object = first
	assert.first = first
	assert.second = second
	unique.RegisterValues(unique.PanicValues(objects), first, second)
	assert.objects = objects.Build()
}

func (assert *ObjectSuite) TestRegistration() {
	if n, ok := assert.objects.GetObject("first"); assert.True(ok) {
		assert.Equal(ident.IdOf("$first"), n.Id())
	}
	if d, ok := assert.objects.GetObject("second"); assert.True(ok) {
		assert.Equal(ident.IdOf("$second"), d.Id())
	}
}

func (assert *ObjectSuite) TestStateAccess() {
	test := func(obj, prop string, value bool) {
		if n, ok := assert.objects.GetObject(obj); assert.True(ok) {
			var res bool
			if e := getValue(n, prop, &res); assert.NoError(e) {
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

func getValue(obj rt.Object, name string, pv interface{}) (err error) {
	if p, ok := obj.Property(name); !ok {
		err = errutil.New("unknown property", name)
	} else {
		err = p.GetValue(pv)
	}
	return
}

// SetValue sets the named property in the passed object to the value.
func setValue(obj rt.Object, name string, v interface{}) (err error) {
	if p, ok := obj.Property(name); !ok {
		err = errutil.New("unknown property", name)
	} else {
		err = p.SetValue(v)
	}
	return
}

func (assert *ObjectSuite) TestStateSet() {
	if n, ok := assert.objects.GetObject("first"); assert.True(ok) {
		var res bool
		// start with yes, it should be true
		getValue(n, "yes", &res)
		if assert.True(res) {
			// try to change the value to maybe
			setValue(n, "maybe", true)
			// yes should now be false.
			getValue(n, "yes", &res)
			if assert.False(res) {
				// and maybe should now be true
				getValue(n, "maybe", &res)
				assert.True(res)
				// try to change states in an illegal way:
				e := setValue(n, "maybe", false)
				assert.Error(e)

				// add verify it didnt change:
				getValue(n, "maybe", &res)
				assert.True(res)
			}
		}
		//
		getValue(n, "yes", &res)
		if assert.False(res) {
			//
			e := setValue(n, "state", "yes")
			if assert.NoError(e) {
				//
				e := getValue(n, "yes", &res)
				if assert.NoError(e) {
					assert.True(res)
				}
			}
		}
	}
	// check, change, and check the labeled bool.
	toggle := func(name, prop string, goal bool) {
		if n, ok := assert.objects.GetObject(name); assert.True(ok) {
			var res bool
			getValue(n, prop, &res)
			if assert.NotEqual(goal, res, "initial value") {
				setValue(n, prop, goal)
				getValue(n, prop, &res)
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
			if e := getValue(n, v.name, v.pv); assert.NoError(e) {
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
