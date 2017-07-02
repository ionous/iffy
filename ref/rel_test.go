package ref_test

import (
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

func (assert *RelSuite) TestSomething() {
	type TooFew struct {
		A *Gremlin `if:"rel:one"`
	}
	type TooMany struct {
		A, B, C *Gremlin `if:"rel:one"`
	}
	type TooManyManys struct {
		A, B *Gremlin `if:"rel:many"`
	}
	type JustRightOne struct {
		A, B *Gremlin `if:"rel:one"`
	}
	type JustRightMany struct {
		A *Gremlin `if:"rel:many"`
		B *Gremlin `if:"rel:one"`
	}
	reg := unique.PanicRegistry(ref.RelationRegistry{make(unique.Types)})
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
		unique.RegisterType(reg, (*JustRightOne)(nil))
	})
	assert.NotPanics(func() {
		unique.RegisterType(reg, (*JustRightMany)(nil))
	})
	_, tooFew := reg.FindType("TooFew")
	assert.False(tooFew)
	_, rightOne := reg.FindType("JustRightOne")
	assert.True(rightOne)
	_, rightMany := reg.FindType("JustRightMany")
	assert.True(rightMany)
}
