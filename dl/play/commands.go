package play

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/rules"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/event/trigger"
	"github.com/ionous/iffy/parser"
)

type Statement interface {
	Define(*Facts) error
}

type Definitions []Statement

type Facts struct {
	Grammar         parser.AnyOf
	ObjectListeners []ListenTo
	ClassListeners  []ListenFor
	Locations       []Location
	Mandates        rules.Mandates
}

type Commands struct {
	Core    core.Commands
	Rules   rules.Commands
	Std     std.Commands
	Parser  parser.Commands
	Trigger trigger.Commands

	// Define:
	*Grammar
	*ListenTo
	*ListenFor
	*Location
	*Mandate
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
