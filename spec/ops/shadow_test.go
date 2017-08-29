package ops_test

import (
	"github.com/ionous/iffy/dl/core" // for interesting evals
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	. "github.com/ionous/iffy/tests" // BaseClass
	"github.com/ionous/sliceOf"
	testify "github.com/stretchr/testify/assert"
	r "reflect"
	"testing"
)

func TestShadows(t *testing.T) {
	assert := testify.New(t)
	classes := make(unique.Types)
	ops := ops.NewOps(classes)

	unique.RegisterBlocks(unique.PanicTypes(ops),
		(*core.Commands)(nil))

	unique.RegisterTypes(unique.PanicTypes(ops.ShadowTypes),
		(*BaseClass)(nil))
	//
	var root struct {
		Num    rt.NumberEval
		Object rt.ObjectEval
	}
	if c, ok := ops.NewBuilder(&root); assert.True(ok) {
		// FIX: without the cmd -- it doesnt error.
		// FIX: and what about using the same param twice?
		c.Cmd("add", 1, 2)
		if c.Cmd("BaseClass").Begin() {
			c.Param("Num").Cmd("add", 1, 2)
			c.Param("Text").Val("3")
			c.Param("Nums").Val(sliceOf.Float(1, 2, 3))
			c.Param("Texts").Val(sliceOf.String("1", "2", "3"))
			c.Param("State").Val(Maybe) // Note: this turns State into a NumEval
			c.Param("Labeled").Cmd("is not", false)
			c.Param("Object").Val("other")
			c.Param("Objects").Val(sliceOf.String("base", "other"))
			c.End()
		}
		//
		if e := c.Build(); e != nil {
			t.Fatal(e)
		}
	}

	objects := ref.NewObjects()
	base, other := &BaseClass{Name: "base"}, &BaseClass{Name: "other"}
	unique.RegisterValues(unique.PanicValues(objects),
		base, other,
	)
	var lines printer.Lines
	run := rtm.New(classes).Objects(objects).Writer(&lines).Rtm()
	// "shadow class tests.BaseClass couldn't create object"
	if obj, e := root.Object.GetObject(run); assert.NoError(e) {
		vals := map[string]struct{ match, fail interface{} }{
			"Num":     {3, 5},
			"Text":    {"3", "5"},
			"Object":  {other, base},
			"Nums":    {sliceOf.Float(1, 2, 3), sliceOf.Float(3, 2, 1)},
			"Texts":   {sliceOf.String("1", "2", "3"), sliceOf.String("3")},
			"Objects": {[]*BaseClass{base, other}, []*BaseClass{base}},
			"State":   {Maybe, Yes},
			"Labeled": {true, false},
		}
		for name, test := range vals {
			cp := r.New(r.ValueOf(test.match).Type()).Elem()
			if e := obj.GetValue(name, cp.Addr().Interface()); !assert.NoError(e) {
				break
			} else if !testify.ObjectsAreEqualValues(test.match, cp.Interface()) {
				t.Fatal("failed to match", name)
				break
			} else if testify.ObjectsAreEqualValues(test.fail, cp.Interface()) {
				t.Fatal("failed to match", name)
				break
			}
		}
	}
}
