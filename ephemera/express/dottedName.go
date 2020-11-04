package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/render"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
)

// used for switching between the automatic printing of an object's name
// and a request for an object of a particular name.
// the full set of possibilities are:
//   {x}                  - print name of an object.
//   {kindOf: x}          - call a function with a local variable or object.
//   {printPluralName: x} - call a user pattern with a local variable or object.
//
// where x can be, for instance:
//    {.target}   - a local variable named target
//    {.lantern}  - a variable or object named lantern
//    {.Lantern}  - an object named lantern

type dottedName struct {
	name string
}

func (on *dottedName) flags() (ret core.TryAsNoun) {
	if name := on.name; lang.IsCapitalized(name) {
		ret = core.TryAsObject
	} else {
		ret = core.TryAsBoth
	}
	return
}

// when dotted names are used as arguments to concrete functions
// 		ex. {numAsWords: .count}
// we cant know the type of the variable .count without keeping a name stack during compilation
// but we can use the existing command GetVar which implements every eval type.
func (on *dottedName) getValueNamed() *core.GetVar {
	return &core.GetVar{
		Name:  T(on.name),
		Flags: on.flags(),
	}
}

// when dotted names are as arguments to patterns:
// 		ex. {printPluralName: .target}
// we dont know the type of "target" ahead of time
// so we just pass it around behind the scenes as an interface.
func (on *dottedName) getFromVar() (ret core.Assignment) {
	return &core.CopyFrom{
		Name:  T(on.name),
		Flags: on.flags(),
	}
}

// when dotted names are used directly:
// 		ex {.lantern}
// first attempting to read from the name as a variable,
// and if that fails, attempting to render the name as an object.
func (on *dottedName) getPrintedName() (ret rt.TextEval) {
	return &render.Name{on.name}
}
