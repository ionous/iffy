package story

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/tables"
)

func (op *Aspect) NewName(k *Importer) (ret ephemera.Named, err error) {
	name := lang.Breakcase(op.Str)
	return k.NewName(name, tables.NAMED_ASPECT, op.At.String()), nil
}

func (op *NounName) NewName(k *Importer) (ret ephemera.Named, err error) {
	name := lang.Breakcase(op.Str)
	return k.NewName(name, tables.NAMED_NOUN, op.At.String()), nil
}

func (op *NounName) AddNameWithCategory(k *Importer, cat string) (ret ephemera.Named, err error) {
	name := lang.Breakcase(op.Str)
	return k.NewName(name, cat, op.At.String()), nil
}

func (op *PatternName) NewName(k *Importer) (ret ephemera.Named, err error) {
	name := lang.Breakcase(op.Str)
	return k.NewName(name, tables.NAMED_PATTERN, op.At.String()), nil
}

func (op *PluralKinds) NewName(k *Importer) (ret ephemera.Named, err error) {
	name := lang.Breakcase(op.Str)
	return k.NewName(name, tables.NAMED_PLURAL_KINDS, op.At.String()), nil
}

func (op *Property) NewName(k *Importer) (ret ephemera.Named, err error) {
	// note: this is linked to NAMED_ASPECT
	// aspect properties in kinds currently must have the same name as the aspect.
	name := lang.Breakcase(op.Str)
	return k.NewName(name, tables.NAMED_FIELD, op.At.String()), nil
}

func (op *RecordSingular) NewName(k *Importer) (ret ephemera.Named, err error) {
	// fix? for now, we leverage the existing kind assembly
	name := lang.LowerBreakcase(op.Str)
	return k.NewName(name, tables.NAMED_KIND, op.At.String()), nil
}

func (op *RecordPlural) NewName(k *Importer) (ret ephemera.Named, err error) {
	// fix? for now, we leverage the existing kind assembly
	name := lang.Breakcase(op.Str)
	return k.NewName(name, tables.NAMED_PLURAL_KINDS, op.At.String()), nil
}

func (op *RelationName) NewName(k *Importer) (ret ephemera.Named, err error) {
	name := lang.Breakcase(op.Str)
	return k.NewName(name, tables.NAMED_VERB, op.At.String()), nil
}

func (op *SingularKind) NewName(k *Importer) (ret ephemera.Named, err error) {
	name := lang.Breakcase(op.Str)
	return k.NewName(name, tables.NAMED_KIND, op.At.String()), nil
}

func (op *TestName) NewName(k *Importer) (ret ephemera.Named, err error) {
	// fix? all names should probably munge their own strings
	// ( see ephemera's NewDomainName for the current hack )
	// things that would need work are:
	// tests, autogen fields, and control over the domain ( ie. NewName vs. NewDowmainName )
	if op.Str == "$CURRENT_TEST" {
		ret = k.StoryEnv.Recent.Test
	} else {
		name := lang.Breakcase(op.Str)
		ret = k.NewName(name, tables.NAMED_TEST, op.At.String())
	}
	return
}

func (op *Trait) NewName(k *Importer) (ret ephemera.Named, err error) {
	name := lang.Breakcase(op.Str)
	return k.NewName(name, tables.NAMED_TRAIT, op.At.String()), nil
}

func (op *VariableName) NewName(k *Importer, cat string) (ret ephemera.Named, err error) {
	name := lang.Breakcase(op.Str)
	return k.NewName(name, cat, op.At.String()), nil
}
