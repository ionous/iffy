package pattern

import (
	"github.com/ionous/iffy/rt"
)

// Filters control whether the eval associated with a pattern will trigger.
// It follows the behavior specified by rt.GetAllTrue
// FIX FIX FIX -- the filters should just be a single bool eval
// use AllTrue when its needed at the handler level yeah_
type Filters []rt.BoolEval
