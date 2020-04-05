package pattern

import (
	"github.com/ionous/iffy/rt"
)

// Filters control whether the eval associated with a pattern will trigger.
// It follows the behavior specified by rt.GetAllTrue
type Filters []rt.BoolEval
