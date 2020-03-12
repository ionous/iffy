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

var generators = map[string]Parse{
	"test": func(r *Parser, item reader.Map) (err error) {
		test := r.namedStr(item, tables.NAMED_TEST, "$TEST_NAME")
		expect := r.getStr(item, "$LINES")

		var prog check.Test // parentItem is { id:..., type:"test", value:... }
		if e := readProg(&prog, reader.Unbox(r.parentItem), export.Runs); e != nil {
			err = e
		} else {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			if e := enc.Encode(prog); e != nil {
				err = e
			} else {
				prog := r.NewProg(r.currType, buf.Bytes())
				r.NewTest(test, prog, expect)
			}
		}
		return
	},
	// story is a bunch of paragraphs
	"story": func(r *Parser, item reader.Map) error {
		return r.parseSlice(item.SliceOf("$PARAGRAPH"))
	},

	// paragraph is a bunch of statements
	"paragraph": func(r *Parser, item reader.Map) error {
		return r.parseSlice(item.SliceOf("$STORY_STATEMENT"))
	},

	"story_statement": func(r *Parser, item reader.Map) error {
		r.nouns.Swap(nil)
		return r.parseItem(item)
	},

	"noun_statement": func(r *Parser, item reader.Map) error {
		return errutil.Append(
			r.parseItem(item.MapOf("$LEDE")),
			r.parseSlice(item.SliceOf("$TAIL")),
			r.parseItem(item.MapOf("$SUMMARY")))
	},

	// run "The {property} of {+noun} is the {[text]:: %lines}"
	// ex. The description of the nets is xxx
	"noun_assignment": func(r *Parser, item reader.Map) error {
		val := r.namedStr(item, tables.PRIM_EXPR, "$LINES")
		prop := r.namedStr(item, tables.NAMED_FIELD, "$PROPERTY")
		defer r.on(tables.NAMED_NOUN, func(noun ephemera.Named) {
			r.NewValue(noun, prop, val)
		})()
		return r.parseSlice(item.SliceOf("$NOUN"))
	},

	// "{relation} {+noun} {are_being} {+noun}."
	// ex. On the beach are shells.
	"relative_to_noun": func(r *Parser, item reader.Map) (err error) {
		relation := r.namedStr(item, tables.NAMED_VERB, "$RELATION")
		//
		if e := r.parseSlice(item.SliceOf("$NOUN")); e != nil {
			err = errutil.Append(err, e)
		}
		leadingNouns := r.nouns.Swap(nil)
		if e := r.parseSlice(item.SliceOf("$NOUN1")); e != nil {
			err = errutil.Append(err, e)
		}
		trailingNouns := r.nouns.Swap(leadingNouns)
		//
		for _, a := range leadingNouns {
			for _, b := range trailingNouns {
				r.NewRelative(a, relation, b)
			}
		}
		return err
	},

	// "{plural_kinds} {are_an} kind of {kind}."
	// ex. "cats are a kind of animal"
	"kinds_of_thing": func(r *Parser, item reader.Map) error {
		kind := r.namedStr(item, tables.NAMED_KIND, "$PLURAL_KINDS")
		parent := r.namedStr(item, tables.NAMED_KIND, "$KIND")
		r.NewKind(kind, parent)
		return nil
	},

	// "{qualities} {are_an} kind of value."
	// ex. colors are a kind of value
	"kinds_of_quality": func(r *Parser, item reader.Map) error {
		aspect := r.namedStr(item, tables.NAMED_ASPECT, "$QUALITY")
		r.NewAspect(aspect)
		return nil
	},

	// "{plural_kinds} {attribute_phrase}"
	// ex. animals are fast or slow.
	"class_attributes": func(r *Parser, item reader.Map) (err error) {
		kind := r.namedStr(item, tables.NAMED_KIND, "$PLURAL_KINDS")
		//
		var traits []ephemera.Named
		defer r.on(tables.NAMED_TRAIT, func(trait ephemera.Named) {
			traits = append(traits, trait)
		})()
		err = r.parseItem(item.MapOf("$ATTRIBUTE_PHRASE"))
		// create an implied aspect named after the first trait
		// fix? maybe we should include the columns of named in the returned struct so we can pick out the source better.
		aspect := r.Named(tables.NAMED_ASPECT, traits[0].String(), item.StrOf(itemId))
		r.NewPrimitive(tables.PRIM_ASPECT, kind, aspect)
		return
	},

	// "{qualities} {attribute_phrase}"
	// (the) colors are red, blue, or green.
	"quality_attributes": func(r *Parser, item reader.Map) error {
		aspect := r.namedStr(item, tables.NAMED_ASPECT, "$QUALITIES")
		rank := 0
		defer r.on(tables.NAMED_TRAIT, func(trait ephemera.Named) {
			r.NewTrait(trait, aspect, rank)
			rank += 1
		})()
		return r.parseItem(item.MapOf("$ATTRIBUTE_PHRASE"))
	},

	//"{plural_kinds} {are_being} {certainty} {attribute}.");
	// horses are usually fast.
	"certainties": func(r *Parser, item reader.Map) error {
		certainty := r.getStr(item, "$CERTAINTY")
		trait := r.namedStr(item, tables.NAMED_TRAIT, "$ATTRIBUTE")
		kind := r.namedStr(item, tables.NAMED_KIND, "$PLURAL_KINDS")
		r.NewCertainty(certainty, trait, kind)
		return nil
	},

	// {plural_kinds} have {determiner} {primitive_type} called {property}.
	// ex. cats have some text called breed.
	// ex. horses have an aspect called speed.
	"kinds_possess_properties": func(r *Parser, item reader.Map) error {
		kind := r.namedStr(item, tables.NAMED_KIND, "$PLURAL_KINDS")
		prop := r.namedStr(item, tables.NAMED_KIND, "$PROPERTY")
		primType := r.getStr(item, "$PRIMITIVE_TYPE")
		r.NewPrimitive(primType, kind, prop)
		return nil
	},

	// run: "{+names} {noun_phrase}."
	"lede": func(r *Parser, item reader.Map) error {
		return errutil.Append(
			r.parseSlice(item.SliceOf("$NOUN")),
			r.parseItem(item.MapOf("$NOUN_PHRASE")))
	},

	// run: "{pronoun} {noun_phrase}."
	"tail": func(r *Parser, item reader.Map) error {
		return errutil.Append(
			r.parseItem(item.MapOf("$PRONOUN")),
			r.parseItem(item.MapOf("$NOUN_PHRASE")))
	},

	// opt: "{kind_of_noun}, {noun_attrs}, or {noun_relation}"
	"noun_phrase": func(r *Parser, item reader.Map) error {
		return r.parseItem(item)
	},

	// run: "{?are_being} {relation} {+noun}"
	// ex. (the cat and the hat) are in (the book)
	// ex. (Hector and Maria) are suspicious of (Santa and Santana).
	"noun_relation": func(r *Parser, item reader.Map) (err error) {
		relation := r.namedStr(item, tables.NAMED_VERB, "$RELATION")
		//
		leadingNouns := r.nouns.Swap(nil)
		err = r.parseSlice(item.SliceOf("$NOUN"))
		trailingNouns := r.nouns.Swap(leadingNouns)
		//
		for _, n := range leadingNouns {
			for _, d := range trailingNouns {
				r.NewRelative(n, relation, d)
			}
		}
		return
	},

	// noun: "{proper_noun} or {common_noun}"
	"noun": func(r *Parser, item reader.Map) error {
		once := "noun"
		if r.once(once) {
			things := r.Named(tables.NAMED_KIND, "things", once)
			nounType := r.Named(tables.NAMED_ASPECT, "nounType", once)
			common := r.Named(tables.NAMED_TRAIT, "common", once)
			proper := r.Named(tables.NAMED_TRAIT, "proper", once)
			r.NewPrimitive(tables.PRIM_ASPECT, things, nounType)
			r.NewTrait(common, nounType, 0)
			r.NewTrait(proper, nounType, 1)
		}
		return r.parseItem(item)
	},

	// run: "{determiner} {common_name}"
	"common_noun": func(r *Parser, item reader.Map) error {
		id := r.currId
		det := r.getStr(item, "$DETERMINER")
		noun := r.namedStr(item, tables.NAMED_NOUN, "$COMMON_NAME")
		r.nouns.Add(noun)
		// set common nounType to true ( implicitly defined by "noun" )
		nounType := r.Named(tables.NAMED_TRAIT, "common", id)
		r.NewValue(noun, nounType, true)
		//
		if det[0] != '$' {
			article := r.Named(tables.NAMED_FIELD, "indefinite article", id)
			r.NewValue(noun, article, det)
			once := "common_noun"
			if r.once(once) {
				indefinite := r.Named(tables.NAMED_FIELD, "indefinite article", once)
				things := r.Named(tables.NAMED_KIND, "things", once)
				r.NewPrimitive(tables.PRIM_TEXT, things, indefinite)
			}
		}
		return nil
	},

	// run: "{proper_name}"
	// common / proper setting
	"proper_noun": func(r *Parser, item reader.Map) error {
		id := r.currId
		noun := r.namedStr(item, tables.NAMED_NOUN, "$PROPER_NAME")
		r.nouns.Add(noun)
		// set proper nounType to true ( implicitly defined by "noun" )
		nounType := r.Named(tables.NAMED_TRAIT, "proper", id)
		r.NewValue(noun, nounType, true)
		return nil
	},

	// run: "{are_an} {*attribute} {kind} {?noun_relation}"
	// ex. "(the box) is a closed container on the beach"
	"kind_of_noun": func(r *Parser, item reader.Map) error {
		//
		kind := r.namedStr(item, tables.NAMED_KIND, "$KIND")
		for _, noun := range r.nouns.Named {
			r.NewNoun(noun, kind)
		}
		parseAttrs(r, item)
		// noun relation takes care of itself --
		// relating the new nouns to the existing nouns.
		return r.parseItem(item.MapOf("$NOUN_RELATION"))
	},

	// run: "{are_being} {+attribute}"
	// ex. "(the box) is closed"
	"noun_attrs": parseAttrs,

	// run: "{are_either} {+attribute}."
	"attribute_phrase": parseAttrs,

	// run: "{The [summary] is:: %lines}"
	"summary": func(r *Parser, item reader.Map) error {
		once := "summary"
		id := r.currId
		if r.once(once) {
			things := r.Named(tables.NAMED_KIND, "things", once)
			appear := r.Named(tables.NAMED_FIELD, "appearance", once)
			r.NewPrimitive(tables.PRIM_EXPR, things, appear)
		}
		prop := r.Named(tables.NAMED_FIELD, "appearance", id)
		noun := r.nouns.Last()
		val := r.namedStr(item, tables.PRIM_EXPR, "$LINES")
		r.NewValue(noun, prop, val)
		return nil
	},
}
