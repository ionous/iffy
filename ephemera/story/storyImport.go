package story

import (
	"log"

	"github.com/ionous/errutil"
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

// (the) colors are red, blue, or green.
func (op *KindOfRelation) ImportPhrase(k *Importer) (err error) {
	log.Println("KindOfRelation not implemented")
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
