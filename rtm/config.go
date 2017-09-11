package rtm

import (
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/pat/rule"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/rel"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/printer"
	"io"
)

type Config struct {
	classes   unique.Types
	ancestors rt.Ancestors
	events    event.EventMap
	grammar   parser.Scanner
	objects   *obj.ObjBuilder
	patterns  *pat.Patterns
	rel       rel.RelationBuilder
	seed      int64
	writer    io.Writer
}

// New to initialize a runtime step-by-step.
// It can be useful for testing to leave some portions of the runtime blank.
func New(classes unique.Types) *Config {
	return &Config{classes: classes}
}

func (c *Config) Ancestors(a rt.Ancestors) *Config {
	c.ancestors = a
	return c
}
func (c *Config) Events(events event.EventMap) *Config {
	c.events = events
	return c
}

func (c *Config) Grammar(r *parser.AnyOf) *Config {
	c.grammar = r
	return c
}

func (c *Config) Objects(o *obj.ObjBuilder) *Config {
	c.objects = o
	return c
}

func (c *Config) Relations(r rel.RelationBuilder) *Config {
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
	var objects obj.ObjectMap
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
	//
	rtm := &Rtm{
		Types:     c.classes,
		ObjectMap: objects,
		Relations: rel,
		Ancestors: a,
		Writer:    w,
		Scanner:   c.grammar,
		Events:    c.events,
	}
	if c.patterns != nil {
		rtm.Patterns = *c.patterns
	}
	//
	seed := c.seed
	if seed == 0 {
		seed = 1
	}
	rtm.Randomizer.Reset(seed) // FIX: time?

	return rtm
}
