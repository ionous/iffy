package internal

import (
	"bytes"
	"encoding/gob"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/tables"
)

var cmds = makeTypeMap(export.Runs)

// FIX -- keeping with the function parser model, swap these out to functions.
// the top level function should match the top level of "story"
// which a switch for parsing "story_statement" slats.
// would remove Importeephemera.NamedGen
var generators = map[string]Parse{
	"test": func(k *Importer, m reader.Map) (err error) {
		test := k.namedStr(m, tables.NAMED_TEST, "$TEST_NAME")
		expect := k.getStr(m, "$LINES")

		var prog check.Test // parentItem is { id:..., type:"test", value:... }
		if e := readProg(&prog, reader.Unbox(k.parentItem), cmds); e != nil {
			err = e
		} else {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if e := enc.Encode(prog); e != nil {
				err = e
			} else {
				prog := k.NewProg(k.currType, buf.Bytes())
				k.NewTest(test, prog, expect)
			}
		}
		return
	},
	// story is a bunch of paragraphs
	"story": func(k *Importer, m reader.Map) error {
		return k.parseSlice(m.SliceOf("$PARAGRAPH"))
	},

	// paragraph is a bunch of statements
	"paragraph": func(k *Importer, m reader.Map) error {
		return k.parseSlice(m.SliceOf("$STORY_STATEMENT"))
	},

	"story_statement": func(k *Importer, m reader.Map) error {
		k.nouns.Swap(nil)
		return k.parseItem(m)
	},

	"noun_statement": func(k *Importer, m reader.Map) error {
		return errutil.Append(
			k.parseItem(m.MapOf("$LEDE")),
			k.parseSlice(m.SliceOf("$TAIL")),
			k.parseItem(m.MapOf("$SUMMARY")))
	},

	// run "The {property} of {+noun} is the {[text]:: %lines}"
	// ex. The description of the nets is xxx
	"noun_assignment": func(k *Importer, m reader.Map) error {
		val := k.namedStr(m, tables.PRIM_EXPR, "$LINES")
		prop := k.namedStr(m, tables.NAMED_FIELD, "$PROPERTY")
		defer k.on(tables.NAMED_NOUN, func(noun ephemera.Named) {
			k.NewValue(noun, prop, val)
		})()
		return k.parseSlice(m.SliceOf("$NOUN"))
	},

	// "{relation} {+noun} {are_being} {+noun}."
	// ex. On the beach are shells.
	"relative_to_noun": func(k *Importer, m reader.Map) (err error) {
		relation := k.namedStr(m, tables.NAMED_VERB, "$RELATION")
		//
		if e := k.parseSlice(m.SliceOf("$NOUN")); e != nil {
			err = errutil.Append(err, e)
		}
		leadingNouns := k.nouns.Swap(nil)
		if e := k.parseSlice(m.SliceOf("$NOUN1")); e != nil {
			err = errutil.Append(err, e)
		}
		trailingNouns := k.nouns.Swap(leadingNouns)
		//
		for _, a := range leadingNouns {
			for _, b := range trailingNouns {
				k.NewRelative(a, relation, b)
			}
		}
		return err
	},

	// "{plural_kinds} {are_an} kind of {kind}."
	// ex. "cats are a kind of animal"
	"kinds_of_thing": func(k *Importer, m reader.Map) error {
		kind := k.namedStr(m, tables.NAMED_KIND, "$PLURAL_KINDS")
		parent := k.namedStr(m, tables.NAMED_KIND, "$KIND")
		k.NewKind(kind, parent)
		return nil
	},

	// "{qualities} {are_an} kind of value."
	// ex. colors are a kind of value
	"kinds_of_quality": func(k *Importer, m reader.Map) error {
		aspect := k.namedStr(m, tables.NAMED_ASPECT, "$QUALITY")
		k.NewAspect(aspect)
		return nil
	},

	// "{plural_kinds} {attribute_phrase}"
	// ex. animals are fast or slow.
	"class_attributes": func(k *Importer, m reader.Map) (err error) {
		kind := k.namedStr(m, tables.NAMED_KIND, "$PLURAL_KINDS")
		//
		var traits []ephemera.Named
		defer k.on(tables.NAMED_TRAIT, func(trait ephemera.Named) {
			traits = append(traits, trait)
		})()
		err = k.parseItem(m.MapOf("$ATTRIBUTE_PHRASE"))
		// create an implied aspect named after the first trait
		// fix? maybe we should include the columns of named in the returned struct so we can pick out the source better.
		aspect := k.Named(tables.NAMED_ASPECT, traits[0].String(), m.StrOf(itemId))
		k.NewPrimitive(tables.PRIM_ASPECT, kind, aspect)
		return
	},

	// "{qualities} {attribute_phrase}"
	// (the) colors are red, blue, or green.
	"quality_attributes": func(k *Importer, m reader.Map) error {
		aspect := k.namedStr(m, tables.NAMED_ASPECT, "$QUALITIES")
		rank := 0
		defer k.on(tables.NAMED_TRAIT, func(trait ephemera.Named) {
			k.NewTrait(trait, aspect, rank)
			rank += 1
		})()
		return k.parseItem(m.MapOf("$ATTRIBUTE_PHRASE"))
	},

	//"{plural_kinds} {are_being} {certainty} {attribute}.");
	// horses are usually fast.
	"certainties": func(k *Importer, m reader.Map) error {
		certainty := k.getStr(m, "$CERTAINTY")
		trait := k.namedStr(m, tables.NAMED_TRAIT, "$ATTRIBUTE")
		kind := k.namedStr(m, tables.NAMED_KIND, "$PLURAL_KINDS")
		k.NewCertainty(certainty, trait, kind)
		return nil
	},

	// {plural_kinds} have {determiner} {primitive_type} called {property}.
	// ex. cats have some text called breed.
	// ex. horses have an aspect called speed.
	"kinds_possess_properties": func(k *Importer, m reader.Map) error {
		kind := k.namedStr(m, tables.NAMED_KIND, "$PLURAL_KINDS")
		prop := k.namedStr(m, tables.NAMED_KIND, "$PROPERTY")
		primType := k.getStr(m, "$PRIMITIVE_TYPE")
		k.NewPrimitive(primType, kind, prop)
		return nil
	},

	// run: "{+names} {noun_phrase}."
	"lede": func(k *Importer, m reader.Map) error {
		return errutil.Append(
			k.parseSlice(m.SliceOf("$NOUN")),
			k.parseItem(m.MapOf("$NOUN_PHRASE")))
	},

	// run: "{pronoun} {noun_phrase}."
	"tail": func(k *Importer, m reader.Map) error {
		return errutil.Append(
			k.parseItem(m.MapOf("$PRONOUN")),
			k.parseItem(m.MapOf("$NOUN_PHRASE")))
	},

	// opt: "{kind_of_noun}, {noun_attrs}, or {noun_relation}"
	"noun_phrase": func(k *Importer, m reader.Map) error {
		return k.parseItem(m)
	},

	// run: "{?are_being} {relation} {+noun}"
	// ex. (the cat and the hat) are in (the book)
	// ex. (Hector and Maria) are suspicious of (Santa and Santana).
	"noun_relation": func(k *Importer, m reader.Map) (err error) {
		relation := k.namedStr(m, tables.NAMED_VERB, "$RELATION")
		//
		leadingNouns := k.nouns.Swap(nil)
		err = k.parseSlice(m.SliceOf("$NOUN"))
		trailingNouns := k.nouns.Swap(leadingNouns)
		//
		for _, n := range leadingNouns {
			for _, d := range trailingNouns {
				k.NewRelative(n, relation, d)
			}
		}
		return
	},

	// noun: "{proper_noun} or {common_noun}"
	"noun": func(k *Importer, m reader.Map) error {
		once := "noun"
		if k.once(once) {
			things := k.Named(tables.NAMED_KIND, "things", once)
			nounType := k.Named(tables.NAMED_ASPECT, "nounType", once)
			common := k.Named(tables.NAMED_TRAIT, "common", once)
			proper := k.Named(tables.NAMED_TRAIT, "proper", once)
			k.NewPrimitive(tables.PRIM_ASPECT, things, nounType)
			k.NewTrait(common, nounType, 0)
			k.NewTrait(proper, nounType, 1)
		}
		return k.parseItem(m)
	},

	// run: "{determiner} {common_name}"
	"common_noun": func(k *Importer, m reader.Map) error {
		id := k.currId
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
			once := "common_noun"
			if k.once(once) {
				indefinite := k.Named(tables.NAMED_FIELD, "indefinite article", once)
				things := k.Named(tables.NAMED_KIND, "things", once)
				k.NewPrimitive(tables.PRIM_TEXT, things, indefinite)
			}
		}
		return nil
	},

	// run: "{proper_name}"
	// common / proper setting
	"proper_noun": func(k *Importer, m reader.Map) error {
		id := k.currId
		noun := k.namedStr(m, tables.NAMED_NOUN, "$PROPER_NAME")
		k.nouns.Add(noun)
		// set proper nounType to true ( implicitly defined by "noun" )
		nounType := k.Named(tables.NAMED_TRAIT, "proper", id)
		k.NewValue(noun, nounType, true)
		return nil
	},

	// run: "{are_an} {*attribute} {kind} {?noun_relation}"
	// ex. "(the box) is a closed container on the beach"
	"kind_of_noun": func(k *Importer, m reader.Map) error {
		//
		kind := k.namedStr(m, tables.NAMED_KIND, "$KIND")
		for _, noun := range k.nouns.Named {
			k.NewNoun(noun, kind)
		}
		parseAttrs(k, m)
		// noun relation takes care of itself --
		// relating the new nouns to the existing nouns.
		return k.parseItem(m.MapOf("$NOUN_RELATION"))
	},

	// run: "{are_being} {+attribute}"
	// ex. "(the box) is closed"
	"noun_attrs": parseAttrs,

	// run: "{are_either} {+attribute}."
	"attribute_phrase": parseAttrs,

	// run: "{The [summary] is:: %lines}"
	"summary": func(k *Importer, m reader.Map) error {
		once := "summary"
		id := k.currId
		if k.once(once) { // declare the existence of the field "appearance"
			things := k.Named(tables.NAMED_KIND, "things", once)
			appear := k.Named(tables.NAMED_FIELD, "appearance", once)
			k.NewPrimitive(tables.PRIM_EXPR, things, appear)
		}
		prop := k.Named(tables.NAMED_FIELD, "appearance", id)
		noun := k.nouns.Last()
		val := k.namedStr(m, tables.PRIM_EXPR, "$LINES")
		k.NewValue(noun, prop, val)
		return nil
	},

	"pattern_decl": func(k *Importer, m reader.Map) (err error) {

		// name := k.namedStr(m, tables.NAMED_PATTERN, "$NAME")

		// currValue := m.MapOf("$TYPE")
		// //tail := m.MapOf("PATTERN_VARIABLES_TAIL")

		// if patternType, e := patternType(k, m.MapOf("$TYPE")); e != nil {
		// 	err = e
		// } else {
		// }

		// can i do something .on()?
		// i need to look for a variable?type with a value of $PRIMITIVE or $OBJECT

		// k.NewPattern(

		return
	},

	"pattern_variables_tail": func(k *Importer, m reader.Map) (err error) {
		return
	},
	// "pattern_variables_decl": patternVariablesDecl,
}

// "The pattern {pattern_name|quote} uses {+variable_decl|comma-and}."
func patternVariablesDecl(k *Importer, m reader.Map) (err error) {
	if m, e := k.expectSlat(m, "pattern_variables_decl"); e != nil {
		err = e
	} else if pat, e := patternName(k, m.MapOf("$PATTERN_NAME")); e != nil {
		err = e
	} else {
		err = k.repeats(m.SliceOf("$VARIABLE_DECL"), func(m reader.Map) (err error) {
			return variableDecl(k, m,
				func(param ephemera.Named, eval string) {
					// fix? maybe itd be better to just manufacture some built in fake "kinds" for this.
					// "numbers", "strings", "booleans", "programs".
					k.NewPatternedEval(pat, param, eval)
				},
				func(param ephemera.Named, kind ephemera.Named) {
					k.NewPatternedKind(pat, param, kind)
				},
			)
		})
	}
	return
}

// "{type:variable_type} ( called {name:variable_name|quote} )"
func variableDecl(k *Importer, m reader.Map,
	primDecl func(ephemera.Named, string),
	objDecl func(ephemera.Named, ephemera.Named)) (err error) {
	if m, e := k.expectSlat(m, "variable_decl"); e != nil {
		err = e
	} else if name, e := variableName(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else {
		err = variableType(k, m.MapOf("$TYPE"),
			func(eval string) { primDecl(name, eval) },
			func(kind ephemera.Named) { objDecl(name, kind) })
	}
	return
}

func patternType(k *Importer, m reader.Map) (ret string, err error) {
	// err = k.expectOpt(m, "pattern_type", map[string]Parse{
	// 	"$ACTIVITY": func(k *Importer, m reader.Map) (err error) {
	// 		ret, err = patternActivity(k, m)
	// 		return
	// 	},
	// 	"$VALUE": func(k *Importer, m reader.Map) (err error) {
	// 		 err = variableType(k, m)
	// 		return
	// 	},
	// })
	err = Unimplemented
	return
}

func patternActivity(k *Importer, m reader.Map) (ret string, err error) {
	if e := k.expectConst(m, "pattern_activity", "$ACTIVITY"); e != nil {
		err = e
	} else {
		ret = tables.EVAL_PROG
	}
	return
}

func variableType(k *Importer, m reader.Map,
	prim func(string),
	obj func(ephemera.Named)) (err error) {
	return k.expectOpt(m, "variable_type", map[string]Parse{
		"$PRIMITIVE": func(k *Importer, m reader.Map) (err error) {
			if n, e := primitiveType(k, m); e != nil {
				err = e
			} else {
				prim(n)
			}
			return
		},
		"$OBJECT": func(k *Importer, m reader.Map) (err error) {
			if n, e := objectType(k, m); e != nil {
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
func primitiveType(k *Importer, m reader.Map) (ret string, err error) {
	if n, e := k.expectStr(m, "primitive_type"); e != nil {
		err = e
	} else {
		switch n {
		case "$NUMBER":
			ret = tables.EVAL_DIGI
		case "$TEXT":
			ret = tables.EVAL_TEXT
		case "$BOOL":
			ret = tables.EVAL_BOOL
		default:
			err = errutil.New("unexpected primitive type", n)
		}
	}
	return
}

// "{an} {kind of%kinds:plural_kinds} object"
// returns the name of "plural_kinds"
func objectType(k *Importer, m reader.Map) (ret ephemera.Named, err error) {
	if m, e := k.expectSlat(m, "object_type"); e != nil {
		err = e
	} else {
		ret, err = pluralKinds(k, m.MapOf("$KINDS"))
	}
	return
}

func unimplemented(k *Importer, m reader.Map) error {
	return Unimplemented
}

const Unimplemented = errutil.Error("github.com/ionous/iffy/cmd/import/internal/Unimplemented")

// make.str("pattern_name");
func patternName(k *Importer, m reader.Map) (ret ephemera.Named, err error) {
	return k.namedType(m, "pattern_name")
}

// make.str("variable_name");
func variableName(k *Importer, m reader.Map) (ret ephemera.Named, err error) {
	return k.namedType(m, "variable_name")
}

func pluralKinds(k *Importer, m reader.Map) (ret ephemera.Named, err error) {
	return k.namedType(m, "plural_kinds")
}
