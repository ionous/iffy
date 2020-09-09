package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
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
		return reader.Repeats(m.SliceOf("$NOUN"), k.Bind(imp_noun))
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
		return reader.Repeats(m.SliceOf("$NOUN"), k.Bind(imp_noun))
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

// noun: "{proper_noun} or {common_noun}"
func imp_noun(k *Importer, r reader.Map) (err error) {
	// declare a noun class that has several default fields
	if once := "noun"; k.Once(once) {
		// common or proper nouns ( rabbit, vs. Roger )
		k.NewImplicitAspect("nounTypes", "things", "common", "proper")
		// whether a player can refer to an object by its name.
		k.NewImplicitAspect("privateNames", "things", "publiclyNamed", "privatelyNamed")
	}
	return reader.Option(r, "noun", reader.ReadMaps{
		"$PROPER_NOUN": k.Bind(imp_proper_noun),
		"$COMMON_NOUN": k.Bind(imp_common_noun),
	})
}

// run: "{determiner} {common_name}"
func imp_common_noun(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "common_noun"); e != nil {
		err = e
	} else if det, e := imp_determiner(k, m.MapOf("$DETERMINER")); e != nil {
		err = e
	} else if noun, e := imp_common_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else {
		k.Recent.Nouns.Add(noun)

		// set common nounType to true ( implicitly defined by "noun" )
		nounType := k.NewName("common", tables.NAMED_TRAIT, reader.At(r))
		k.NewValue(noun, nounType, true)
		//
		if det[0] != '$' {
			article := k.NewName("indefinite article", tables.NAMED_FIELD, reader.At(r))
			k.NewValue(noun, article, det)
			if once := "common_noun"; k.Once(once) {
				domain := k.gameDomain()
				indefinite := k.NewDomainName(domain, "indefinite article", tables.NAMED_FIELD, once)
				things := k.NewDomainName(domain, "things", tables.NAMED_KINDS, once)
				k.NewField(things, indefinite, tables.PRIM_TEXT)
			}
		}
	}
	return
}

func imp_determiner(k *Importer, r reader.Map) (ret string, err error) {
	return reader.String(r, "determiner")
}

// run: "{proper_name}"
// common / proper setting
func imp_proper_noun(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "proper_noun"); e != nil {
		err = e
	} else if noun, e := imp_proper_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else {
		k.Recent.Nouns.Add(noun)
		// set proper nounType to true ( implicitly defined by "noun" )
		nounType := k.NewName("proper", tables.NAMED_TRAIT, reader.At(m))
		k.NewValue(noun, nounType, true)
	}
	return
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
	} else if lines, e := imp_line_expr(k, m.MapOf("$LINES")); e != nil {
		err = e
	} else {
		// declare the existence of the field "appearance"
		if once := "summary"; k.Once(once) {
			domain := k.gameDomain()
			things := k.NewDomainName(domain, "things", tables.NAMED_KINDS, once)
			appear := k.NewDomainName(domain, "appearance", tables.NAMED_FIELD, once)
			k.NewField(things, appear, tables.PRIM_EXPR)
		}
		prop := k.NewName("appearance", tables.NAMED_FIELD, reader.At(m))
		noun := LastNameOf(k.Recent.Nouns.Subjects)
		k.NewValue(noun, prop, lines)
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
	err = reader.Repeats(r.SliceOf("$TRAIT"), func(el reader.Map) (err error) {
		if trait, e := imp_trait(k, el); e != nil {
			err = e
		} else {
			ret = append(ret, trait)
		}
		return
	})
	return
}

// "{type:variable_type} ( called {name:variable_name|quote} )"
func imp_variable_decl(k *Importer, r reader.Map) (retName, retType ephemera.Named, err error) {
	if m, e := reader.Unpack(r, "variable_decl"); e != nil {
		err = e
	} else if n, e := imp_variable_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else if t, e := imp_variable_type(k, m.MapOf("$TYPE")); e != nil {
		err = e
	} else {
		retName, retType = n, t
	}
	return
}

func imp_variable_type(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	err = reader.Option(r, "variable_type", reader.ReadMaps{
		"$PRIMITIVE": func(m reader.Map) (err error) {
			ret, err = imp_primitive_var(k, m)
			return
		},
		"$OBJECT": func(m reader.Map) (err error) {
			ret, err = imp_object_type(k, m)
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
