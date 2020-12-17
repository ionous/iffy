package story

import (
	"github.com/ionous/errutil"
)

// ex. colors are a kind of value
func (op *KindsOfAspect) ImportPhrase(k *Importer) (err error) {
	if a, e := op.Aspect.NewName(k); e != nil {
		err = e
	} else {
		k.NewAspect(a)
	}
	return
}

// ex. "cats are a kind of animal"
func (op *KindsOfKind) ImportPhrase(k *Importer) (err error) {
	if kind, e := op.PluralKinds.NewName(k); e != nil {
		err = e
	} else if parent, e := op.SingularKind.NewName(k); e != nil {
		err = e
	} else {
		k.NewKind(kind, parent)
	}
	return
}

// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func (op *KindsPossessProperties) ImportPhrase(k *Importer) (err error) {
	if kind, e := op.PluralKinds.NewName(k); e != nil {
		err = e
	} else {
		for _, n := range op.PropertyDecl {
			if e := n.ImportProperty(k, kind); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}
