package ref

import (
	"github.com/ionous/iffy/id"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	r "reflect"
)

type RefClass struct {
	id        string
	rtype     r.Type
	meta      unique.Metadata
	parent    *RefClass
	parentIdx int           // index of parent aggregate in rtype; valid if parent!= nil
	props     []rt.Property // RefProp, RefEnum, etc.
}

// findId returns the field name in rtype containing id.
func (cls *RefClass) findId() (ret string) {
	for cls != nil {
		if id, ok := cls.meta["id"]; ok {
			ret = id
			break
		}
		cls = cls.parent
	}
	return
}

func (cls *RefClass) Type() r.Type {
	return cls.rtype
}

// String implements fmt.Stringer
func (c *RefClass) String() string {
	return c.id
}

// GetId returns the unique identifier for this classes.
func (c *RefClass) GetId() string {
	return c.id
}

// GetParentType returns false for classes if no parent;
func (c *RefClass) GetParent() (rt.Class, bool) {
	return c.parent, c.parent != nil
}

// Number returns the number of indexable properties.
// The number of available properties for a given Class never changes at runtime.
func (c *RefClass) NumProperty() int {
	return len(c.props)
}

// PropertyNum returns the indexed property.
// Panics if the index is greater than Number.
func (c *RefClass) PropertyNum(i int) rt.Property {
	return c.props[i]
}

// GetProperty by name.
func (c *RefClass) GetProperty(name string) (ret rt.Property, okay bool) {
	id := id.MakeId(name)
	if p, _, ok := c.getProperty(id); ok {
		ret, okay = p, true
	}
	return
}

func (c *RefClass) getProperty(id string) (ret rt.Property, path []int, okay bool) {
	if out, ok := c.FindProperty(func(p rt.Property) (found bool) {
		if id == p.GetId() {
			ret, found = p, true
		}
		return
	}); ok {
		path = out
		okay = true
	}
	return
}

// GetPropertyByChoice evaluates all properties to find an enumeration which can store the passed choice
func (c *RefClass) GetPropertyByChoice(choice string) (rt.Property, bool) {
	id := id.MakeId(choice)
	r, _, _ := c.getPropertyByChoice(id)
	return r, r != nil
}

func (c *RefClass) FindProperty(match func(p rt.Property) bool) (path []int, okay bool) {
	var partial []int
	type getFieldIndex interface {
		getFieldIndex() int
	}
	for {
		for _, p := range c.props {
			if match(p) {
				idx := p.(getFieldIndex).getFieldIndex()
				path = append(partial, idx)
				okay = true
				break
			}
		}
		if okay || c.parent == nil {
			break
		} else {
			c, partial = c.parent, append(partial, c.parentIdx)
		}
	}
	return
}

func (c *RefClass) getPropertyByChoice(id string) (ret *RefEnum, path []int, value int) {
	if out, ok := c.FindProperty(func(p rt.Property) (found bool) {
		if p, ok := p.(*RefEnum); ok {
			if i := p.choiceToIndex(id); i >= 0 {
				ret = p
				value = i
				found = true
			}
		}
		return
	}); ok {
		path = out
	}
	return
}

func (c *RefClass) IsCompatible(name string) (okay bool) {
	id := id.MakeId(name)
	if c.id == id {
		okay = true
	} else {
		for p := c.parent; p != nil; p = p.parent {
			if p.id == id {
				okay = true
				break
			}
		}
	}
	return
}
