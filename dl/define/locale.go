package define

import (
	"github.com/ionous/iffy/dl/locate"
)

// type Relation struct {
// 	Name string
// 	Type index.Type
// }
// type Value struct {
// 	Obj  string
// 	Prop string
// 	Val  interface{}
// }
// type State struct {
// 	Obj  string
// 	Name string
// }
type Location struct {
	Parent string
	Locale
	Child string
}

// Locale, temporary. will replace with enum branch changes
type Locale interface {
	Locale() locate.Containment
}

type Supports struct{}

func (Supports) Locale() locate.Containment {
	return locate.Supports
}

type Contains struct{}

func (Contains) Locale() locate.Containment {
	return locate.Contains
}

type Wears struct{}

func (Wears) Locale() locate.Containment {
	return locate.Wears
}

type Carries struct{}

func (Carries) Locale() locate.Containment {
	return locate.Carries
}

type Holds struct{}

func (Holds) Locale() locate.Containment {
	return locate.Holds
}

// func (r *Relation) Assert(f **Facts) (nil error) {
// 	f.Relations = append(f.Relations, *r)
// 	return
// }
// func (v *Value) Assert(f **Facts) (nil error) {
// 	f.Values = append(f.Values, *v)
// 	return
// }
// func (s *State) Assert(f **Facts) (nil error) {
// 	f.Values = append(f.Values, Value{s.Obj, s.Name, true})
// 	return
// }
func (l *Location) Define(f *Facts) (nil error) {
	// FIX: check for parent-child loops.
	f.Locations = append(f.Locations, *l)
	return
}
