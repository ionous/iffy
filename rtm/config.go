package rtm

import (
	"bytes"
	"github.com/ionous/iffy/event"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/pat"
	"github.com/ionous/iffy/ref/obj"
	"github.com/ionous/iffy/ref/rel"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"io"
)

type Config struct {
	classes   unique.Types
	ancestors rt.Ancestors
	events    event.EventMap
	grammar   parser.Scanner
	objects   *obj.ObjBuilder
	ruleMap   pat.Rulebook
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
func (c *Config) Rules(p pat.Contract) *Config {
	p.Sort()
	c.ruleMap = p.Rulebook
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

	//
	rel := c.rel.Build()
	var w io.Writer
	if cw := c.writer; cw != nil {
		w = cw
	} else {
		var cw bytes.Buffer
		w = &cw
	}
	//
	rtm := &Rtm{
		Types:     c.classes,
		Relations: rel,
		Ancestors: a,
		writer:    w,
		Scanner:   c.grammar,
		Events:    c.events,
		rules:     c.ruleMap,
	}
	if c.objects != nil {
		rtm.Objects = c.objects.Build(rtm)
	}
	//
	seed := c.seed
	if seed == 0 {
		seed = 1
	}
	rtm.Randomizer.Reset(seed) // FIX: time?

	return rtm
}
