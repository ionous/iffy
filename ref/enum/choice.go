package enum

import (
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref/unique"
	r "reflect"
)

// ChoiceToIndex returns the index of the passed id in the passed set of choices.
func ChoiceToIndex(cid string, cs []string) (ret int, okay bool) {
	for i, q := range cs {
		if id.MakeId(q) == cid {
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
func PropertyPath(rtype r.Type, pid string) (ret []int, idx int) {
	fn := func(f *r.StructField, path []int) (done bool) {
		if id.MakeId(f.Name) == pid {
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
