package debug

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// DoNothing implements Execute, but .... does nothing.
type DoNothing struct {
	Reason string
}

func (*DoNothing) Compose() composer.Spec {
	return composer.Spec{
		Group: "exec",
		Desc:  "Do Nothing: Statement which does nothing.",
	}
}

func (DoNothing) Execute(rt.Runtime) error { return nil }
