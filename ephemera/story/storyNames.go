package story

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

func (op *Aspect) NewName(k *Importer) (ret ephemera.Named, err error) {
	return k.NewName(op.Str, tables.NAMED_ASPECT, op.At.String()), nil
}

func (op *NounName) NewName(k *Importer) (ret ephemera.Named, err error) {
	return k.NewName(op.Str, tables.NAMED_NOUN, op.At.String()), nil
}

func (op *NounName) AddNameWithCategory(k *Importer, cat string) (ret ephemera.Named, err error) {
	return k.NewName(op.Str, cat, op.At.String()), nil
}

func (op *PatternName) NewName(k *Importer) (ret ephemera.Named, err error) {
	return k.NewName(op.Str, tables.NAMED_PATTERN, op.At.String()), nil
}

func (op *PluralKinds) NewName(k *Importer) (ret ephemera.Named, err error) {
	return k.NewName(op.Str, tables.NAMED_PLURAL_KINDS, op.At.String()), nil
}

func (op *Property) NewName(k *Importer) (ret ephemera.Named, err error) {
	return k.NewName(op.Str, tables.NAMED_FIELD, op.At.String()), nil
}

func (op *Relation) NewName(k *Importer) (ret ephemera.Named, err error) {
	return k.NewName(op.Str, tables.NAMED_VERB, op.At.String()), nil
}

func (op *SingularKind) NewName(k *Importer) (ret ephemera.Named, err error) {
	return k.NewName(op.Str, tables.NAMED_KIND, op.At.String()), nil
}

func (op *TestName) NewName(k *Importer) (ret ephemera.Named, err error) {
	k.OverrideNameDuring("$CURRENT_TEST", k.StoryEnv.Recent.Test, func() {
		ret = k.NewName(op.Str, tables.NAMED_TEST, op.At.String())
	})
	return
}

func (op *Trait) NewName(k *Importer) (ret ephemera.Named, err error) {
	return k.NewName(op.Str, tables.NAMED_TRAIT, op.At.String()), nil
}

func (op *VariableName) NewName(k *Importer, cat string) (ret ephemera.Named, err error) {
	return k.NewName(op.Str, cat, op.At.String()), nil
}
