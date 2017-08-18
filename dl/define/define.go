package define

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/std"
	"github.com/ionous/iffy/parser"
	"github.com/ionous/iffy/pat/patspec"
)

type Statement interface {
	Define(*Facts) error
}

type Definitions []Statement

type Facts struct {
	Grammar   parser.AnyOf
	Pattern   []patspec.PatternSpecs
	Listeners []When
	Locations []Location
}

type Commands struct {
	Core     core.Commands
	Patterns patspec.Commands
	Std      std.Commands
	Parser   parser.Commands

	Define struct {
		*The
		*When
		*Pattern
		*Grammar
		*Plural
		*Location
	}
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
