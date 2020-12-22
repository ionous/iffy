package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

// these definitions mirror modeling statements in iffy.js
// right now, common noun and proper noun implement ObjectEval directly.
// it'd be nice to make them swaps just like they are in the modeling section.
// future: make.swap("noun", "{proper_noun} or {named_noun}");

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

func (op *SimpleNoun) GetObject(run rt.Runtime) (ret g.Value, err error) {
	return getObjectNamed(run, string(op.Name))
}
