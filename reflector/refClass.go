package reflector

import (
	"github.com/ionous/iffy/ref"
	r "reflect"
)

type RefClass struct {
	id        string
	rtype     r.Type // mainly for equality building tests
	meta      Metadata
	parent    *RefClass
	parentIdx int            // index of parent aggregate in rtype; valid if parent!= nil
	props     []ref.Property // RefProp, RefEnum, etc.
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

// String implements fmt.Stringer
func (c *RefClass) String() string {
	return c.id
}

// GetId returns the unique identifier for this classes.
func (c *RefClass) GetId() string {
	return c.id
}

// GetParentType returns false for classes if no parent;
func (c *RefClass) GetParent() (ref.Class, bool) {
	return c.parent, c.parent != nil
}

// Number returns the number of indexable properties.
// The number of available properties for a given Class never changes at runtime.
func (c *RefClass) NumProperty() int {
	return len(c.props)
}

// PropertyNum returns the indexed property.
// Panics if the index is greater than Number.
func (c *RefClass) PropertyNum(i int) ref.Property {
	return c.props[i]
}

// GetProperty by name.
func (c *RefClass) GetProperty(name string) (ret ref.Property, okay bool) {
	id := MakeId(name)
	for _, p := range c.props {
		if p.GetId() == id {
			ret, okay = p, true
			break
		}
	}
	return
}

// GetPropertyByChoice evaluates all properties to find an enumeration which can store the passed choice
func (c *RefClass) GetPropertyByChoice(choice string) (ref.Property, bool) {
	id := MakeId(choice)
	r, _, _ := c.getPropertyByChoice(id)
	return r, r != nil
}

func (c *RefClass) getPropertyByChoice(id string) (ret *RefEnum, path []int, value int) {
	for {
		for _, p := range c.props {
			if p, ok := p.(*RefEnum); ok {
				if i := p.choiceToIndex(id); i >= 0 {
					ret = p
					path = append(path, p.fieldIdx)
					value = i
					break
				}
			}
		}
		if ret != nil && c.parent == nil {
			break
		} else {
			c = c.parent
			path = append(path, c.parentIdx)
		}
	}
	return
}

func (c *RefClass) IsCompatible(name string) (okay bool) {
	id := MakeId(name)
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
