package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// Activity wraps a block of multiple execute statements
// This is primarily useful for the composer so it can display blocks in a uniform manner.
type Activity struct {
	Exe []rt.Execute
}

func (*Activity) Compose() composer.Spec {
	return composer.Spec{
		Name:  "activity",
		Group: "hidden",
		Spec:  "{exe*execute}",
	}
}

// Execute statements
func (op *Activity) Execute(run rt.Runtime) error {
	return rt.RunAll(run, op.Exe)
}
