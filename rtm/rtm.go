package rtm

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/pat/rule"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"github.com/ionous/iffy/rt/scope"
	"io"
)

type Rtm struct {
	unique.Types
	*ref.Objects
	ref.Relations
	ScopeStack
	io.Writer
	Randomizer
	rt.Ancestors
	pat.Patterns
	Plurals
}

// GetPatterns mainly for testing.
func (rtm *Rtm) GetPatterns() *pat.Patterns {
	return &rtm.Patterns
}

// GetClass with the passed name.
func (rtm *Rtm) GetClass(name string) (ret rt.Class, okay bool) {
	id := ident.IdOf(name)
	if cls, ok := rtm.Types[id]; ok {
		ret, okay = cls, ok
	}
	return
}

type Config struct {
	classes   unique.Types
	objects   *ref.ObjBuilder
	rel       ref.RelationBuilder
	ancestors rt.Ancestors
	patterns  *pat.Patterns
	seed      int64
	writer    io.Writer
}

// New to initialize a runtime step-by-step.
// It can be useful for testing to leave some portions of the runtime blank.
func New(classes unique.Types) *Config {
	return &Config{classes: classes}
}

func (c *Config) Objects(o *ref.ObjBuilder) *Config {
	c.objects = o
	return c
}

func (c *Config) Ancestors(a rt.Ancestors) *Config {
	c.ancestors = a
	return c
}

func (c *Config) Relations(r ref.RelationBuilder) *Config {
	c.rel = r
	return c
}

func (c *Config) Randomize(seed int64) *Config {
	c.seed = seed
	return c
}

func (c *Config) Rules(p rule.Rules) *Config {
	p.Sort()
	c.patterns = &p.Patterns
	return c
}

func (c *Config) Writer(w io.Writer) *Config {
	c.writer = w
	return c
}

func (c *Config) Rtm() *Rtm {
	a := c.ancestors
	if a == nil {
		a = NoAncestors{}
	}
	var objects *ref.Objects
	if c.objects != nil {
		objects = c.objects.Build()
	}
	//
	rel := c.rel.Build()
	var w io.Writer
	if cw := c.writer; cw != nil {
		w = cw
	} else {
		var cw printer.Lines
		w = &cw
	}
	rtm := &Rtm{
		Types:     c.classes,
		Objects:   objects,
		Relations: rel,
		Ancestors: a,
		Writer:    w,
	}
	if c.patterns != nil {
		rtm.Patterns = *c.patterns
	}
	//
	seed := c.seed
	if seed == 0 {
		seed = 1
	}
	rtm.PushScope(scope.ModelFinder(rtm))
	rtm.Randomizer.Reset(seed) // FIX: time?

	return rtm
}
