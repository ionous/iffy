package std

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
)

// Player returns the object of the current viewpoint
// ChangePlayer is not yet supported.
type Player struct{}

func (*Player) GetObject(run rt.Runtime) (ret rt.Object, err error) {
	if pawn, ok := run.GetObject("pawn"); !ok {
		err = errutil.New("couldnt find pawn")
	} else {
		err = pawn.GetValue("actor", &ret)
	}
	return
}
