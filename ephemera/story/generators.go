package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

// FIX -- keeping with the function parser model, swap these out to functions.
// the top level function should match the top level of "story"
// which a switch for parsing "story_statement" slats.

// story is a bunch of paragraphs
//make.run("story", "{+paragraph|ghost}");
func imp_story(k *imp.Porter, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "story"); e != nil {
		err = e
	} else {
		err = reader.Repeats(m.SliceOf("$PARAGRAPH"), k.Bind(imp_paragraph))
	}
	return
}

// paragraph is a bunch of statements
func imp_paragraph(k *imp.Porter, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "paragraph"); e != nil {
		err = e
	} else {
		storyNouns.Swap(nil)
		err = reader.Repeats(m.SliceOf("$STORY_STATEMENT"), k.Bind(imp_story_statement))
	}
	return
}

// run: "{+names} {noun_phrase}."
func imp_lede(k *imp.Porter, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "lede"); e != nil {
		err = e
	} else if e := reader.Repeats(m.SliceOf("$NOUN"), k.Bind(imp_noun)); e != nil {
		err = e
	} else if e := imp_noun_phrase(k, m.MapOf("$NOUN_PHRASE")); e != nil {
		err = e
	}
	return
}

// run: "{pronoun} {noun_phrase}."
func imp_tail(k *imp.Porter, r reader.Map) (err error) {
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
func imp_pronoun(k *imp.Porter, r reader.Map) (err error) {
	if _, e := reader.String(r, "pronoun"); e != nil {
		err = e
	} else {
		// FIX: this can indicate plurality
	}
	return
}

// opt: "{kind_of_noun}, {noun_attrs}, or {noun_relation}"
func imp_noun_phrase(k *imp.Porter, r reader.Map) (err error) {
	return reader.Option(r, "noun_phrase", reader.ReadMaps{
		"$KIND_OF_NOUN":  k.Bind(imp_kind_of_noun),
		"$NOUN_ATTRS":    k.Bind(imp_noun_attrs),
		"$NOUN_RELATION": k.Bind(imp_noun_relation),
	})
}

// run: "{?are_being} {relation} {+noun}"
// ex. (the cat and the hat) are in (the book)
// ex. (Hector and Maria) are suspicious of (Santa and Santana).
func imp_noun_relation(k *imp.Porter, r reader.Map) (err error) {
	// unexpected type  wanted noun_relation at
	if m, e := reader.Unpack(r, "noun_relation"); e != nil {
		err = e
	} else if relation, e := imp_relation(k, m.MapOf("$RELATION")); e != nil {
		err = e
	} else {
		// fix? parse are_being
		leadinstoryNouns := storyNouns.Swap(nil)
		if e := reader.Repeats(m.SliceOf("$NOUN"), k.Bind(imp_noun)); e != nil {
			err = e
		} else {
			trailinstoryNouns := storyNouns.Swap(leadinstoryNouns)
			//
			for _, n := range leadinstoryNouns {
				for _, d := range trailinstoryNouns {
					k.NewRelative(n, relation, d)
				}
			}
		}
	}
	return
}

// noun: "{proper_noun} or {common_noun}"
func imp_noun(k *imp.Porter, r reader.Map) (err error) {
	// declare a noun class that has several default fields
	if once := "noun"; k.Once(once) {
		things := k.NewName(tables.NAMED_KIND, "things", once)
		nounType := k.NewName(tables.NAMED_ASPECT, "nounType", once)
		common := k.NewName(tables.NAMED_TRAIT, "common", once)
		proper := k.NewName(tables.NAMED_TRAIT, "proper", once)
		k.NewPrimitive(tables.PRIM_ASPECT, things, nounType)
		k.NewTrait(common, nounType, 0)
		k.NewTrait(proper, nounType, 1)
	}
	return reader.Option(r, "noun", reader.ReadMaps{
		"$PROPER_NOUN": k.Bind(imp_proper_noun),
		"$COMMON_NOUN": k.Bind(imp_common_noun),
	})
}

// run: "{determiner} {common_name}"
func imp_common_noun(k *imp.Porter, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "common_noun"); e != nil {
		err = e
	} else if det, e := imp_determiner(k, m.MapOf("$DETERMINER")); e != nil {
		err = e
	} else if noun, e := imp_common_name(k, m.MapOf("$COMMON_NAME")); e != nil {
		err = e
	} else {
		storyNouns.Add(noun)
		// set common nounType to true ( implicitly defined by "noun" )
		nounType := k.NewName(tables.NAMED_TRAIT, "common", reader.At(r))
		k.NewValue(noun, nounType, true)
		//
		if det[0] != '$' {
			article := k.NewName(tables.NAMED_FIELD, "indefinite article", reader.At(r))
			k.NewValue(noun, article, det)
			if once := "common_noun"; k.Once(once) {
				indefinite := k.NewName(tables.NAMED_FIELD, "indefinite article", once)
				things := k.NewName(tables.NAMED_KIND, "things", once)
				k.NewPrimitive(tables.PRIM_TEXT, things, indefinite)
			}
		}
	}
	return
}

func imp_determiner(k *imp.Porter, r reader.Map) (ret string, err error) {
	return reader.String(r, "determiner")
}

// run: "{proper_name}"
// common / proper setting
func imp_proper_noun(k *imp.Porter, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "proper_noun"); e != nil {
		err = e
	} else if noun, e := imp_proper_name(k, m.MapOf("$PROPER_NAME")); e != nil {
		err = e
	} else {
		storyNouns.Add(noun)
		// set proper nounType to true ( implicitly defined by "noun" )
		nounType := k.NewName(tables.NAMED_TRAIT, "proper", reader.At(m))
		k.NewValue(noun, nounType, true)
	}
	return
}

// run: "{are_an} {*attribute:*trait} {kind:singular_kind} {?noun_relation}"
// ex. "(the box) is a closed container on the beach"
func imp_kind_of_noun(k *imp.Porter, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "kind_of_noun"); e != nil {
		err = e
	} else if kind, e := imp_singular_kind(k, m.MapOf("$KIND")); e != nil {
		err = e
	} else {
		var traits []ephemera.Named
		if e := reader.Repeats(m.SliceOf("$ATTRIBUTE"), func(el reader.Map) (err error) {
			if trait, e := imp_trait(k, el); e != nil {
				err = e
			} else {
				traits = append(traits, trait)
			}
			return
		}); e != nil {
			err = e
		} else {
			// we collect the nouns, but delay processing them till now.
			for _, noun := range storyNouns.names {
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
func imp_summary(k *imp.Porter, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "summary"); e != nil {
		err = e
	} else if lines, e := imp_line_expr(k, m.MapOf("$LINES")); e != nil {
		err = e
	} else {
		// declare the existence of the field "appearance"
		if once := "summary"; k.Once(once) {
			things := k.NewName(tables.NAMED_KIND, "things", once)
			appear := k.NewName(tables.NAMED_FIELD, "appearance", once)
			k.NewPrimitive(tables.PRIM_EXPR, things, appear)
		}
		prop := k.NewName(tables.NAMED_FIELD, "appearance", reader.At(m))
		noun := storyNouns.Last()
		k.NewValue(noun, prop, lines)
	}
	return
}

// run: "{are_being} {+attribute:trait}"
// ex. "(the box) is closed"
func imp_noun_attrs(k *imp.Porter, r reader.Map) (err error) {
	return reader.Repeats(r.SliceOf("$ATTRIBUTE"), func(el reader.Map) (err error) {
		if trait, e := imp_trait(k, el); e != nil {
			err = e
		} else {
			for _, noun := range storyNouns.names {
				k.NewValue(noun, trait, true) // the value of the trait for the noun is true
			}
		}
		return
	})
}

// fix... part of class attributes
func imp_attribute_phrase(k *imp.Porter, r reader.Map) (ret []ephemera.Named, err error) {
	err = reader.Repeats(r.SliceOf("$ATTRIBUTE"), func(el reader.Map) (err error) {
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
func imp_variable_decl(k *imp.Porter, r reader.Map) (retName, retType ephemera.Named, err error) {
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

func imp_variable_type(k *imp.Porter, r reader.Map) (ret ephemera.Named, err error) {
	err = reader.Option(r, "variable_type", reader.ReadMaps{
		"$PRIMITIVE": func(m reader.Map) (err error) {
			ret, err = imp_primitive_type(k, m)
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
// returns one of the evalType(s)
func imp_primitive_type(k *imp.Porter, r reader.Map) (ret ephemera.Named, err error) {
	if evalType, e := reader.Enum(r, "primitive_type", reader.Map{
		"$NUMBER": "number_eval",
		"$TEXT":   "text_eval",
		"$BOOL":   "bool_eval",
	}); e != nil {
		err = e
	} else {
		ret = k.NewName(tables.NAMED_TYPE, evalType.(string), reader.At(r))
	}
	return
}

// "{an} {kind of%kinds:plural_kinds} object"
// returns the name of "plural_kinds"
func imp_object_type(k *imp.Porter, r reader.Map) (ret ephemera.Named, err error) {
	if m, e := reader.Unpack(r, "object_type"); e != nil {
		err = e
	} else {
		ret, err = imp_plural_kinds(k, m.MapOf("$KINDS"))
	}
	return
}

func unimplemented(k *imp.Porter, r reader.Map) (err error) {
	return Unimplemented
}

const Unimplemented = errutil.Error("github.com/ionous/iffy/cmd/import/internal/Unimplemented")
