package pattern

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/term"
)

// Rules contained by this package.
// fix: would it be better to list rule sets?
// the rule set elements could be used to find the individual rule types.
var Support = []interface{}{
	(*BoolRule)(nil),
	(*NumberRule)(nil),
	(*TextRule)(nil),
	(*NumListRule)(nil),
	(*TextListRule)(nil),
	(*ExecuteRule)(nil),
	//
	//(*term.Preparer)(nil),
	(*term.Terms)(nil),
	(*term.Number)(nil),
	(*term.Bool)(nil),
	(*term.Text)(nil),
	(*term.Object)(nil),
	(*term.NumList)(nil),
	(*term.TextList)(nil),
}

var Slats = []composer.Slat{
	(*DetermineAct)(nil),
	(*DetermineNum)(nil),
	(*DetermineText)(nil),
	(*DetermineBool)(nil),
	(*DetermineNumList)(nil),
	(*DetermineTextList)(nil),
	(*Arguments)(nil),
	(*Argument)(nil),
}
