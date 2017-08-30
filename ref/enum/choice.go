package enum

import (
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	r "reflect"
)

// ChoiceToIndex returns the index of the passed id in the passed set of choices.
func ChoiceToIndex(id ident.Id, cs []string) (ret int, okay bool) {
	for i, q := range cs {
		if ident.IdOf(q) == id {
			ret, okay = i, true
			break
		}
	}
	return
}

// PropertyPath mimics class.PropertyPath, searching the passed type for a field with the passed id, except: it also tests the possible values of enumerated fields.
// ex. PropertyPath( TypeOf(Door), "Openable" )
// matches:
// struct Door { Openable bool }, and
// struct Door { OpenableState } with:
//
// go:generate stringer -type=OpenableState
// type OpenableState int

// const (
// 	Openable OpenableState = iota
// 	NotOpenable
// )
//
func PropertyPath(rtype r.Type, name string) (ret []int, idx int) {
	pid := ident.IdOf(name)
	fn := func(f *r.StructField, path []int) (done bool) {
		if ident.IdOf(f.Name) == pid {
			ret, idx, done = path, -1, true
		} else if choices := Enumerate(f.Type); len(choices) > 0 {
			if c, ok := ChoiceToIndex(pid, choices); ok {
				ret, idx, done = path, c, true
			}
		}
		return
	}
	unique.WalkProperties(rtype, fn)
	return
}
