package ref

import (
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

type Gremlin struct {
	Name string `if:"id,plural:base classes"`
	// pets [] *Rock
}

type Rock struct {
	Name string `if:"id,plural:base classes"`
	// BeneficentOne *Gremlin
}

// test a simple one to many relation
func TestOneToMany(t *testing.T) {
	assert := testify.New(t)
	classes := NewClasses()
	unique.RegisterTypes(unique.PanicTypes(classes),
		(*Gremlin)(nil),
		(*Rock)(nil))

	objbuilder := NewObjects(classes)
	unique.RegisterValues(unique.PanicValues(objbuilder),
		&Gremlin{Name: "claire"},
		&Rock{Name: "loofa"},
		&Rock{Name: "rocky"},
		&Rock{Name: "petra"})

	relbuilder := NewRelations()
	relbuilder.NewRelation("GremlinRocks", index.OneToMany)

	// this test doesnt use runtime, so build manually
	objects := objbuilder.Build()
	relations := relbuilder.Build()
	//
	Object := func(name string) rt.Object {
		ret, ok := objects.GetObject(name)
		if !ok {
			assert.Fail("couldnt find object", name)
		}
		return ret
	}

	if gr, ok := relations.GetRelation("GremlinRocks"); assert.True(ok) {
		gr := gr.(*RefRelation)
		contains := func(i index.Index, n string) bool {
			_, ok := i.FindFirst(0, n)
			return ok
		}
		assert.Equal(index.OneToMany, gr.GetType())

		claire, loofa, petra := Object("claire"), Object("loofa"), Object("petra")
		assert.False(contains(gr.Table.Primary, "$claire"))
		assert.False(contains(gr.Table.Secondary, "$loofa"))
		assert.False(contains(gr.Table.Secondary, "$petra"))
		//
		if c, e := gr.Relate(claire, loofa, index.NoData); assert.NoError(e) && assert.True(c) {
			assert.True(contains(gr.Table.Primary, "$claire"))
			assert.True(contains(gr.Table.Secondary, "$loofa"))
			assert.False(contains(gr.Table.Secondary, "$petra"))

			if c, e := gr.Relate(nil, loofa, nil); assert.NoError(e) && assert.True(c) {
				assert.False(contains(gr.Table.Primary, "$claire"))
				assert.False(contains(gr.Table.Secondary, "$loofa"))
				assert.False(contains(gr.Table.Secondary, "$petra"))

				if c, e := gr.Relate(claire, petra, index.NoData); assert.NoError(e) && assert.True(c) {
					assert.True(contains(gr.Table.Primary, "$claire"))
					assert.False(contains(gr.Table.Secondary, "$loofa"))
					assert.True(contains(gr.Table.Secondary, "$petra"))

					if c, e := gr.Relate(claire, loofa, index.NoData); assert.NoError(e) && assert.True(c) {
						assert.True(contains(gr.Table.Primary, "$claire"))
						assert.True(contains(gr.Table.Secondary, "$loofa"))
						assert.True(contains(gr.Table.Secondary, "$petra"))
					}
				}
			}
		}
	}
}
