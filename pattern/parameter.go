package pattern

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
)

type Parameter struct {
	Name string // parameter name
	From core.Assignment
}

type Parameters struct {
	Params []*Parameter
}

func (*Parameter) Compose() composer.Spec {
	return composer.Spec{
		Name:  "parameter",
		Spec:  "its {name:variable_name} is {from:assignment}",
		Group: "patterns",
	}
}

func (*Parameters) Compose() composer.Spec {
	return composer.Spec{
		Name:  "parameters",
		Spec:  " when {parameters%params+parameter}",
		Group: "patterns",
	}
}
