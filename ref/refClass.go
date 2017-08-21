package ref

import (
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/ref/enum"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

// note: for speed of some of these operations, we could cache the results in a pointer to a struct pooled in a map of rtype->RefClass.
// equality would work because the pointers are the same,
// just as equaity works by aliasing r.Type directly.
type RefClass struct {
	r.Type
}

func makeClass(rtype r.Type) RefClass {
	return RefClass{rtype}
}

// GetId returns the unique identifier for this classes.
func (c RefClass) GetId() string {
	name := c.Type.Name()
	return id.MakeId(name)
}

// GetName returns a friendly name: spaces, no caps.
func (c RefClass) GetName() string {
	name := c.Type.Name()
	return lang.Lowerspace(name)
}

// String implements fmt.Stringer
func (c RefClass) String() string {
	return c.Type.String()
}

// GetParentType returns false for classes if no parent;
func (c RefClass) GetParent() (ret rt.Class, okay bool) {
	if path, ok := unique.PathOf(c.Type, "parent"); ok {
		field := c.Type.FieldByIndex(path)
		ret, okay = makeClass(field.Type), true
	}
	return
}

func (c RefClass) IsCompatible(name string) (okay bool) {
	if idn := id.MakeId(name); c.GetId() == idn {
		okay = true
	} else {
		for {
			if p, ok := c.GetParent(); !ok {
				break
			} else if p.GetId() == idn {
				okay = true
				break
			}
		}
	}
	return
}

func (c RefClass) getProperty(pid string) (ret []int) {
	fn := func(f *r.StructField, path []int) (done bool) {
		if id.MakeId(f.Name) == pid {
			ret, done = path, true
		}
		return
	}
	unique.WalkProperties(c.Type, fn)
	return
}

// when we dont know if pid is a boolean property or a choice.
// ex. struct { Openable bool } or struct { Openable }
func (c RefClass) getAmbigiousProperty(idt string) (ret []int, idx int) {
	fn := func(f *r.StructField, path []int) (done bool) {
		if id.MakeId(f.Name) == idt {
			ret, idx, done = path, -1, true
		} else if choices := enum.Enumerate(f.Type); len(choices) > 0 {
			if c, ok := choiceToIndex(idt, choices); ok {
				ret, idx, done = path, c, true
			}
		}
		return
	}
	unique.WalkProperties(c.Type, fn)
	return
}

func choiceToIndex(cid string, cs []string) (ret int, okay bool) {
	for i, q := range cs {
		if id.MakeId(q) == cid {
			ret, okay = i, true
			break
		}
	}
	return
}
