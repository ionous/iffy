package story

import "github.com/ionous/iffy/ephemera"

// a list of recently referenced nouns
// helps simplify some aspects of importing
var storyNouns nounList

type nounList struct {
	names []ephemera.Named
}

func (n *nounList) Add(name ephemera.Named) {
	n.names = append(n.names, name)
}

// return the most recently specified noun, or blank if none.
func (n *nounList) Last() (ret ephemera.Named) {
	if cnt := len(n.names); cnt > 0 {
		ret = n.names[cnt-1]
	}
	return
}

func (n *nounList) Swap(v []ephemera.Named) []ephemera.Named {
	prev := n.names
	n.names = v
	return prev
}
