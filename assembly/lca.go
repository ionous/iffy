package assembly

import "strings"

// by default, lca contains kind --
// but it might not once new kinds are merged into it.
type hierarchy struct {
	name    string
	parents string   // mdl hierarchy of kind
	lca     []string // root is on the right.
	valid   bool     // valid if lca is a named part
}

// normalize name, parents into an array of kinds.
func (h *hierarchy) getAncestry() (ret []string) {
	parts := []string{h.name}
	if len(h.parents) > 0 {
		ret = append(parts, strings.Split(h.parents, ",")...)
	} else {
		ret = parts
	}
	return
}

func (h *hierarchy) set(lca []string) {
	h.lca, h.valid = lca, len(lca) > 0
}

func (h *hierarchy) update(other *hierarchy) {
	if h.name != other.name {
		cmp, lca := findOverlap(h.lca, other.getAncestry())
		h.name = other.name
		h.valid = cmp != 0
		h.lca = lca
	}
}
