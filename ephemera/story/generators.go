package story

import (
	"unicode"
	"unicode/utf8"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/tables"
)

// FIX -- keeping with the function parser model, swap these out to functions.
// the top level function should match the top level of "story"
// which a switch for parsing "story_statement" slats.

// story is a bunch of paragraphs
//make.run("story", "{+paragraph}");
func imp_story(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "story"); e != nil {
		err = e
	} else {
		err = reader.Repeats(m.SliceOf("$PARAGRAPH"), k.Bind(imp_paragraph))
	}
	return
}

// paragraph is a bunch of statements on the same line
func imp_paragraph(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "paragraph"); e != nil {
		err = e
	} else {
		// k.StoryEnv = StoryEnv{}
		err = reader.Repeats(m.SliceOf("$STORY_STATEMENT"), k.Bind(imp_story_statement))
	}
	return
}

// run: "{+noun} {noun_phrase}."
func imp_lede(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "lede"); e != nil {
		err = e
	} else if e := k.Recent.Nouns.CollectSubjects(func() error {
		return reader.Repeats(m.SliceOf("$NOUNS"), k.Bind(imp_named_noun))
	}); e != nil {
		err = e
	} else if e := imp_noun_phrase(k, m.MapOf("$NOUN_PHRASE")); e != nil {
		err = e
	}
	return
}

// run: "{pronoun} {noun_phrase}."
func imp_tail(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "tail"); e != nil {
		err = e
	} else if e := imp_pronoun(k, m.MapOf("$PRONOUN")); e != nil {
		err = e
	} else if e := imp_noun_phrase(k, m.MapOf("$NOUN_PHRASE")); e != nil {
		err = e
	}
	return
}

// make.str("pronoun",  "{it}, {they}, or {pronoun}");
func imp_pronoun(k *Importer, r reader.Map) (err error) {
	if _, e := reader.String(r, "pronoun"); e != nil {
		err = e
	} else {
		// FIX: this can indicate plurality
	}
	return
}

// opt: "{kind_of_noun}, {noun_traits}, or {noun_relation}"
func imp_noun_phrase(k *Importer, r reader.Map) (err error) {
	return reader.Option(r, "noun_phrase", reader.ReadMaps{
		"$KIND_OF_NOUN":  k.Bind(imp_kind_of_noun),
		"$NOUN_TRAITS":   k.Bind(imp_noun_traits),
		"$NOUN_RELATION": k.Bind(imp_noun_relation),
	})
}

// run: "{?are_being} {relation} {+noun}"
// ex. [the cat and the hat] (are) (in) (the book)
// ex. [Hector and Maria] (are) (suspicious of) (Santa and Santana).
// fix? parse are_being
func imp_noun_relation(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "noun_relation"); e != nil {
		err = e
	} else if relation, e := imp_relation(k, m.MapOf("$RELATION")); e != nil {
		err = e
	} else if e := k.Recent.Nouns.CollectObjects(func() error {
		return reader.Repeats(m.SliceOf("$NOUNS"), k.Bind(imp_named_noun))
	}); e != nil {
		err = e
	} else {
		for _, subject := range k.Recent.Nouns.Subjects {
			for _, object := range k.Recent.Nouns.Objects {
				k.NewRelative(subject, relation, object)
			}
		}
	}
	return
}

// run: "{determiner} {noun_name}"
func imp_named_noun(k *Importer, r reader.Map) (err error) {
	// declare a noun class that has several default fields
	if once := "noun"; k.Once(once) {
		// common or proper nouns ( rabbit, vs. Roger )
		k.NewImplicitAspect("nounTypes", "things", "commonNamed", "properNamed", "counted")
		// whether a player can refer to an object by its name.
		k.NewImplicitAspect("privateNames", "things", "publiclyNamed", "privatelyNamed")
	}
	//
	if m, e := reader.Unpack(r, "named_noun"); e != nil {
		err = e
	} else if det, e := imp_determiner(k, m.MapOf("$DETERMINER")); e != nil {
		err = e
	} else {
		name := m.MapOf("$NAME")
		if cnt, ok := lang.WordsToNum(det); !ok {
			err = read_named_noun(k, det, name)
		} else {
			err = read_counted_noun(k, cnt, name)
		}
	}
	return
}

func read_counted_noun(k *Importer, cnt int, r reader.Map) (err error) {
	kind := "singular_kind"
	if cnt > 1 {
		kind = "plural_kinds"
	}
	if countedKind, e := importName(k, r, "noun_name", kind); e != nil {
		err = e
	} else {
		typeTrait := k.NewName("counted", tables.NAMED_TRAIT, reader.At(r))
		nameTrait := k.NewName("privatelyNamed", tables.NAMED_TRAIT, reader.At(r))
		// this isnt right.
		// even with an in memory map its not quite right because technically ephemera can come from multiple sources
		// fix: something something noun stacks, not individually duplicated nouns
		baseName := countedKind.String()
		for i := 0; i < cnt; i++ {
			countedNoun := k.autoCounter.next(baseName)
			noun := k.NewName(countedNoun, "noun", reader.At(r))
			k.Recent.Nouns.Add(noun)
			k.NewValue(noun, nameTrait, true)
			k.NewValue(noun, typeTrait, true)
		}
	}
	return
}

func read_named_noun(k *Importer, det string, r reader.Map) (err error) {
	if noun, e := imp_noun_name(k, r); e != nil {
		err = e
	} else {
		k.Recent.Nouns.Add(noun)
		// pick common or proper based on noun capitalization.
		// fix: implicitly generated facts should be considered preliminary
		// so that authors can override them.
		traitStr := "commonNamed"
		if first, _ := utf8.DecodeRuneInString(noun.String()); unicode.ToUpper(first) == first {
			traitStr = "properNamed"
		}
		typeTrait := k.NewName(traitStr, tables.NAMED_TRAIT, reader.At(r))
		k.NewValue(noun, typeTrait, true)

		// record any custom determiner
		if usesKeyWord := det[0] == '$'; !usesKeyWord {
			// set the indefinite article field
			article := k.NewName("indefinite article", tables.NAMED_FIELD, reader.At(r))
			k.NewValue(noun, article, det)

			// create a "indefinite article" field for all "things"
			if once := "named_noun"; k.Once(once) {
				domain := k.gameDomain()
				things := k.NewDomainName(domain, "things", tables.NAMED_KINDS, once)
				indefinite := k.NewDomainName(domain, "indefinite article", tables.NAMED_FIELD, once)
				k.NewField(things, indefinite, tables.PRIM_TEXT)
			}
		}
	}
	return
}

// set proper nounType to true ( implicitly defined by "noun" )
func imp_determiner(k *Importer, r reader.Map) (ret string, err error) {
	return reader.String(r, "determiner")
}

// run: "{are_an} {*trait:*trait} {kind:singular_kind} {?noun_relation}"
// ex. "[the box] (is a) (closed) (container) ((on) (the beach))"
func imp_kind_of_noun(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "kind_of_noun"); e != nil {
		err = e
	} else if kind, e := imp_singular_kind(k, m.MapOf("$KIND")); e != nil {
		err = e
	} else {
		var traits []ephemera.Named
		if e := reader.Repeats(m.SliceOf("$TRAIT"), func(el reader.Map) (err error) {
			if trait, e := imp_trait(k, el); e != nil {
				err = e
			} else {
				traits = append(traits, trait)
			}
			return
		}); e != nil {
			err = e
		} else {
			// we collected the nouns and delayed processing them till now.
			for _, noun := range k.Recent.Nouns.Subjects {
				k.NewNoun(noun, kind)
				for _, trait := range traits {
					k.NewValue(noun, trait, true) // the value of the trait for the noun is true
				}
			}
			if v := m.MapOf("$NOUN_RELATION"); len(v) != 0 {
				err = imp_noun_relation(k, v)
			}
		}
	}
	return
}

// run: "{The [summary] is:: %lines}"
func imp_summary(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "summary"); e != nil {
		err = e
	} else if lines, e := imp_lines(k, m.MapOf("$LINES")); e != nil {
		err = e
	} else if text, e := convert_text_or_template(lines); e != nil {
		err = e
	} else {
		// declare the existence of the field "appearance"
		if once := "summary"; k.Once(once) {
			domain := k.gameDomain()
			things := k.NewDomainName(domain, "things", tables.NAMED_KINDS, once)
			appear := k.NewDomainName(domain, "appearance", tables.NAMED_FIELD, once)
			k.NewField(things, appear, tables.PRIM_TEXT)
		}
		prop := k.NewName("appearance", tables.NAMED_FIELD, reader.At(m))
		noun := LastNameOf(k.Recent.Nouns.Subjects)
		k.NewValue(noun, prop, text)
	}
	return
}

// run: "{are_being} {+trait:trait}"
// ex. "(the box) is closed"
func imp_noun_traits(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "noun_traits"); e != nil {
		err = e
	} else {
		err = reader.Repeats(m.SliceOf("$TRAIT"), func(el reader.Map) (err error) {
			if trait, e := imp_trait(k, el); e != nil {
				err = e
			} else {
				for _, noun := range k.Recent.Nouns.Subjects {
					k.NewValue(noun, trait, true) // the value of the trait for the noun is true
				}
			}
			return
		})
	}
	return
}

// fix... part of class traits
func imp_trait_phrase(k *Importer, r reader.Map) (ret []ephemera.Named, err error) {
	if m, e := reader.Unpack(r, "trait_phrase"); e != nil {
		err = e
	} else {
		err = reader.Repeats(m.SliceOf("$TRAIT"), func(el reader.Map) (err error) {
			if trait, e := imp_trait(k, el); e != nil {
				err = e
			} else {
				ret = append(ret, trait)
			}
			return
		})
	}
	return
}

type variableDecl struct {
	name, typeName ephemera.Named
	affinity       string
}

// "{type:variable_type} ( called {name:variable_name|quote} )"
func imp_variable_decl(k *Importer, cat string, r reader.Map) (ret variableDecl, err error) {
	if m, e := reader.Unpack(r, "variable_decl"); e != nil {
		err = e
	} else if n, e := imp_variable_name(k, cat, m.MapOf("$NAME")); e != nil {
		err = e
	} else if t, aff, e := imp_variable_type(k, m.MapOf("$TYPE")); e != nil {
		err = e
	} else {
		ret = variableDecl{n, t, aff}
	}
	return
}

func imp_variable_type(k *Importer, r reader.Map) (retType ephemera.Named, retAffinity string, err error) {
	err = reader.Option(r, "variable_type", reader.ReadMaps{
		"$PRIMITIVE": func(m reader.Map) (err error) {
			retType, err = imp_primitive_var(k, m)
			return
		},
		"$OBJECT": func(m reader.Map) (err error) {
			retType, err = imp_object_type(k, m)
			retAffinity = affine.Object.String()
			return
		},
		"$EXT": func(m reader.Map) (err error) {
			err = errutil.New("extension types not implemented")
			return
		},
	})
	return
}

// "{a number%number}, {some text%text}, or {a true/false value%bool}"
// returns one of the evalType(s) as a "Named" value --
// we return a name to normalize references to object kinds which are also used as variables
func imp_primitive_var(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	if evalType, e := reader.Enum(r, "primitive_type", reader.Map{
		"$NUMBER": "number_eval",
		"$TEXT":   "text_eval",
		"$BOOL":   "bool_eval",
	}); e != nil {
		err = e
	} else {
		ret = k.NewName(evalType.(string), tables.NAMED_TYPE, reader.At(r))
	}
	return
}

// ick. fix. see imp_primitive_phrase.
func imp_primitive_prop(k *Importer, r reader.Map) (string, error) {
	p, e := reader.Enum(r, "primitive_type", reader.Map{
		"$NUMBER": tables.PRIM_DIGI,
		"$TEXT":   tables.PRIM_TEXT,
		"$BOOL":   tables.PRIM_BOOL,
	})
	return p.(string), e
}

// "{an} {kind of%kind:singular_kind} object"
// returns the name of "singular_kind"
func imp_object_type(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	if m, e := reader.Unpack(r, "object_type"); e != nil {
		err = e
	} else {
		ret, err = imp_singular_kind(k, m.MapOf("$KIND"))
	}
	return
}

func Unimplemented(k *Importer, r reader.Map) (err error) {
	return errutil.New("unimplemented", r.StrOf(reader.ItemType), reader.At(r))
}
