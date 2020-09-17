package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/rt"
)

// used for switching between the automatic printing of an object's name
// and a request for an object of a particular name.
// ex. {.lantern} vs. {kindOf: .lantern}
type dottedName struct {
	name rt.TextEval // searches for an object by name and returns an id.
	// the Name might be a chain of text evals: .my.lastQuip
}

// when dotted names are used as arguments to concrete functions
// 		ex. {numAsWords: .count}
// we cant know the type of the variable .count without keeping a name stack during compilation
// but we can use the existing command GetVar which implements every eval type.
func (on dottedName) getVariableNamed() *core.GetVar {
	return &core.GetVar{
		on.name,
	}
}

// when dotted names are as arguments to patterns:
// 		ex. {printPluralName: .target}
// we dont know the type of "target" ahead of time
// so we just pass it around behind the scenes as an interface.
func (on dottedName) getFromVar() core.Assignment {
	return &core.FromVar{
		on.name,
	}
}

// when dotted names are used directly:
// 		ex {.lantern}
func (on dottedName) getPrintedName() rt.TextEval {
	return &core.Buffer{core.NewActivity(
		&pattern.DetermineAct{
			"printAName",
			pattern.NewArgs(&core.FromText{on.name}),
		})}
}
