package rtm

import (
	"github.com/ionous/iffy/dl/text"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/pat/patbuilder"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/scope"
	"io"
)

type Rtm struct {
	ref.ClassMap
	*ref.Objects
	*ref.Relations
	ScopeStack
	OutputStack
	Randomizer
	Ancestors
	Thunk
	text.Plurals
}

// GetPatterns mainly for testing.
func (rtm *Rtm) GetPatterns() *pat.Patterns {
	return &rtm.Thunk.Patterns
}

type Config struct {
	classes   ref.ClassMap
	objects   *ref.ObjBuilder
	rel       *ref.RelBuilder
	ancestors Ancestors
	patterns  *patbuilder.Patterns
	seed      int64
	writer    io.Writer
}

// New to initialize a runtime step-by-step.
// It can be useful for testing to leave some portions of the runtime blank.
// Classes are the only "required" element.
func New(classes *ref.ClassBuilder) *Config {
	return &Config{classes: classes.ClassMap}
}

func (c *Config) Objects(o *ref.ObjBuilder) *Config {
	c.objects = o
	return c
}

func (c *Config) Ancestors(a Ancestors) *Config {
	c.ancestors = a
	return c
}

func (c *Config) Relations(r *ref.RelBuilder) *Config {
	c.rel = r
	return c
}

func (c *Config) Randomize(seed int64) *Config {
	c.seed = seed
	return c
}

func (c *Config) Patterns(p *patbuilder.Patterns) *Config {
	c.patterns = p
	return c
}

func (c *Config) Writer(w io.Writer) *Config {
	c.writer = w
	return c
}

func (c *Config) NewRtm() *Rtm {
	a := c.ancestors
	if a == nil {
		a = NoAncestors{}
	}
	var objects *ref.Objects
	var rel *ref.Relations
	//
	if c.objects != nil {
		objects = c.objects.Build()
		// /
		if c.rel != nil {
			c.rel.Build(objects)
		}
	}
	rtm := &Rtm{
		ClassMap:  c.classes,
		Objects:   objects,
		Relations: rel,
		Ancestors: a,
	}
	if c.patterns != nil {
		rtm.Thunk = Thunk{rtm, c.patterns.Build()}
	}
	//
	seed := c.seed
	if seed == 0 {
		seed = 1
	}
	rtm.PushScope(scope.ModelFinder(rtm))
	rtm.Randomizer.Reset(seed) // FIX: time?
	if w := c.writer; w != nil {
		rtm.PushWriter(w)
	}
	return rtm
}

// Ancestors is compatible with the rt.Runtime
type Ancestors interface {
	GetAncestors(rt.Object) (rt.ObjectStream, error)
}
