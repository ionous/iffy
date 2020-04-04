package composer

type Spec struct {
	Name, Spec, Group, Desc string
	Locals                  []string
}

type Specification interface {
	Compose() Spec
}
