package story

import "github.com/ionous/iffy/ephemera"

type StoryEnv struct {
	Recent struct {
		// Scene, Aspect, Test string
		// Nouns[]? Relation, Trait
		// string or ephemera.Named
		Nouns Nouns
		Test  ephemera.Named
	}
	Current struct {
		// eventually, a stack.
		Domain ephemera.Named
	}
}

func (n *StoryEnv) CollectTest(test ephemera.Named, during func() error) (err error) {
	lastScene := n.Current.Domain
	n.Current.Domain = test
	// the most recent test might become the last popped test value
	// ( once domains and tests are stackable )
	n.Recent.Test = test
	err = during()
	n.Current.Domain = lastScene
	return
}

type Nouns struct {
	Subjects, Objects []ephemera.Named
	Objectifying      bool // phrases discuss noun subjects by default
}

// add to the known recent nouns over the course of the passed function.
// subjects are the main focus of the sentence, often the ones mentioned first (lhs).
func (n *Nouns) CollectSubjects(fn func() error) error {
	n.Subjects = nil
	n.Objectifying = false
	return fn()
}

// add to the known recent nouns over the course of the passed function.
// objects are the support nouns in a sentence, often mentioned last (rhs).
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
