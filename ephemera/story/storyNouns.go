package story

import (
	"unicode"
	"unicode/utf8"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/tables"
)

func (*Pronoun) Import(*Importer) (err error) {
	// FIX: pronoun(s) can indicate plurality
	return
}

func (op *NounPhrase) Import(k *Importer) (err error) {
	if imp, ok := op.Opt.(GenericImport); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	} else {
		err = imp.Import(k)
	}
	return
}

func ImportNamedNouns(k *Importer, els []NamedNoun) (err error) {
	for _, el := range els {
		if e := el.Import(k); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (op *NamedNoun) Import(k *Importer) (err error) {
	// declare a noun class that has several default fields
	if once := "noun"; k.Once(once) {
		// common or proper nouns ( rabbit, vs. Roger )
		k.NewImplicitAspect("noun_types", "things", "common_named", "proper_named", "counted")
		// whether a player can refer to an object by its name.
		k.NewImplicitAspect("private_names", "things", "publicly_named", "privately_named")
	}
	//
	if cnt, ok := lang.WordsToNum(op.Determiner.Str); !ok {
		err = op.ReadNamedNoun(k)
	} else {
		err = op.ReadCountedNoun(k, cnt)
	}
	return
}

func (op *NamedNoun) ReadCountedNoun(k *Importer, cnt int) (err error) {
	// probably? should have a specific counted noun phrase b/c
	// and two things are things doesnt make much sense
	//
	// cat := tables.NAMED_KIND
	// if cnt > 1 {
	// 	cat = tables.NAMED_PLURAL_KINDS
	// }
	// op.Name.AddNameWithCategory(k, cat)
	typeTrait := k.NewName("counted", tables.NAMED_TRAIT, op.At.String())
	// fix: something something noun stacks, not individually duplicated nouns
	baseName := op.Name.String()
	for i := 0; i < cnt; i++ {
		countedNoun := k.autoCounter.Next(baseName)
		noun := k.NewName(countedNoun, "noun", op.At.String())
		k.Recent.Nouns.Add(noun)
		k.NewValue(noun, typeTrait, true)
	}
	return
}

func (op *NamedNoun) ReadNamedNoun(k *Importer) (err error) {
	if noun, e := op.Name.NewName(k); e != nil {
		err = e
	} else {
		k.Recent.Nouns.Add(noun)
		// pick common or proper based on noun capitalization.
		// fix: implicitly generated facts should be considered preliminary
		// so that authors can override them.
		traitStr := "common_named"
		if first, _ := utf8.DecodeRuneInString(noun.String()); unicode.ToUpper(first) == first {
			traitStr = "proper_named"
		}
		typeTrait := k.NewName(traitStr, tables.NAMED_TRAIT, op.At.String())
		k.NewValue(noun, typeTrait, true)

		// record any custom determiner
		if str, ok := decode.FindChoice(&op.Determiner, op.Determiner.Str); ok && len(str) == 0 {
			// set the indefinite article field
			article := k.NewName("indefinite_article", tables.NAMED_FIELD, op.At.String())
			k.NewValue(noun, article, op.Determiner.Str)

			// create a "indefinite article" field for all "things"
			if once := "named_noun"; k.Once(once) {
				domain := k.gameDomain()
				things := k.NewDomainName(domain, "things", tables.NAMED_KINDS, once)
				indefinite := k.NewDomainName(domain, "indefinite_article", tables.NAMED_FIELD, once)
				k.NewField(things, indefinite, tables.PRIM_TEXT, "")
			}

		}
	}
	return
}

// ex. "[the box] (is a) (closed) (container) ((on) (the beach))"
func (op *KindOfNoun) Import(k *Importer) (err error) {
	if kind, e := op.Kind.NewName(k); e != nil {
		err = e
	} else {
		//
		var traits []ephemera.Named
		if ts := op.Trait; ts != nil {
			for _, t := range *ts {
				if t, e := t.NewName(k); e != nil {
					err = errutil.Append(err, e)
				} else {
					traits = append(traits, t)
				}
			}
		}
		if err == nil {
			// we collected the nouns and delayed processing them till now.
			for _, noun := range k.Recent.Nouns.Subjects {
				k.NewNoun(noun, kind)
				for _, trait := range traits {
					k.NewValue(noun, trait, true) // the value of the trait for the noun is true
				}
			}
		}
	}
	return
}

// ex. [the cat and the hat] (are) (in) (the book)
// ex. [Hector and Maria] (are) (suspicious of) (Santa and Santana).
func (op *NounRelation) Import(k *Importer) (err error) {
	if rel, e := op.Relation.NewName(k); e != nil {
		err = e
	} else if e := k.Recent.Nouns.CollectObjects(func() (err error) {
		return ImportNamedNouns(k, op.Nouns)
	}); e != nil {
		err = e
	} else {
		for _, subject := range k.Recent.Nouns.Subjects {
			for _, object := range k.Recent.Nouns.Objects {
				k.NewRelative(subject, rel, object)
			}
		}
	}
	return
}

//
func (op *NounTraits) Import(k *Importer) (err error) {
	for _, t := range op.Trait {
		if trait, e := t.NewName(k); e != nil {
			err = e
			break
		} else {
			for _, noun := range k.Recent.Nouns.Subjects {
				k.NewValue(noun, trait, true) // the value of the trait for the noun is true
			}
		}
	}
	return
}
