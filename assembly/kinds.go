package assembly

import (
	"database/sql"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
)

type cachedKind struct {
	name   string
	parent *cachedKind
}

// String returns the name of this kind.
func (c *cachedKind) String() string {
	return c.name
}

// GetAncestors returns a comma separated string of this kind's parent, and its parent's search.
func (c *cachedKind) GetAncestors() string {
	list := make([]string, 0, 0)
	for t := c.parent; t != nil; t = t.parent {
		list = append(list, t.String())
	}
	return strings.Join(list, ",")
}

// HasAncestor returns true if req is the parent of this kind, or a parent of any parent.
func (c *cachedKind) HasAncestor(req *cachedKind) (ret bool) {
	for t := c.parent; t != nil; t = t.parent {
		if req == t {
			ret = true
			break
		}
	}
	return
}

type cachedKinds struct {
	cache cacheMap
}

// helper for cached kinds
type cacheMap map[string]*cachedKind

// Get a cachedKind of n
func (c *cachedKinds) Get(n string) (ret *cachedKind) {
	if el, ok := c.cache[n]; ok {
		ret = el
	} else {
		if c.cache == nil {
			c.cache = make(cacheMap)
		}
		el := &cachedKind{name: n}
		c.cache[n] = el
		ret = el
	}
	return
}

// work backwards from k to ensure a defined root.
// fix? revisit? is there a more db friendly way to do this?
func (c *cachedKinds) AddAncestorsOf(db *sql.DB, k string) (err error) {
	search := []*cachedKind{c.Get(k)}
	pairs := make([]*cachedKind, 0)
	//
	for len(search) > 0 {
		var req *cachedKind
		last := len(search) - 1
		req, search = search[last], search[:last]
		// for each kid paired with req in the db, register req as the kid's ancestor.
		if e := ephemera.KidsOf(db, req.name, func(k string) {
			kid := c.Get(k)
			if e, ok := tryPair(kid, req); e != nil {
				err = errutil.Append(err, e)
			} else if !ok {
				pairs = append(pairs, kid, req)
			} else {
				search = append(search, kid)
			}
		}); e != nil {
			err = e
			break
		}
	}
	// (re)try pairs until nothing matches
	for keepGoing := err == nil; keepGoing; {
		keepGoing = false // provisionally
		for i := 0; i < len(pairs); i += 2 {
			kid, req := pairs[i], pairs[i+1]
			if e, ok := tryPair(kid, req); e != nil {
				err = e
				break
			} else if ok {
				// slice out successful pairs
				pairs = append(pairs[:i], pairs[i+3:]...)
				keepGoing = true
				break
			}
		}
	}
	// log errors
	for i := 0; i < len(pairs); i += 2 {
		kid, req := pairs[i], pairs[i+1]
		e := errutil.New("couldn't add", req, req.GetAncestors(), "as ancestor of", kid, kid.GetAncestors())
		err = errutil.Append(err, e)
	}
	return
}

func tryPair(kid, req *cachedKind) (err error, okay bool) {
	// check for cycles
	if req.HasAncestor(kid) {
		err = errutil.New("cycle detected", req, req.GetAncestors(), "<>", kid, req)
	} else {
		// the kid has no existing constraints
		if kid.parent == nil {
			kid.parent = req
			okay = true
		} else {
			// the kid's existing parent exists in the hierarchy of the requested kind
			// setting req as the new parent, keeps the old parent in the hierarchy.
			if req.HasAncestor(kid.parent) {
				kid.parent = req
				okay = true
			} else {
				// it's possible that we might already be a child of req
				// ex. we might have been be a child of an earlier pair
				if kid.HasAncestor(req) {
					okay = true
				}
			}
		}
	}
	return
}
