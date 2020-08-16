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
		Group: "internal",
	}
}

func (op *CommonNoun) GetText(run rt.Runtime) (ret string, err error) {
	return getObjectFullName(run, op)
}

func (op *CommonNoun) GetBool(run rt.Runtime) (ret bool, err error) {
	return getObjectExists(run, op)
}

func (op *CommonNoun) GetObjectRef(run rt.Runtime) (ret string, exact bool, err error) {
	ret, exact = string(op.Name), false
	return
}

// internal because, currently, iffy.js defines the spec.
func (*ProperNoun) Compose() composer.Spec {
	return composer.Spec{
		Group: "internal",
	}
}

func (op *ProperNoun) GetText(run rt.Runtime) (ret string, err error) {
	return getObjectFullName(run, op)
}

func (op *ProperNoun) GetBool(run rt.Runtime) (ret bool, err error) {
	return getObjectExists(run, op)
}

func (op *ProperNoun) GetObjectRef(run rt.Runtime) (ret string, exact bool, err error) {
	ret, exact = string(op.Name), false
	return
}

// note: we dont support primitive type specs yet; and even if we did,
// we'd have to move these and the rest of iffy.js into go specs.

// func (Determiner) Compose() Spec {
// 	return Spec{
// 		Desc: `A qualifying word preceding a noun.
// 		For example: the definite article ( the ), or an indefinite article ( an, some ).`,
// 		Spec: "{a}, {an}, {the}, or {other determiner%determiner}",
// 	}
// }

// func (ProperName) Compose() Spec {
// 	return Spec{
// 		Desc: `Proper Name: A name given to some specific person, place, or thing.
// Proper names are usually capitalized. For example, maybe: 'Haruki', 'Jane', or 'Toronto'.`,
// 	}
// }

// func (CommonName) Compose() Spec {
// 	return Spec{
// 		Desc: `Common Name: A generalized name given to some specific item, place, or thing.
//     Common names are usually not capitalized. For example, maybe: 'table', 'chair', or 'dog park'.`,
// 	}
// }
