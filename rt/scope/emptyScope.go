package scope

import (
	"github.com/ionous/iffy/rt"
)

// EmptyScope allows use as a perpetually erroring scope.
type EmptyScope struct{}

func (EmptyScope) GetVariable(n string) (rt.Value, error) {
	return nil, rt.UnknownVariable(n)
}
