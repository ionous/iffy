package core

import (
	"github.com/ionous/iffy/dl/composer"
)

type Argument struct {
	Name string // argument name
	From Assignment
}

type Arguments struct {
	Args []*Argument
}

func (*Argument) Compose() composer.Spec {
	return composer.Spec{
		Name:  "argument",
		Spec:  "its {name:variable_name} is {from:assignment}",
		Group: "patterns",
	}
}

func (*Arguments) Compose() composer.Spec {
	return composer.Spec{
		Name:  "arguments",
		Spec:  " when {arguments%args+argument}",
		Group: "patterns",
	}
}
