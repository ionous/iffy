package builder

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/spec"
)

// Build the memento tree to finialize values
func Build(n *Memento) (ret interface{}, err error) {
	if n.spec != nil {
		if e := buildSpec(n); e != nil {
			err = e
		} else {
			ret = n.spec
		}
	} else if n.specs != nil {
		if e := buildSpecs(n); e != nil {
			err = e
		} else {
			ret = n.specs
		}
	} else if n.val != nil {
		if e := buildValue(n); e != nil {
			err = e
		} else {
			ret = n.val
		}
	} else {
		err = errutil.New(n.pos, "memento is empty?")
	}
	return
}

func buildValue(n *Memento) (err error) {
	if !n.kids.IsEmpty() {
		err = errutil.New("values should not have children")
	}
	return
}

func buildSpec(n *Memento) (err error) {
	var keys map[string]*Memento
	for _, kid := range n.kids.list {
		if v, e := Build(kid); e != nil {
			err = errutil.New(kid.pos, e)
			break
		} else {
			// verify keyword expectations:
			if k := kid.key; len(k) > 0 {
				if keys == nil {
					keys = map[string]*Memento{k: kid}
				} else if was := keys[k]; was != nil {
					err = errutil.New(kid.pos, "duplicate keyword detected", was.pos)
					break
				} else {
					keys[k] = kid
				}

				if e := n.spec.Assign(k, v); e != nil {
					err = errutil.New(kid.pos, e)
					break
				}
			} else {
				if keys != nil {
					err = errutil.New(kid.pos, "positional arguments should appear before all keyword arguments")
					break
				}
				if e := n.spec.Position(v); e != nil {
					err = errutil.New(kid.pos, e)
					break
				}
			}
		}
	}
	return
}

func buildSpecs(n *Memento) (err error) {
	for _, kid := range n.kids.list {
		if len(kid.key) > 0 {
			err = errutil.New(kid.pos, "array elements shouldnt use keyword arguments")
			break
		} else if res, e := Build(kid); e != nil {
			err = errutil.New(kid.pos, e)
			break
		} else if spec, ok := res.(spec.Spec); !ok {
			err = errutil.New(kid.pos, "only commands should be used in arrays")
			break
		} else if e := n.specs.AddElement(spec); e != nil {
			err = errutil.New(kid.pos, e)
			break
		}
	}
	return
}
