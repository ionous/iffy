package rel

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
)

var Slats = []composer.Composer{
	(*Relation)(nil),
	(*Relate)(nil),
	(*RelativeOf)(nil),
	(*RelativesOf)(nil),
	(*ReciprocalOf)(nil),
	(*ReciprocalsOf)(nil),
}

func cmdError(op composer.Composer, err error) error {
	return errutil.Append(err, &core.CommandError{Cmd: op})
}
