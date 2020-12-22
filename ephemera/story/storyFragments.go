package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

func (op *Lede) Import(k *Importer) (err error) {
	if e := k.Recent.Nouns.CollectSubjects(func() (err error) {
		for _, nn := range op.Nouns {
			if e := nn.Import(k); e != nil {
				err = errutil.Append(err, e)
			}
		}
		return
	}); e != nil {
		err = e
	} else {
		err = op.NounPhrase.Import(k)
	}
	return
}

func (op *Summary) Import(k *Importer) (err error) {
	if text, e := op.Lines.ConvertText(); e != nil {
		err = e
	} else {
		// declare the existence of the field "appearance"
		if once := "summary"; k.Once(once) {
			domain := k.gameDomain()
			things := k.NewDomainName(domain, "things", tables.NAMED_KINDS, once)
			appear := k.NewDomainName(domain, "appearance", tables.NAMED_FIELD, once)
			k.NewField(things, appear, tables.PRIM_TEXT, "")
		}
		prop := k.NewName("appearance", tables.NAMED_FIELD, op.At.String())
		noun := LastNameOf(k.Recent.Nouns.Subjects)
		k.NewValue(noun, prop, text)
	}
	return
}

func (op *Tail) Import(k *Importer) (err error) {
	if e := op.Pronoun.Import(k); e != nil {
		err = e
	} else if e := op.NounPhrase.Import(k); e != nil {
		err = e
	}
	return
}

func (op *TraitPhrase) ImportTraits(k *Importer, aspect ephemera.Named) (err error) {
	for rank, trait := range op.Trait {
		if t, e := trait.NewName(k); e != nil {
			err = errutil.Append(err, e)
		} else {
			k.NewTrait(t, aspect, rank)
		}
	}
	return
}
