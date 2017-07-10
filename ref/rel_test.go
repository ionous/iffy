package ref

import (
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/stretchr/testify/suite"
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

type GremlinRocks struct {
	BeneficentOne *Gremlin `if:"rel:one"`
	Pet           *Rock    `if:"rel:many"`
	Nickname      string
}

func TestRelationSuite(t *testing.T) {
	suite.Run(t, new(RelSuite))
}

type RelSuite struct {
	suite.Suite
	classes   *Classes
	objects   *Objects
	relations *Relations
}

func (assert *RelSuite) SetupTest() {
	assert.classes = NewClasses()
	assert.objects = NewObjects(assert.classes)
	assert.relations = NewRelations(assert.classes, assert.objects)
}

func (assert *RelSuite) TestRegistration() {
	type TooFew struct {
		A *Gremlin `if:"rel:one"`
	}
	type TooMany struct {
		A, B, C *Gremlin `if:"rel:one"`
	}
	type TooManyManys struct {
		A, B, C *Gremlin `if:"rel:many"`
	}
	type OneToOne struct {
		A, B *Gremlin `if:"rel:one"`
	}
	type ManyToOne struct {
		A *Gremlin `if:"rel:many"`
		B *Gremlin `if:"rel:one"`
	}
	type OneToMany struct {
		A *Gremlin `if:"rel:one"`
		B *Gremlin `if:"rel:many"`
	}
	type ManyToMany struct {
		A, B *Gremlin `if:"rel:many"`
	}
	type JustRight struct {
		*OneToOne
		*ManyToOne
		*OneToMany
		*ManyToMany
	}

	unique.RegisterTypes(
		unique.PanicTypes(assert.classes),
		(*Gremlin)(nil))

	reg := unique.PanicTypes(assert.relations)
	assert.Panics(func() {
		unique.RegisterTypes(reg, (*TooFew)(nil))
	})
	assert.Panics(func() {
		unique.RegisterTypes(reg, (*TooMany)(nil))
	})
	assert.Panics(func() {
		unique.RegisterTypes(reg, (*TooManyManys)(nil))
	})
	assert.NotPanics(func() {
		unique.RegisterBlocks(reg, (*JustRight)(nil))
	})
	_, tooFew := reg.FindType("TooFew")
	assert.False(tooFew)

	for i := 0; i < 4; i++ {
		t := index.Type(i)
		if r, ok := assert.relations.GetRelation(t.String()); assert.True(ok) {
			assert.Equal(t, r.GetType(), t.String())
		}
	}
}

// test a simple one to many relation
func (assert *RelSuite) TestOneToMany() {
	unique.RegisterTypes(unique.PanicTypes(assert.classes),
		(*Gremlin)(nil),
		(*Rock)(nil))

	unique.RegisterTypes(unique.PanicTypes(assert.relations),
		(*GremlinRocks)(nil))

	unique.RegisterValues(unique.PanicValues(assert.objects),
		&Gremlin{Name: "claire"},
		&Rock{Name: "loofa"},
		&Rock{Name: "rocky"},
		&Rock{Name: "petra"})

	Object := func(name string) rt.Object {
		ret, ok := assert.objects.GetObject(name)
		if !ok {
			assert.Fail("couldnt find obect", name)
		}
		return ret
	}

	if gr, ok := assert.relations.GetRelation("GremlinRocks"); assert.True(ok) {
		assert.Equal("$gremlinRocks", gr.GetId())
		assert.Equal(index.OneToMany, gr.GetType())

		if rel, e := gr.Relate(Object("claire"), Object("loofa")); assert.NoError(e) {
			gr := gr.(*RefRelation)
			contains := func(i index.Column, n string) bool {
				_, ok := gr.table.Index[i].FindFirst(0, n)
				return ok
			}
			var bene, pet rt.Object
			if e := rel.GetValue("beneficent one", &bene); assert.NoError(e) {
				assert.Equal(Object("claire"), bene)
				if e := rel.GetValue("pet", &pet); assert.NoError(e) {
					assert.Equal(Object("loofa"), pet)
					assert.True(contains(index.Primary, "$claire"))
					//
					if e := rel.SetValue("pet", nil); assert.NoError(e) {
						assert.False(contains(index.Primary, "$claire"))
						// testing: nil pointers return error
						assert.Error(rel.GetValue("pet", &pet))
					}
					//
					if e := rel.SetValue("pet", Object("petra")); assert.NoError(e) {
						assert.True(contains(index.Primary, "$claire"))
						// NOTE: if you clear via Relates(), you wont see the change from a relation.
						if _, e := gr.Relate(Object("claire"), Object("loofa")); assert.NoError(e) {
							assert.True(contains(index.Secondary, "$loofa"))
							assert.True(contains(index.Secondary, "$petra"))
						}
					}
				}
			}
		}
	}
}
