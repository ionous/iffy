package composer

type Spec struct {
	Name, Spec, Group, Desc string
}

type SpecInterface interface {
	Compose() Spec
}
