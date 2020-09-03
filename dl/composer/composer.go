package composer

// Spec for display in composer
type Spec struct {
	Name, Spec, Group, Desc string
	Locals                  []string
}

type Slat interface {
	Compose() Spec
}

type Slot struct {
	Name  string
	Type  interface{} // nil instance, ex. (*core.Comparator)(nil)
	Desc  string
	Group string // display group(s)
}
