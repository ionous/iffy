package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// these definitions mirror modeling statements in iffy.js
// right now, common noun and proper noun implement ObjectEval directly.
// it'd be nice to make them swaps just like they are in the modeling section.
// future: make.opt("noun", "{proper_noun} or {named_noun}");

// SimpleNoun implements ObjectEval
type SimpleNoun struct {
	Determiner Determiner // determiners are used for modeling hints
	Name       NounName
}

type Determiner string
type NounName string

// internal because, currently, iffy.js defines the spec.
func (*SimpleNoun) Compose() composer.Spec {
	return composer.Spec{
		Name:  "named_noun",
		Spec:  "{determiner} {name%noun_name}",
		Group: "internal",
	}
}

// can be used as text, returns the object.id
// func (op *SimpleNoun) GetText(run rt.Runtime) (ret string, err error) {
// 	return op.GetObjectValue(run)
// }

func (op *SimpleNoun) GetObjectValue(run rt.Runtime) (ret rt.Value, err error) {
	return getObjectInexactly(run, string(op.Name))
}
