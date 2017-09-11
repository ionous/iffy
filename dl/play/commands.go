package play

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/dl/story"
	"github.com/ionous/iffy/event/trigger"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/pat/rule"
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
	Mandates        rule.Mandates
}

type Commands struct {
	Core    core.Commands
	Rules   rule.Commands
	Std     std.Commands
	Story   story.Commands
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
