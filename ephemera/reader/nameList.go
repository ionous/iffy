package reader

import "github.com/ionous/ephemera"

type NameList struct {
	Named []ephemera.Named
}

func (n *NameList) Add(name ephemera.Named) {
	n.Named = append(n.Named, name)
}

// return the most recently specified noun, or blank if none.
func (n *NameList) Last() (ret ephemera.Named) {
	if cnt := len(n.Named); cnt > 0 {
		ret = n.Named[cnt-1]
	}
	return
}

func (n *NameList) Swap(v []ephemera.Named) []ephemera.Named {
	prev := n.Named
	n.Named = v
	return prev
}
