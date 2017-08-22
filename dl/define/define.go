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
	Listeners []ListenTo
	Locations []Location
	Mandates  rule.Mandates
}

type Commands struct {
	Core   core.Commands
	Rules  rule.Commands
	Std    std.Commands
	Parser parser.Commands

	// Define:
	*Grammar
	*ListenTo
	*Location
	*Mandate
	// Locale:
	*Supports
	*Contains
	*Wears
	*Carries
	*Holds
	// EventOptions
	*Capture
	*RunAfter
	*TargetOnly
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
