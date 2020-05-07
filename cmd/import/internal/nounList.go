package internal

import "github.com/ionous/iffy/ephemera"

type nounList struct {
	Named []ephemera.Named
}

func (n *nounList) Add(name ephemera.Named) {
	n.Named = append(n.Named, name)
}

// return the most recently specified noun, or blank if none.
func (n *nounList) Last() (ret ephemera.Named) {
	if cnt := len(n.Named); cnt > 0 {
		ret = n.Named[cnt-1]
	}
	return
}

func (n *nounList) Swap(v []ephemera.Named) []ephemera.Named {
	prev := n.Named
	n.Named = v
	return prev
}
