package term

import (
	"github.com/ionous/iffy/rt"
)

// parameters in the future might have defaults...
// something similar might be used for local variables.
// we could also -- in some far future land -- code generate things.
type Preparer interface {
	Prepare(rt.Runtime, *Terms) error
	String() string
}
