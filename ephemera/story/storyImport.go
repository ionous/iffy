package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

type GenericImport interface {
	Import(*Importer) error
}

type StoryStatement interface {
	ImportPhrase(k *Importer) error
}

// story is a bunch of paragraphs
func (op *Story) ImportStory(k *Importer) (err error) {
	if els := op.Paragraph; els != nil {
		for _, el := range *els {
			if e := el.ImportParagraph(k); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

// paragraph is a bunch of statements on the same line
func (op *Paragraph) ImportParagraph(k *Importer) (err error) {
	if els := op.StoryStatement; els != nil {
		for _, el := range *els {
			if e := el.ImportPhrase(k); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

// (the) colors are red, blue, or green.
func (op *AspectTraits) ImportPhrase(k *Importer) (err error) {
	if aspect, e := op.Aspect.NewName(k); e != nil {
		err = e
	} else {
		err = op.TraitPhrase.ImportTraits(k, aspect)
	}
	return
}

// horses are usually fast.
func (op *Certainties) ImportPhrase(k *Importer) (err error) {
	if certainty, e := op.Certainty.ImportString(k); e != nil {
		err = e
	} else if trait, e := op.Trait.NewName(k); e != nil {
		err = e
	} else if kind, e := op.PluralKinds.NewName(k); e != nil {
		err = e
	} else {
		k.NewCertainty(certainty, trait, kind)
	}
	return
}

func (op *Comment) ImportPhrase(k *Importer) (err error) {
	// do nothing for now.
	return
}

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
		err = op.PropertyPhrase.ImportProperties(k, kind)
	}
	// fix: handle determiner?
	return
}

// ex. The description of the nets is xxx
func (op *NounAssignment) ImportPhrase(k *Importer) (err error) {
	if prop, e := op.Property.NewName(k); e != nil {
		err = e
	} else if text, e := op.Lines.ConvertText(); e != nil {
		err = e
	} else if e := k.Recent.Nouns.CollectSubjects(func() (err error) {
		for _, n := range op.Nouns {
			if e := n.Import(k); e != nil {
				err = errutil.Append(err, e)
			}
		}
		return
	}); e != nil {
		err = e
	} else {
		for _, noun := range k.Recent.Nouns.Subjects {
			k.NewValue(noun, prop, text)
		}
	}
	return
}
func (op *NounStatement) ImportPhrase(k *Importer) (err error) {
	if e := op.Lede.Import(k); e != nil {
		err = e
	} else {
		if els := op.Tail; els != nil {
			for _, el := range *els {
				if e := el.Import(k); e != nil {
					err = errutil.Append(err, e)
				}
			}
		}
		if err == nil && op.Summary != nil {
			err = op.Summary.Import(k)
		}
	}
	return
}
func (op *PatternActions) ImportPhrase(k *Importer) (err error) {
	if patternName, e := op.Name.NewName(k); e != nil {
		err = e
	} else if e := op.PatternRules.ImportPattern(k, patternName); e != nil {
		err = e
	} else {
		if els := op.PatternLocals; els != nil {
			err = els.ImportPattern(k, patternName)
		}
	}
	return
}

// Adds a new pattern declaration and optionally some associated pattern parameters.
func (op *PatternDecl) ImportPhrase(k *Importer) (err error) {
	if patternName, e := op.Name.NewName(k); e != nil {
		err = e
	} else if patternType, e := op.Type.ImportType(k); e != nil {
		err = e
	} else {
		k.NewPatternDecl(patternName, patternName, patternType, "", ephemera.Prog{})
		//
		if els := op.Optvars; els != nil {
			for _, el := range els.VariableDecl {
				if val, e := el.ImportVariable(k, tables.NAMED_PARAMETER); e != nil {
					err = errutil.Append(err, e)
				} else {
					k.NewPatternDecl(patternName, val.name, val.typeName, val.affinity, ephemera.Prog{})
				}
			}
		}
	}
	return
}
func (op *PatternVariablesDecl) ImportPhrase(k *Importer) (err error) {
	if patternName, e := op.PatternName.NewName(k); e != nil {
		err = e
	} else {
		// fix: shouldnt this be called pattern parameters?
		for _, el := range op.VariableDecl {
			if val, e := el.ImportVariable(k, tables.NAMED_PARAMETER); e != nil {
				err = errutil.Append(err, e)
			} else {
				k.NewPatternDecl(patternName, val.name, val.typeName, val.affinity, ephemera.Prog{})
			}
		}
	}
	return
}

// ex. On the beach are shells.
func (op *RelativeToNoun) ImportPhrase(k *Importer) (err error) {
	if relation, e := op.Relation.NewName(k); e != nil {
		err = e
	} else if e := k.Recent.Nouns.CollectObjects(func() error {
		return ImportNamedNouns(k, op.Nouns)
	}); e != nil {
		err = e
	} else if e := k.Recent.Nouns.CollectSubjects(func() error {
		return ImportNamedNouns(k, op.Nouns1)
	}); e != nil {
		err = e
	} else {
		for _, object := range k.Recent.Nouns.Objects {
			for _, subject := range k.Recent.Nouns.Subjects {
				k.NewRelative(object, relation, subject)
			}
		}
	}
	return
}
