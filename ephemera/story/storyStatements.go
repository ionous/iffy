package story

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
)

func imp_story_statement(k *Importer, r reader.Map) (err error) {
	return reader.Slot(r, "story_statement", reader.ReadMaps{
		"certainties":              k.Bind(imp_certainties),
		"aspect_traits":            k.Bind(imp_aspect_traits),
		"kinds_of_aspect":          k.Bind(imp_kinds_of_aspect),
		"kinds_of_kind":            k.Bind(imp_kinds_of_kind),
		"kinds_possess_properties": k.Bind(imp_kinds_possess_properties),
		"noun_assignment":          k.Bind(imp_noun_assignment),
		"noun_statement":           k.Bind(imp_noun_statement),
		"pattern_decl":             k.Bind(imp_pattern_decl),
		"pattern_variables_decl":   k.Bind(imp_pattern_variables_decl),
		"relative_to_noun":         k.Bind(imp_relative_to_noun),
		"test_statement":           k.Bind(imp_test_statement),
		"test_scene":               k.Bind(imp_test_scene),
		"pattern_handler":          k.Bind(imp_pattern_handler),
		"pattern_actions":          k.Bind(imp_pattern_actions),
	})
}

//"{plural_kinds} {are_being} {certainty} {trait}.");
// horses are usually fast.
func imp_certainties(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "certainties"); e != nil {
		err = e
	} else if certainty, e := imp_certainty(k, m.MapOf("$CERTAINTY")); e != nil {
		err = e
	} else if trait, e := imp_trait(k, m.MapOf("$TRAIT")); e != nil {
		err = e
	} else if kind, e := imp_plural_kinds(k, m.MapOf("$PLURAL_KINDS")); e != nil {
		err = e
	} else {
		k.NewCertainty(certainty, trait, kind)
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

// "{aspects} {are_an} kind of value."
// ex. colors are a kind of value
func imp_kinds_of_aspect(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "kinds_of_aspect"); e != nil {
		err = e
	} else if _, e := imp_aspect(k, m.MapOf("$ASPECT")); e != nil {
		err = e
	}
	return
}

// "{plural_kinds} {are_an} kind of {kind}."
// ex. "cats are a kind of animal"
func imp_kinds_of_kind(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "kinds_of_kind"); e != nil {
		err = e
	} else if kind, e := imp_plural_kinds(k, m.MapOf("$PLURAL_KINDS")); e != nil {
		err = e
	} else if parent, e := imp_singular_kind(k, m.MapOf("$SINGULAR_KIND")); e != nil {
		err = e
	} else {
		k.NewKind(kind, parent)
	}
	return
}

// {plural_kinds} have {determiner} {property_phrase}.
// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func imp_kinds_possess_properties(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "kinds_possess_properties"); e != nil {
		err = e
	} else if kind, e := imp_plural_kinds(k, m.MapOf("$PLURAL_KINDS")); e != nil {
		err = e
	} else /*if _, e := imp_determiner(k, m.MapOf("$DETERMINER")); e != nil {
		err = e
	} else */if e := imp_property_phrase(k, kind, m.MapOf("$PROPERTY_PHRASE")); e != nil {
		err = e
	}
	return
}

// run "The {property} of {+noun} is the {[text]:: %lines}"
// ex. The description of the nets is xxx
func imp_noun_assignment(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "noun_assignment"); e != nil {
		err = e
	} else if lines, e := imp_line_expr(k, m.MapOf("$LINES")); e != nil {
		err = e
	} else if prop, e := imp_property(k, m.MapOf("$PROPERTY")); e != nil {
		err = e
	} else if e := k.Recent.Nouns.CollectSubjects(func() error {
		return reader.Repeats(m.SliceOf("$NOUN"), k.Bind(imp_noun))
	}); e != nil {
		err = e
	} else {
		for _, noun := range k.Recent.Nouns.Subjects {
			k.NewValue(noun, prop, lines)
		}
	}
	return
}

func imp_noun_statement(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "noun_statement"); e != nil {
		err = e
	} else if e := imp_lede(k, m.MapOf("$LEDE")); e != nil {
		err = e
	} else if e := reader.Repeats(m.SliceOf("$TAIL"), k.Bind(imp_tail)); e != nil {
		err = e
	} else if v := m.MapOf("$SUMMARY"); len(v) != 0 {
		err = imp_summary(k, v)
	}
	return
}

// Adds a new pattern declaration and optionally some associated pattern parameters.
// {name:pattern_name|quote} determines {type:pattern_type}.
// {optvars?pattern_variables_tail}
func imp_pattern_decl(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "pattern_decl"); e != nil {
		err = e
	} else if patternName, e := imp_pattern_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else if patternType, e := imp_pattern_type(k, m.MapOf("$TYPE")); e != nil {
		err = e
	} else {
		k.NewPatternDecl(patternName, patternName, patternType)
		if tail := m.MapOf("$OPTVARS"); len(tail) > 0 {
			err = imp_pattern_variables_tail(k, patternName, tail)
		}
	}
	return
}

// `Pattern variables: Storage for values used during the execution of a pattern.`)
// {+variable_decl|comma-and}.",
func imp_pattern_variables_tail(k *Importer, patternName ephemera.Named, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "pattern_variables_tail"); e != nil {
		err = e
	} else {
		err = rep_variable_decl(k, patternName, m)
	}
	return
}

func rep_variable_decl(k *Importer, patternName ephemera.Named, r reader.Map) error {
	return reader.Repeats(r.SliceOf("$VARIABLE_DECL"),
		func(m reader.Map) (err error) {
			if paramName, paramType, e := imp_variable_decl(k, m); e != nil {
				err = e
			} else {
				k.NewPatternDecl(patternName, paramName, paramType)
			}
			return
		})
}

// "The pattern {pattern_name|quote} uses {+variable_decl|comma-and}."
func imp_pattern_variables_decl(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "pattern_variables_decl"); e != nil {
		err = e
	} else if patternName, e := imp_pattern_name(k, m.MapOf("$PATTERN_NAME")); e != nil {
		err = e
	} else {
		err = rep_variable_decl(k, patternName, m)
	}
	return
}

// "{aspects} {trait_phrase}"
// (the) colors are red, blue, or green.
func imp_aspect_traits(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "aspect_traits"); e != nil {
		err = e
	} else if aspect, e := imp_aspect(k, r.MapOf("$ASPECT")); e != nil {
		err = e
	} else if traits, e := imp_trait_phrase(k, m.MapOf("$TRAIT_PHRASE")); e != nil {
		err = e
	} else {
		for rank, trait := range traits {
			k.NewTrait(trait, aspect, rank)
		}
	}
	return
}

// "{relation} {+noun} {are_being} {+noun}."
// ex. On the beach are shells.
func imp_relative_to_noun(k *Importer, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "relative_to_noun"); e != nil {
		err = e
	} else if relation, e := imp_relation(k, m.MapOf("$RELATION")); e != nil {
		err = e
	} else if e := k.Recent.Nouns.CollectObjects(func() error {
		return reader.Repeats(m.SliceOf("$NOUN"), k.Bind(imp_noun))
	}); e != nil {
		err = e
	} else if e := k.Recent.Nouns.CollectSubjects(func() error {
		return reader.Repeats(m.SliceOf("$NOUN1"), k.Bind(imp_noun))
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
