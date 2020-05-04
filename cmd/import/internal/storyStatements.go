package internal

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

func imp_story_statement(k *Importer, r reader.Map) (err error) {
	return k.expectSlot(r, "story_statement", map[string]Parse{
		"certainties":              imp_certainties,
		"class_attributes":         imp_class_attributes,
		"kinds_of_quality":         imp_kinds_of_quality,
		"kinds_of_thing":           imp_kinds_of_thing,
		"kinds_possess_properties": imp_kinds_possess_properties,
		"noun_assignment":          imp_noun_assignment,
		"noun_statement":           imp_noun_statement,
		// "pattern_decl":             imp_pattern_decl,
		"pattern_variables_decl": imp_pattern_variables_decl,
		"quality_attributes":     imp_quality_attributes,
		"relative_to_noun":       imp_relative_to_noun,
		"test_statement":         imp_test_statement,
	})
}

//"{plural_kinds} {are_being} {certainty} {trait}.");
// horses are usually fast.
func imp_certainties(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "certainties"); e != nil {
		err = e
	} else if certainty, e := imp_certainty(k, m.MapOf("$CERTAINTY")); e != nil {
		err = e
	} else if trait, e := imp_trait(k, m.MapOf("$ATTRIBUTE")); e != nil {
		err = e
	} else if kind, e := imp_plural_kinds(k, m.MapOf("$PLURAL_KINDS")); e != nil {
		err = e
	} else {
		k.NewCertainty(certainty, trait, kind)
	}
	return
}

func imp_certainty(k *Importer, r reader.Map) (ret string, err error) {
	if n, e := k.expectEnum(r, "certainty", map[string]interface{}{
		"$ALWAYS":  "always",
		"$NEVER":   "never",
		"$SELDOM":  "seldom",
		"$USUALLY": "usually",
	}); e != nil {
		err = e
	} else {
		ret = n.(string)
	}
	return
}

// "{plural_kinds} {attribute_phrase}"
// ex. animals are fast or slow.
func imp_class_attributes(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "class_attributes"); e != nil {
		err = e
	} else if kind, e := imp_plural_kinds(k, m.MapOf("$PLURAL_KINDS")); e != nil {
		err = e
	} else {
		//
		var traits []ephemera.Named
		defer k.on(tables.NAMED_TRAIT, func(trait ephemera.Named) {
			traits = append(traits, trait)
		})()
		if e := imp_attribute_phrase(k, m.MapOf("$ATTRIBUTE_PHRASE")); e != nil {
			err = e
		} else {
			// create an implied aspect named after the first trait
			// fix? maybe we should include the columns of named in the returned struct so we can pick out the source better.
			aspect := k.Named(tables.NAMED_ASPECT, traits[0].String(), m.StrOf(itemId))
			k.NewPrimitive(tables.PRIM_ASPECT, kind, aspect)
		}
	}
	return
}

// "{qualities} {are_an} kind of value."
// ex. colors are a kind of value
func imp_kinds_of_quality(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "kinds_of_quality"); e != nil {
		err = e
	} else {
		aspect := k.namedStr(m, tables.NAMED_ASPECT, "$QUALITY")
		k.NewAspect(aspect)
	}
	return
}

// "{plural_kinds} {are_an} kind of {kind}."
// ex. "cats are a kind of animal"
func imp_kinds_of_thing(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "kinds_of_thing"); e != nil {
		err = e
	} else if kind, e := imp_plural_kinds(k, m.MapOf("$PLURAL_KINDS")); e != nil {
		err = e
	} else {
		parent := k.namedStr(m, tables.NAMED_KIND, "$KIND")
		k.NewKind(kind, parent)
	}
	return
}

// {plural_kinds} have {determiner} {primitive_type} called {property}.
// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func imp_kinds_possess_properties(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "kinds_possess_properties"); e != nil {
		err = e
	} else if kind, e := imp_plural_kinds(k, m.MapOf("$PLURAL_KINDS")); e != nil {
		err = e
	} else {
		prop := k.namedStr(m, tables.NAMED_KIND, "$PROPERTY")
		primType := k.getStr(m, "$PRIMITIVE_TYPE")
		k.NewPrimitive(primType, kind, prop)
	}
	return
}

// run "The {property} of {+noun} is the {[text]:: %lines}"
// ex. The description of the nets is xxx
func imp_noun_assignment(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "noun_assignment"); e != nil {
		err = e
	} else {
		val := k.namedStr(m, tables.PRIM_EXPR, "$LINES")
		prop := k.namedStr(m, tables.NAMED_FIELD, "$PROPERTY")
		defer k.on(tables.NAMED_NOUN, func(noun ephemera.Named) {
			k.NewValue(noun, prop, val)
		})()
		err = k.repeats(m.SliceOf("$NOUN"), imp_noun)
	}
	return
}

func imp_noun_statement(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "noun_statement"); e != nil {
		err = e
	} else if e := imp_lede(k, m.MapOf("$LEDE")); e != nil {
		err = e
	} else if e := k.repeats(m.SliceOf("$TAIL"), imp_tail); e != nil {
		err = e
	} else if v := m.MapOf("$SUMMARY"); len(v) != 0 {
		err = imp_summary(k, v)
	}
	return
}

// Adds a NewPattern, and optionally some associated pattern parameters.
// {name:pattern_name|quote} determines {type:pattern_type}.
// {optvars?pattern_variables_tail}
func imp_pattern_decl(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "pattern_decl"); e != nil {
		err = e
	} else if pat, e := imp_pattern_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else if typ, e := imp_pattern_type(k, m.MapOf("$TYPE")); e != nil {
		err = e
	} else {
		k.NewPattern(pat, typ)
		if tail := m.MapOf("$OPTVARS"); len(tail) > 0 {
			err = imp_pattern_variables_tail(k, pat, tail)
		}
	}
	return
}

// `Pattern variables: Storage for values used during the execution of a pattern.`)
// {+variable_decl|comma-and}.",
func imp_pattern_variables_tail(k *Importer, pat ephemera.Named, r reader.Map) (err error) {
	return k.repeats(r.SliceOf("$VARIABLE_DECL"),
		func(k *Importer, m reader.Map) (err error) {
			if paramName, paramType, e := imp_variable_decl(k, m); e != nil {
				err = e
			} else {
				k.NewPatternParam(pat, paramName, paramType)
			}
			return
		})
}

// "The pattern {pattern_name|quote} uses {+variable_decl|comma-and}."
func imp_pattern_variables_decl(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "pattern_variables_decl"); e != nil {
		err = e
	} else if pat, e := imp_pattern_name(k, m.MapOf("$PATTERN_NAME")); e != nil {
		err = e
	} else {
		// reuse, works because they have the same $name.
		err = imp_pattern_variables_tail(k, pat, m)
	}
	return
}

// "{qualities} {attribute_phrase}"
// (the) colors are red, blue, or green.
func imp_quality_attributes(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "quality_attributes"); e != nil {
		err = e
	} else {
		aspect := k.namedStr(m, tables.NAMED_ASPECT, "$QUALITIES")
		rank := 0
		defer k.on(tables.NAMED_TRAIT, func(trait ephemera.Named) {
			k.NewTrait(trait, aspect, rank)
			rank += 1
		})()
		if e := imp_attribute_phrase(k, m.MapOf("$ATTRIBUTE_PHRASE")); e != nil {
			err = e
		}
	}
	return
}

// "{relation} {+noun} {are_being} {+noun}."
// ex. On the beach are shells.
func imp_relative_to_noun(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "relative_to_noun"); e != nil {
		err = e
	} else {
		relation := k.namedStr(m, tables.NAMED_VERB, "$RELATION")
		if e := k.repeats(m.SliceOf("$NOUN"), imp_noun); e != nil {
			err = errutil.Append(err, e)
		} else {
			leadingNouns := k.nouns.Swap(nil)
			if e := k.repeats(m.SliceOf("$NOUN1"), imp_noun); e != nil {
				err = e
			} else {
				trailingNouns := k.nouns.Swap(leadingNouns)
				for _, a := range leadingNouns {
					for _, b := range trailingNouns {
						k.NewRelative(a, relation, b)
					}
				}
			}
		}
	}
	return
}

func imp_test_statement(k *Importer, outer reader.Map) (err error) {
	return k.expectSlot(outer, "testing", map[string]Parse{
		"test_output": imp_test_output,
	})
}

func imp_test_output(k *Importer, outer reader.Map) (err error) {
	if m, e := k.expectSlat(outer, "test_output"); e != nil {
		err = e
	} else {
		test := k.namedStr(m, tables.NAMED_TEST, "$TEST_NAME")
		expect := k.getStr(m, "$LINES")
		var prog check.TestOutput
		if e := ReadProg(&prog, reader.Unbox(outer), cmds); e != nil {
			err = e
		} else if p, e := k.newProg("test", prog); e != nil {
			err = e
		} else {
			k.NewTest(test, p, expect)
		}
	}
	return
}
