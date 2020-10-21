package list

import "github.com/ionous/iffy/dl/composer"

type FrontOrBack bool

func (op FrontOrBack) Compose() composer.Spec {
	//make.str("bool", "{true} or {false}");
	return composer.Spec{
		Name:  "list_edge",
		Group: "list",
		Spec:  "{front} or {back}",
		Desc:  "List Edge: Indicate elements at the front or back of a list.",
	}
}
