package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
)

// used for switching between the automatic printing of an object's name
// and a request for an object of a particular name.
// ex. {.lantern} vs. {kindOf: .lantern}
type dottedName struct {
	name string
}

// when dotted names are used as arguments to concrete functions
// 		ex. {numAsWords: .count}
// we cant know the type of the variable .count without keeping a name stack during compilation
// but we can use the existing command GetVar which implements every eval type.
func (on *dottedName) getVariableNamed() *core.GetVar {
	return &core.GetVar{
		Name:            on.getName(),
		TryTextAsObject: true,
	}
}

// when dotted names are as arguments to patterns:
// 		ex. {printPluralName: .target}
// we dont know the type of "target" ahead of time
// so we just pass it around behind the scenes as an interface.
func (on *dottedName) getFromVar() core.Assignment {
	return &core.FromVar{
		Name:            on.getName(),
		TryTextAsObject: true,
	}
}

func (on *dottedName) getName() (ret rt.TextEval) {
	if name := on.name; lang.IsCapitalized(name) {
		ret = &core.ObjectName{T(name)}
	} else {
		ret = T(name)
	}
	return
}

// when dotted names are used directly:
// 		ex {.lantern}
// first attempting to read from the name as a variable,
// and if that fails, attempting to render the name as an object.
func (on *dottedName) getPrintedName() (ret rt.TextEval) {
	return &RenderName{on.name}
}
