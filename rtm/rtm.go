package rtm

import (
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
)

type Rtm struct {
	*ref.Classes
	*ref.Objects
	*ref.Relations
	ScopeStack
	OutputStack
	Randomizer
	Ancestors
	rt.Patterns
}

type Config struct {
	classes   *ref.Classes
	objects   *ref.Objects
	rel       *ref.Relations
	ancestors Ancestors
	patterns  pat.Patterns
	seed      int64
}

// New to initialize a runtime step-by-step.
// It can be useful for testing to leave some portions of the runtime blank.
// Classes are the only "required" element.
func New(classes *ref.Classes) *Config {
	return &Config{classes: classes}
}

func (c *Config) Objects(o *ref.Objects) *Config {
	c.objects = o
	return c
}

func (c *Config) Ancestors(a Ancestors) *Config {
	c.ancestors = a
	return c
}

func (c *Config) Relations(r *ref.Relations) *Config {
	c.rel = r
	return c
}

func (c *Config) Randomize(seed int64) *Config {
	c.seed = seed
	return c
}

func (c *Config) Patterns(p pat.Patterns) *Config {
	c.patterns = p
	return c
}

func (c *Config) Rtm() *Rtm {
	a := c.ancestors
	if a == nil {
		a = NoAncestors{}
	}
	rtm := &Rtm{
		Classes:   c.classes,
		Objects:   c.objects,
		Relations: c.rel,
		Ancestors: a,
	}
	//
	rtm.Patterns = Thunk{rtm, c.patterns}
	//
	seed := c.seed
	if seed == 0 {
		seed = 1
	}
	rtm.PushScope(scope.ModelFinder(rtm))
	rtm.Randomizer.Reset(seed) // FIX: time?
	return rtm
}

// Ancestors is compatible with the rt.Runtime
type Ancestors interface {
	GetAncestors(rt.Object) (rt.ObjectStream, error)
}
