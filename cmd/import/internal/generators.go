package internal

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/tables"
)

var cmds = makeTypeMap(export.Runs)

// FIX -- keeping with the function parser model, swap these out to functions.
// the top level function should match the top level of "story"
// which a switch for parsing "story_statement" slats.

// story is a bunch of paragraphs
//make.run("story", "{+paragraph|ghost}");
func imp_story(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "story"); e != nil {
		err = e
	} else {
		err = k.repeats(m.SliceOf("$PARAGRAPH"), imp_paragraph)
	}
	return
}

// paragraph is a bunch of statements
func imp_paragraph(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "paragraph"); e != nil {
		err = e
	} else {
		k.nouns.Swap(nil)
		err = k.repeats(m.SliceOf("$STORY_STATEMENT"), imp_story_statement)
	}
	return
}

// run: "{+names} {noun_phrase}."
func imp_lede(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "lede"); e != nil {
		err = e
	} else if e := k.repeats(m.SliceOf("$NOUN"), imp_noun); e != nil {
		err = e
	} else if e := imp_noun_phrase(k, m.MapOf("$NOUN_PHRASE")); e != nil {
		err = e
	}
	return
}

// run: "{pronoun} {noun_phrase}."
func imp_tail(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "tail"); e != nil {
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
	if _, e := k.expectStr(r, "pronoun"); e != nil {
		err = e
	} else {
		// FIX: this can indicate plurality
	}
	return
}

// opt: "{kind_of_noun}, {noun_attrs}, or {noun_relation}"
func imp_noun_phrase(k *Importer, r reader.Map) (err error) {
	return k.expectOpt(r, "noun_phrase", map[string]Parse{
		"$KIND_OF_NOUN":  imp_kind_of_noun,
		"$NOUN_ATTRS":    imp_noun_attrs,
		"$NOUN_RELATION": imp_noun_relation,
	})
}

// run: "{?are_being} {relation} {+noun}"
// ex. (the cat and the hat) are in (the book)
// ex. (Hector and Maria) are suspicious of (Santa and Santana).
func imp_noun_relation(k *Importer, r reader.Map) (err error) {
	// unexpected type  wanted noun_relation at
	if m, e := k.expectSlat(r, "noun_relation"); e != nil {
		err = e
	} else {
		// fix? parse are_being
		relation := k.namedStr(m, tables.NAMED_VERB, "$RELATION")
		leadingNouns := k.nouns.Swap(nil)
		if e := k.repeats(m.SliceOf("$NOUN"), imp_noun); e != nil {
			err = e
		} else {
			trailingNouns := k.nouns.Swap(leadingNouns)
			//
			for _, n := range leadingNouns {
				for _, d := range trailingNouns {
					k.NewRelative(n, relation, d)
				}
			}
		}
	}
	return
}

// noun: "{proper_noun} or {common_noun}"
func imp_noun(k *Importer, r reader.Map) (err error) {
	// declare a noun class that has several default fields
	if once := "noun"; k.once(once) {
		things := k.Named(tables.NAMED_KIND, "things", once)
		nounType := k.Named(tables.NAMED_ASPECT, "nounType", once)
		common := k.Named(tables.NAMED_TRAIT, "common", once)
		proper := k.Named(tables.NAMED_TRAIT, "proper", once)
		k.NewPrimitive(tables.PRIM_ASPECT, things, nounType)
		k.NewTrait(common, nounType, 0)
		k.NewTrait(proper, nounType, 1)
	}
	return k.expectOpt(r, "noun", map[string]Parse{
		"$PROPER_NOUN": imp_proper_noun,
		"$COMMON_NOUN": imp_common_noun,
	})
}

// run: "{determiner} {common_name}"
func imp_common_noun(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "common_noun"); e != nil {
		err = e
	} else {
		id := r.StrOf(itemId)
		det := k.getStr(m, "$DETERMINER")
		noun := k.namedStr(m, tables.NAMED_NOUN, "$COMMON_NAME")
		k.nouns.Add(noun)
		// set common nounType to true ( implicitly defined by "noun" )
		nounType := k.Named(tables.NAMED_TRAIT, "common", id)
		k.NewValue(noun, nounType, true)
		//
		if det[0] != '$' {
			article := k.Named(tables.NAMED_FIELD, "indefinite article", id)
			k.NewValue(noun, article, det)
			if once := "common_noun"; k.once(once) {
				indefinite := k.Named(tables.NAMED_FIELD, "indefinite article", once)
				things := k.Named(tables.NAMED_KIND, "things", once)
				k.NewPrimitive(tables.PRIM_TEXT, things, indefinite)
			}
		}
	}
	return
}

// run: "{proper_name}"
// common / proper setting
func imp_proper_noun(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "proper_noun"); e != nil {
		err = e
	} else {
		id := r.StrOf(itemId)
		noun := k.namedStr(m, tables.NAMED_NOUN, "$PROPER_NAME")
		k.nouns.Add(noun)
		// set proper nounType to true ( implicitly defined by "noun" )
		nounType := k.Named(tables.NAMED_TRAIT, "proper", id)
		k.NewValue(noun, nounType, true)
	}
	return nil
}

// run: "{are_an} {*attribute} {kind} {?noun_relation}"
// ex. "(the box) is a closed container on the beach"
func imp_kind_of_noun(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "kind_of_noun"); e != nil {
		err = e
	} else {
		kind := k.namedStr(m, tables.NAMED_KIND, "$KIND")
		for _, noun := range k.nouns.Named {
			k.NewNoun(noun, kind)
		}
		if e := imp_noun_attrs(k, m); e != nil {
			err = e
		} else if v := m.MapOf("$NOUN_RELATION"); len(v) != 0 {
			err = imp_noun_relation(k, v)
		}
	}
	return
}

// run: "{The [summary] is:: %lines}"
func imp_summary(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "summary"); e != nil {
		err = e
	} else {
		id := r.StrOf(itemId)
		// declare the existence of the field "appearance"
		if once := "summary"; k.once(once) {
			things := k.Named(tables.NAMED_KIND, "things", once)
			appear := k.Named(tables.NAMED_FIELD, "appearance", once)
			k.NewPrimitive(tables.PRIM_EXPR, things, appear)
		}
		prop := k.Named(tables.NAMED_FIELD, "appearance", id)
		noun := k.nouns.Last()
		val := k.namedStr(m, tables.PRIM_EXPR, "$LINES")
		k.NewValue(noun, prop, val)
	}
	return
}

// run: "{are_being} {+attribute}"
// ex. "(the box) is closed"
func imp_noun_attrs(k *Importer, r reader.Map) (err error) {
	defer k.on(tables.NAMED_TRAIT, func(trait ephemera.Named) {
		for _, noun := range k.nouns.Named {
			k.NewValue(noun, trait, true)
		}
	})()
	for _, it := range r.SliceOf("$ATTRIBUTE") {
		k.catStr(reader.Box(it), tables.NAMED_TRAIT)
	}
	return
}

func imp_attribute_phrase(k *Importer, r reader.Map) error {
	return imp_noun_attrs(k, r)
}

// "{type:variable_type} ( called {name:variable_name|quote} )"
func imp_variable_decl(k *Importer, r reader.Map,
	primDecl func(ephemera.Named, string),
	objDecl func(ephemera.Named, ephemera.Named)) (err error) {
	if m, e := k.expectSlat(r, "variable_decl"); e != nil {
		err = e
	} else if name, e := imp_variable_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else {
		err = imp_variable_type(k, m.MapOf("$TYPE"),
			func(eval string) { primDecl(name, eval) },
			func(kind ephemera.Named) { objDecl(name, kind) })
	}
	return
}

func imp_pattern_type(k *Importer, r reader.Map) (ret string, err error) {
	// err = k.expectOpt(m, "pattern_type", map[string]Parse{
	// 	"$ACTIVITY": func(k *Importer, r reader.Map) (err error) {
	// 		ret, err = imp_pattern_activity(k, m)
	// 		return
	// 	},
	// 	"$VALUE": func(k *Importer, r reader.Map) (err error) {
	// 		 err = imp_variable_type(k, m)
	// 		return
	// 	},
	// })
	err = Unimplemented
	return
}

func imp_pattern_activity(k *Importer, r reader.Map) (ret string, err error) {
	if e := k.expectConst(r, "pattern_activity", "$ACTIVITY"); e != nil {
		err = e
	} else {
		ret = tables.EVAL_PROG
	}
	return
}

func imp_variable_type(k *Importer, r reader.Map,
	prim func(string),
	obj func(ephemera.Named)) (err error) {
	return k.expectOpt(r, "variable_type", map[string]Parse{
		"$PRIMITIVE": func(k *Importer, m reader.Map) (err error) {
			if n, e := imp_primitve_type(k, m); e != nil {
				err = e
			} else {
				prim(n)
			}
			return
		},
		"$OBJECT": func(k *Importer, m reader.Map) (err error) {
			if n, e := imp_object_type(k, m); e != nil {
				err = e
			} else {
				obj(n)
			}
			return
		},
	})
}

// "{a number%number}, {some text%text}, or {a true/false value%bool}"
// returns one of the evalType(s)
func imp_primitve_type(k *Importer, r reader.Map) (ret string, err error) {
	if n, e := k.expectEnum(r, "primitive_type", map[string]interface{}{
		"$NUMBER": tables.EVAL_DIGI,
		"$TEXT":   tables.EVAL_TEXT,
		"$BOOL":   tables.EVAL_BOOL,
	}); e != nil {
		err = e
	} else {
		ret = n.(string)
	}
	return
}

// "{an} {kind of%kinds:plural_kinds} object"
// returns the name of "plural_kinds"
func imp_object_type(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	if m, e := k.expectSlat(r, "object_type"); e != nil {
		err = e
	} else {
		ret, err = imp_plural_kinds(k, m.MapOf("$KINDS"))
	}
	return
}

func unimplemented(k *Importer, r reader.Map) (err error) {
	return Unimplemented
}

const Unimplemented = errutil.Error("github.com/ionous/iffy/cmd/import/internal/Unimplemented")
