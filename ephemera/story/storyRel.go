package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

func (op *KindOfRelation) ImportPhrase(k *Importer) (err error) {
	// rec.NewRelation(r, k, q, c)
	if rel, e := op.Relation.NewName(k); e != nil {
		err = e
	} else if card, e := op.RelationCardinality.ImportCardinality(k); e != nil {
		err = e
	} else {
		k.NewRelation(rel, card.firstKind, card.secondKind, card.cardinality)
	}
	return
}

type importedCardinality struct {
	cardinality           string // tables.ONE_TO_ONE
	firstKind, secondKind ephemera.Named
}

func (op *RelationCardinality) ImportCardinality(k *Importer) (ret importedCardinality, err error) {
	type cardinalityImporter interface {
		ImportCardinality(k *Importer) (importedCardinality, error)
	}
	if c, ok := op.Opt.(cardinalityImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	} else {
		ret, err = c.ImportCardinality(k)
	}
	return
}

func (op *OneToOne) ImportCardinality(k *Importer) (ret importedCardinality, err error) {
	if first, e := op.Kind.FixPlurals(k); e != nil {
		err = e
	} else if second, e := op.OtherKind.FixPlurals(k); e != nil {
		err = e
	} else {
		ret = importedCardinality{tables.ONE_TO_ONE, first, second}
	}
	return
}
func (op *OneToMany) ImportCardinality(k *Importer) (ret importedCardinality, err error) {
	if first, e := op.Kind.FixPlurals(k); e != nil {
		err = e
	} else if second, e := op.Kinds.NewName(k); e != nil {
		err = e
	} else {
		ret = importedCardinality{tables.ONE_TO_MANY, first, second}
	}
	return
}
func (op *ManyToOne) ImportCardinality(k *Importer) (ret importedCardinality, err error) {
	if first, e := op.Kinds.NewName(k); e != nil {
		err = e
	} else if second, e := op.Kind.FixPlurals(k); e != nil {
		err = e
	} else {
		ret = importedCardinality{tables.MANY_TO_ONE, first, second}
	}
	return
}
func (op *ManyToMany) ImportCardinality(k *Importer) (ret importedCardinality, err error) {
	if first, e := op.Kinds.NewName(k); e != nil {
		err = e
	} else if second, e := op.OtherKinds.NewName(k); e != nil {
		err = e
	} else {
		ret = importedCardinality{tables.MANY_TO_MANY, first, second}
	}
	return
}
