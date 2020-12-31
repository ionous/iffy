package list

import "github.com/ionous/iffy/dl/composer"

type Edge bool

func (*Edge) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_edge",
		Uses:  "str",
		Group: "list",
		Spec:  "{atFront%front} or {atBack%back}",
		Desc:  "List Edge: Indicate elements at the front or back of a list.",
		Stub:  true,
	}
}
