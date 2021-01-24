package express

import (
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/render"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
)

// single dotted name supporting the printing of an object's name or
// a request for an object of a particular name.
// the full set of possibilities are:
//   {x}                  - print name of an object.
//   {kind_of: x}          - call a function with a local variable or object.
//   {print_plural_name: x} - call a user pattern with a local variable or object.
//
// where x can be, for instance:
//    {.target}   - a local variable named target
//    {.lantern}  - a variable or object named lantern
//    {.Lantern}  - an object named lantern
type dotName string

func (on dotName) flags() (ret render.TryAsNoun) {
	if name := string(on); lang.IsCapitalized(name) {
		ret = render.TryAsObject
	} else {
		ret = render.TryAsBoth
	}
	return
}

// when dotted names are used as arguments to concrete functions
// 		ex. {numAsWords: .count}
// we cant know the type of the variable .count without keeping a name stack during compilation
// but we can use the existing command Var which implements every eval type.
func (on dotName) getValueNamed() *render.RenderRef {
	return &render.RenderRef{core.Var{string(on)}, on.flags()}
}

// when dotted names are as arguments to patterns:
// 		ex. {printPluralName: .target}
// we dont know the type of "target" ahead of time
// so we just pass it around behind the scenes as an interface.
func (on dotName) getFromVar() core.Assignment {
	return on.getValueNamed()
}

// when dotted names are used directly:
// 		ex {.target} or {.Lantern} or {.text}
// first attempting to read from the name as a variable,
// and if that fails, attempting to render the name as an object.
func (on dotName) getPrintedName() rt.TextEval {
	// the render.Name function itself handles the capitalization check
	// one thing missing here: if the text in a variable is not already an id
	// this will just print the text.
	return &render.Name{string(on)}
}
