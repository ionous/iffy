package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// these definitions mirror modeling statements in iffy.js
// right now, common noun and proper noun implement ObjectRef directly.
// it'd be nice to make them swaps just like they are in the modeling section.
// future: make.opt("noun", "{proper_noun} or {common_noun}");

// CommonNoun implements ObjectRef
type CommonNoun struct {
	Determiner Determiner // determiners are used for modeling hints
	Name       CommonName
}

// ProperNoun implements ObjectRef
type ProperNoun struct {
	Name ProperName
}

type ProperName string
type CommonName string
type Determiner string

// internal because, currently, iffy.js defines the spec.
func (*CommonNoun) Compose() composer.Spec {
	return composer.Spec{
		Name:  "common_noun",
		Spec:  "{determiner} {name%common_name}",
		Group: "internal",
	}
}

// can be used as text, returns the object.id
func (op *CommonNoun) GetText(run rt.Runtime) (ret string, err error) {
	return op.GetObjectRef(run)
}

func (op *CommonNoun) GetObjectRef(run rt.Runtime) (retId string, err error) {
	return getObjectInexactly(run, string(op.Name))
}

// internal because, currently, iffy.js defines the spec.
func (*ProperNoun) Compose() composer.Spec {
	return composer.Spec{
		Name:  "proper_noun",
		Spec:  "{name%proper_name}",
		Group: "internal",
	}
}

// can be used as text, returns the object.id
func (op *ProperNoun) GetText(run rt.Runtime) (ret string, err error) {
	return op.GetObjectRef(run)
}

func (op *ProperNoun) GetObjectRef(run rt.Runtime) (retId string, err error) {
	return getObjectInexactly(run, string(op.Name))
}
