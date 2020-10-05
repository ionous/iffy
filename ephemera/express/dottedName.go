package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
)

// var dots dottedName

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
		Name:            on.getName(false),
		TryTextAsObject: true,
	}
}

// when dotted names are as arguments to patterns:
// 		ex. {printPluralName: .target}
// we dont know the type of "target" ahead of time
// so we just pass it around behind the scenes as an interface.
func (on *dottedName) getFromVar() core.Assignment {
	return &core.FromVar{
		Name:            on.getName(false),
		TryTextAsObject: true,
	}
}

// when dotted names are used directly:
// 		ex {.lantern}
func (on *dottedName) getPrintedName() rt.TextEval {
	return &core.Buffer{core.NewActivity(
		&pattern.DetermineAct{
			"printName",
			// on.name is already setup with an object id lookup
			// so from text is being given an object.
			pattern.NewArgs(&core.FromText{on.getName(true)}),
		})}
}

func (on *dottedName) getName(b bool) (ret rt.TextEval) {
	if t := on.name; lang.IsCapitalized(t) {
		ret = &core.ObjectName{T(t)}
	} else if b {
		ret = on.getVariableNamed()
	} else {
		ret = T(t)
	}
	return
}
