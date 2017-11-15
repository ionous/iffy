package std

import (
	"github.com/ionous/iffy/rt"
)

// Parents pattern which returns a list of the objects which enclose the targeted object, from inner most to outermost.
// See also: Children.
type Parents struct {
	Object rt.ObjectEval
}

// Children pattern which returns a list of the objects directly enclosed by targeted object.
// See also: Parents.
type Children struct {
	Object rt.ObjectEval
}
