package std_test

import (
	"github.com/ionous/iffy/dl/core"
	. "github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/dl/text"
	"github.com/ionous/iffy/pat/patbuilder"
	"github.com/ionous/iffy/pat/patspec"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/spec/ops"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

// FIX: things limiting reuse across tests:
// . concrete object pointers;
//   alt: get by name and return a nothing object
// . many "builders" use map, you cant re-make that map just by calling Rtm.New;
//   alt: revisit builders
// . object needs Objects for get/pointer, which ties objects to classes early;
//   alt: get rid of pointers/object lookup.
func xTestGroup(t *testing.T) {
	assert := testify.New(t)

	classes := ref.NewClasses()
	unique.RegisterBlocks(unique.PanicTypes(classes),
		(*Classes)(nil))

	objects := ref.NewObjects(classes)
	unique.RegisterValues(unique.PanicValues(objects),
		thingList...)
	ops := ops.NewOps()
	unique.RegisterBlocks(unique.PanicTypes(ops),
		(*core.Commands)(nil),
		(*text.Commands)(nil),
		(*Commands)(nil),
		(*patspec.Commands)(nil),
	)
	unique.RegisterBlocks(unique.PanicTypes(ops.ShadowTypes),
		(*Patterns)(nil),
	)
	//
	patterns, e := patbuilder.NewPatternMaster(ops, classes,
		(*Patterns)(nil)).Build(BuildPatterns)
	assert.NoError(e)

	var lines printer.Lines
	run := rtm.New(classes).Objects(objects).Patterns(patterns).Writer(&lines).Rtm()

	prn := &PrintNondescriptObjects{&core.Objects{nameList}}
	if e := prn.Execute(run); assert.NoError(e) {
		t.Log(lines.Lines())
	}
}

var thingMap = map[string]*Thing{
	// some unnamed things
	// this relies on the internal means of naming unnamed objects
	"thing#1": &Thing{},
	"thing#2": &Thing{},
	// a named thing
	"pen": &Thing{
		Kind: Kind{Name: "pen"},
	},
	// a thing with a printed name
	"sword": &Thing{
		Kind: Kind{Name: "sword", PrintedName: "plastic sword"},
	},
}

var thingList = func(src map[string]*Thing) (ret []interface{}) {
	for _, v := range src {
		ret = append(ret, v)
	}
	return
}(thingMap)

var nameList = func(src map[string]*Thing) (ret []string) {
	for n, _ := range src {
		ret = append(ret, n)
	}
	return
}(thingMap)
