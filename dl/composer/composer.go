package composer

// Spec for display in composer
type Spec struct {
	Name, Spec, Group, Desc, Uses string
	Stub                          bool // generate a custom loading struct.
	Locals                        []string
	Fluent                        *Fluency
}

type Composer interface {
	Compose() Spec
}

type Fluency struct {
	Name  string // if empty, it will use the first syllable of the type name
	RunIn bool   // is the first parameter unlabeled
}
type Slot struct {
	Name  string
	Type  interface{} // nil instance, ex. (*core.Comparator)(nil)
	Desc  string
	Group string // display group(s)
}
