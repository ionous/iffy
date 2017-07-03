package ref_test

import (
	"github.com/ionous/iffy/index"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
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

func TestRel(t *testing.T) {
	suite.Run(t, new(RelSuite))
}

type RelSuite struct {
	suite.Suite
	test *testing.T
}

func (assert *RelSuite) SetupTest() {}

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

	classes := make(ref.Classes)
	unique.RegisterType(classes, (*Gremlin)(nil))
	rel := ref.MakeRelations(classes)
	reg := unique.PanicRegistry(rel)
	assert.Panics(func() {
		unique.RegisterType(reg, (*TooFew)(nil))
	})
	assert.Panics(func() {
		unique.RegisterType(reg, (*TooMany)(nil))
	})
	assert.Panics(func() {
		unique.RegisterType(reg, (*TooManyManys)(nil))
	})
	assert.NotPanics(func() {
		unique.RegisterBlock(reg, (*JustRight)(nil))
	})
	_, tooFew := reg.FindType("TooFew")
	assert.False(tooFew)

	for i := 0; i < 4; i++ {
		t := index.Type(i)
		if r, ok := rel.GetRelation(t.String()); assert.True(ok) {
			assert.Equal(t, r.GetType(), t.String())
		}
	}
}
