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
	Parent   string
	Relative locate.Containment
	Child    string
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
