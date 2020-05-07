package internal

import (
	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

func imp_story_statement(k *Importer, r reader.Map) (err error) {
	return reader.Slot(r, "story_statement", reader.ReadMaps{
		"certainties":              k.bind(imp_certainties),
		"class_attributes":         k.bind(imp_class_attributes),
		"kinds_of_quality":         k.bind(imp_kinds_of_quality),
		"kinds_of_thing":           k.bind(imp_kinds_of_thing),
		"kinds_possess_properties": k.bind(imp_kinds_possess_properties),
		"noun_assignment":          k.bind(imp_noun_assignment),
		"noun_statement":           k.bind(imp_noun_statement),
		// "pattern_decl":             k.bind(imp_pattern_decl),
		"pattern_variables_decl": k.bind(imp_pattern_variables_decl),
		"quality_attributes":     k.bind(imp_quality_attributes),
		"relative_to_noun":       k.bind(imp_relative_to_noun),
		"test_statement":         k.bind(imp_test_statement),
	})
}

//"{plural_kinds} {are_being} {certainty} {trait}.");
// horses are usually fast.
func imp_certainties(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Slat(r, "certainties"); e != nil {
		err = e
	} else if certainty, e := imp_certainty(k, m.MapOf("$CERTAINTY")); e != nil {
		err = e
	} else if trait, e := imp_trait(k, m.MapOf("$TRAIT")); e != nil {
		err = e
	} else if kind, e := imp_plural_kinds(k, m.MapOf("$PLURAL_KINDS")); e != nil {
		err = e
	} else {
		k.eph.NewCertainty(certainty, trait, kind)
	}
	return
}

func imp_certainty(k *Importer, r reader.Map) (ret string, err error) {
	if n, e := reader.Enum(r, "certainty", map[string]interface{}{
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
	if m, e := reader.Slat(r, "class_attributes"); e != nil {
		err = e
	} else if kind, e := imp_plural_kinds(k, m.MapOf("$PLURAL_KINDS")); e != nil {
		err = e
	} else if traits, e := imp_attribute_phrase(k, m.MapOf("$ATTRIBUTE_PHRASE")); e != nil {
		err = e
	} else {
		// create an implied aspect named after the first trait
		// fix? maybe we should include the columns of named in the returned struct so we can pick out the source better.
		aspect := k.eph.Named(tables.NAMED_ASPECT, traits[0].String(), m.StrOf(reader.ItemId))
		k.eph.NewPrimitive(tables.PRIM_ASPECT, kind, aspect)
	}
	return
}

// "{qualities} {are_an} kind of value."
// ex. colors are a kind of value
func imp_kinds_of_quality(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Slat(r, "kinds_of_quality"); e != nil {
		err = e
	} else if _, e := imp_qualities(k, m.MapOf("$QUALITIES")); e != nil {
		err = e
	}
	return
}

// "{plural_kinds} {are_an} kind of {kind}."
// ex. "cats are a kind of animal"
func imp_kinds_of_thing(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Slat(r, "kinds_of_thing"); e != nil {
		err = e
	} else if kind, e := imp_plural_kinds(k, m.MapOf("$PLURAL_KINDS")); e != nil {
		err = e
	} else if parent, e := imp_singular_kind(k, m.MapOf("$KIND")); e != nil {
		err = e
	} else {
		k.eph.NewKind(kind, parent)
	}
	return
}

// {plural_kinds} have {determiner} {primitive_type} called {property}.
// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func imp_kinds_possess_properties(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Slat(r, "kinds_possess_properties"); e != nil {
		err = e
	} else if _, e := imp_plural_kinds(k, m.MapOf("$PLURAL_KINDS")); e != nil {
		err = e
	} else if _, e := imp_determiner(k, m.MapOf("$DETERMINER")); e != nil {
		err = e
	} else {
		// opt.property_phrase missing
		err = Unimplemented
	}
	return
}

// run "The {property} of {+noun} is the {[text]:: %lines}"
// ex. The description of the nets is xxx
func imp_noun_assignment(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Slat(r, "noun_assignment"); e != nil {
		err = e
	} else if lines, e := imp_line_expr(k, m.MapOf("$LINES")); e != nil {
		err = e
	} else if prop, e := imp_property(k, m.MapOf("$PROPERTY")); e != nil {
		err = e
	} else if e := reader.Repeats(m.SliceOf("$NOUN"), k.bind(imp_noun)); e != nil {
		err = e
	} else {
		for _, noun := range k.nouns.Named {
			k.eph.NewValue(noun, prop, lines)
		}
	}
	return
}

func imp_noun_statement(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Slat(r, "noun_statement"); e != nil {
		err = e
	} else if e := imp_lede(k, m.MapOf("$LEDE")); e != nil {
		err = e
	} else if e := reader.Repeats(m.SliceOf("$TAIL"), k.bind(imp_tail)); e != nil {
		err = e
	} else if v := m.MapOf("$SUMMARY"); len(v) != 0 {
		err = imp_summary(k, v)
	}
	return
}

// Adds a NewPatternType, and optionally some associated pattern parameters.
// {name:pattern_name|quote} determines {type:pattern_type}.
// {optvars?pattern_variables_tail}
func imp_pattern_decl(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Slat(r, "pattern_decl"); e != nil {
		err = e
	} else if pat, e := imp_pattern_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else if typ, e := imp_pattern_type(k, m.MapOf("$TYPE")); e != nil {
		err = e
	} else {
		k.eph.NewPatternType(pat, typ)
		if tail := m.MapOf("$OPTVARS"); len(tail) > 0 {
			err = imp_pattern_variables_tail(k, pat, tail)
		}
	}
	return
}

// `Pattern variables: Storage for values used during the execution of a pattern.`)
// {+variable_decl|comma-and}.",
func imp_pattern_variables_tail(k *Importer, pat ephemera.Named, r reader.Map) (err error) {
	return reader.Repeats(r.SliceOf("$VARIABLE_DECL"),
		func(m reader.Map) (err error) {
			if paramName, paramType, e := imp_variable_decl(k, m); e != nil {
				err = e
			} else {
				k.eph.NewPatternParam(pat, paramName, paramType)
			}
			return
		})
}

// "The pattern {pattern_name|quote} uses {+variable_decl|comma-and}."
func imp_pattern_variables_decl(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Slat(r, "pattern_variables_decl"); e != nil {
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
	if m, e := reader.Slat(r, "quality_attributes"); e != nil {
		err = e
	} else if aspect, e := imp_qualities(k, r.MapOf("$QUALITIES")); e != nil {
		err = e
	} else if traits, e := imp_attribute_phrase(k, m.MapOf("$ATTRIBUTE_PHRASE")); e != nil {
		err = e
	} else {
		for rank, trait := range traits {
			k.eph.NewTrait(trait, aspect, rank)
		}
	}
	return
}

func imp_qualities(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	if n, e := reader.String(r.MapOf("$QUALITIES"), "qualities"); e != nil {
		err = e
	} else {
		ret = k.eph.Named(tables.NAMED_ASPECT, n, reader.At(r))
	}
	return
}

// "{relation} {+noun} {are_being} {+noun}."
// ex. On the beach are shells.
func imp_relative_to_noun(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Slat(r, "relative_to_noun"); e != nil {
		err = e
	} else if relation, e := imp_relation(k, m.MapOf("$RELATION")); e != nil {
		err = e
	} else if e := reader.Repeats(m.SliceOf("$NOUN"), k.bind(imp_noun)); e != nil {
		err = e
	} else {
		leadingNouns := k.nouns.Swap(nil)
		if e := reader.Repeats(m.SliceOf("$NOUN1"), k.bind(imp_noun)); e != nil {
			err = e
		} else {
			trailingNouns := k.nouns.Swap(leadingNouns)
			for _, a := range leadingNouns {
				for _, b := range trailingNouns {
					k.eph.NewRelative(a, relation, b)
				}
			}
		}
	}
	return
}

func imp_test_statement(k *Importer, r reader.Map) (err error) {
	if n, e := imp_test_name(k, r.MapOf("name")); e != nil {
		err = e
	} else {
		err = reader.Slot(r, "testing", reader.ReadMaps{
			"test_output": func(m reader.Map) error {
				return imp_test_output(k, n, m)
			},
		})
	}
	return
}

func imp_test_output(k *Importer, test ephemera.Named, r reader.Map) (err error) {
	var prog check.TestOutput
	if m, e := reader.Slat(r, "test_output"); e != nil {
		err = e
	} else if expect, e := imp_lines(k, m.MapOf("$LINES")); e != nil {
		err = e
	} else if e := ReadProg(&prog, reader.Unbox(r), cmds); e != nil {
		err = e
	} else if p, e := k.newProg("test", prog); e != nil {
		err = e
	} else {
		k.eph.NewTest(test, p, expect)
	}
	return
}

func imp_lines(k *Importer, r reader.Map) (ret string, err error) {
	return reader.String(r, "lines")
}
