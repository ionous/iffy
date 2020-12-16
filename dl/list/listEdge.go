package list

import "github.com/ionous/iffy/dl/composer"

type Front bool

func (*Front) Compose() composer.Spec {
	return composer.Spec{
		Name:  "list_edge",
		Uses:  "str",
		Group: "list",
		Spec:  "{front} or {back}",
		Desc:  "List Edge: Indicate elements at the front or back of a list.",
		Stub:  true,
	}
}
