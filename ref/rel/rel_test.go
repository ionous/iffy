package rel

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	testify "github.com/stretchr/testify/assert"
	r "reflect"
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
	classes := make(unique.Types)
	unique.PanicTypes(classes,
		(*Gremlin)(nil),
		(*Rock)(nil))
	relbuilder := NewRelations()
	relbuilder.NewRelation("GremlinRocks", index.OneToMany)

	relations := relbuilder.Build()
	//
	Object := func(name string) rt.Object {
		return ObjectMock(ident.IdOf(name))
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
		if c, e := gr.Relate(claire, loofa, index.NoData); assert.NoError(e) && assert.False(c) {
			assert.True(contains(gr.Table.Primary, "$claire"))
			assert.True(contains(gr.Table.Secondary, "$loofa"))
			assert.False(contains(gr.Table.Secondary, "$petra"))

			if c, e := gr.Relate(nil, loofa, nil); assert.NoError(e) && assert.True(c) {
				assert.False(contains(gr.Table.Primary, "$claire"))
				assert.False(contains(gr.Table.Secondary, "$loofa"))
				assert.False(contains(gr.Table.Secondary, "$petra"))

				if c, e := gr.Relate(claire, petra, index.NoData); assert.NoError(e) && assert.False(c) {
					assert.True(contains(gr.Table.Primary, "$claire"))
					assert.False(contains(gr.Table.Secondary, "$loofa"))
					assert.True(contains(gr.Table.Secondary, "$petra"))

					if c, e := gr.Relate(claire, loofa, index.NoData); assert.NoError(e) && assert.False(c) {
						assert.True(contains(gr.Table.Primary, "$claire"))
						assert.True(contains(gr.Table.Secondary, "$loofa"))
						assert.True(contains(gr.Table.Secondary, "$petra"))
					}
				}
			}
		}
	}
}

type ObjectMock ident.Id

func (om ObjectMock) Id() ident.Id {
	return ident.Id(om)
}
func (om ObjectMock) Type() r.Type {
	return r.TypeOf(om)
}
func (om ObjectMock) GetValue(prop string, pv interface{}) error {
	return errutil.New("object mock doesnt support get value")
}
func (om ObjectMock) SetValue(prop string, v interface{}) error {
	return errutil.New("object mock doesnt support set value")
}
