package define

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/pat/rule"
)

type Statement interface {
	Define(*Facts) error
}

type Definitions []Statement

type Facts struct {
	Grammar   parser.AnyOf
	Rules     rule.Mandates
	Listeners []When
	Locations []Location
}

type Commands struct {
	Core     core.Commands
	Patterns rule.Commands
	Std      std.Commands
	Parser   parser.Commands

	// Define:
	*The
	*When
	*Pattern
	*Grammar
	*Plural
	*Location
	// Locale:
	*Supports
	*Contains
	*Wears
	*Carries
	*Holds
}

type Classes struct {
	Core core.Classes
	Std  std.Classes
}

type Patterns struct {
	Std std.Patterns
}

func (d Definitions) Define(f *Facts) (err error) {
	for _, v := range d {
		if e := v.Define(f); e != nil {
			err = e
			break
		}
	}
	return
}
