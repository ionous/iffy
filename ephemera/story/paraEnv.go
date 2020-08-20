package story

import "github.com/ionous/iffy/ephemera"

type ParagraphEnv struct {
	Recent struct {
		// Scene, Aspect, Test string
		// Nouns[]? Relation, Trait
		// string or ephemera.Named
		Nouns Nouns
	}
	Current struct {
		Scene ephemera.Named
	}
}

type Nouns struct {
	Subjects, Objects []ephemera.Named
	Objectifying      bool // phrases discuss noun subjects by default
}

func (n *Nouns) CollectSubjects(fn func() error) error {
	n.Subjects = nil
	n.Objectifying = false
	return fn()
}

func (n *Nouns) CollectObjects(fn func() error) error {
	n.Objects = nil
	n.Objectifying = true
	err := fn()
	n.Objectifying = false
	return err
}

func (n *Nouns) pList() (ret *[]ephemera.Named) {
	if n.Objectifying {
		ret = &n.Objects
	} else {
		ret = &n.Subjects
	}
	return
}

func (n *Nouns) Add(name ephemera.Named) {
	pn := n.pList()
	(*pn) = append((*pn), name)
}

func LastNameOf(n []ephemera.Named) (ret ephemera.Named) {
	if cnt := len(n); cnt > 0 {
		ret = (n)[cnt-1]
	}
	return
}
