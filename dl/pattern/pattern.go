package pattern

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/term"
)

// Rules contained by this package.
// fix: would it be better to list rule sets?
// the rule set elements could be used to find the individual rule types.
var Support = []interface{}{
	(*ExecuteRule)(nil),
	//
	//(*term.Preparer)(nil),
	(*term.Number)(nil),
	(*term.Bool)(nil),
	(*term.Text)(nil),
	(*term.Record)(nil),
	(*term.Object)(nil),
	(*term.NumList)(nil),
	(*term.TextList)(nil),
	(*term.RecordList)(nil),
}

var Slats = []composer.Composer{
	(*Determine)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &core.CommandError{Cmd: op})
}

func cmdErrorCtx(op composer.Composer, ctx string, err error) error {
	return errutil.Append(err, &core.CommandError{Cmd: op, Ctx: ctx})
}
