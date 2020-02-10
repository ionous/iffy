package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os/user"
	"path"

	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/kr/pretty"
	_ "github.com/mattn/go-sqlite3"
)

type Parse func(*Parser, reader.Map)

// fix: queue could handle this. stack an event queue.
// maybe then we could stack "last noun" in there?
type CategoryEvent func(ephemera.Named)

type Parser struct {
	*ephemera.Recorder
	table      map[string]Parse
	oneTime    map[string]bool
	categories map[string]CategoryEvent // category
	nouns      reader.NameList
	lastId     string
	lastType   string
	lastValue  interface{}
}

// return deferred removal
func (r *Parser) on(cat string, handler CategoryEvent) func() {
	was := r.categories[cat]
	r.categories[cat] = handler
	return func() {
		if was != nil {
			r.categories[cat] = was
		} else {
			delete(r.categories, cat)
		}
	}
}

// return true if item is the first time once has been called with the specified string.
func (r *Parser) once(s string) (ret bool) {
	if !r.oneTime[s] {
		r.oneTime[s] = true
		ret = true
	}
	return
}

// return m[key]["value" as a string
func (*Parser) getStr(m reader.Map, param string) string {
	return m.MapOf(param).StrOf(itemValue)
}

// return { m[key] } as a new Named entry
func (r *Parser) namedStr(m reader.Map, cat, key string) ephemera.Named {
	return r.catStr(m.MapOf(key), cat)
}

func (r *Parser) catStr(item reader.Map, cat string) ephemera.Named {
	id, str := item.StrOf(itemId), item.StrOf(itemValue)
	named := r.Named(cat, str, id)
	if h, ok := r.categories[cat]; ok {
		h(named)
	}
	return named
}

// helper to process m["type"]
func (r *Parser) parse(m reader.Map) {
	if len(m) > 0 {
		currId := m.StrOf(itemId)
		currType := m.StrOf(itemType)
		currValue := m.MapOf(itemValue)

		log.Println("parsing", currId, currType)
		if fn, ok := r.table[currType]; !ok {
			log.Fatalln("unknown type", currType)
		} else {
			r.lastId, r.lastType, r.lastValue = currId, currType, currValue
			fn(r, currValue)
		}
	}
}

func (r *Parser) parseSlice(ms []interface{}) {
	for _, it := range ms {
		r.parse(reader.Cast(it))
	}
}

// -----------------------------------
const (
	itemId    = "id"
	itemType  = "type"
	itemValue = "value"
)

func parseAttrs(r *Parser, item reader.Map) {
	defer r.on(ephemera.NAMED_TRAIT, func(trait ephemera.Named) {
		for _, noun := range r.nouns.Named {
			r.NewValue(noun, trait, true)
		}
	})()
	for _, it := range item.SliceOf("$ATTRIBUTE") {
		r.catStr(reader.Cast(it), ephemera.NAMED_TRAIT)
	}
}

var fns = map[string]Parse{
	// story is a bunch of paragraphs
	"story": func(r *Parser, item reader.Map) {
		r.parseSlice(item.SliceOf("$PARAGRAPH"))
	},

	// paragraph is a bunch of statements
	"paragraph": func(r *Parser, item reader.Map) {
		r.parseSlice(item.SliceOf("$STORY_STATEMENT"))
	},

	"story_statement": func(r *Parser, item reader.Map) {
		r.nouns.Swap(nil)
		r.parse(item)
	},

	"noun_statement": func(r *Parser, item reader.Map) {
		r.parse(item.MapOf("$LEDE"))
		r.parseSlice(item.SliceOf("$TAIL"))
		r.parse(item.MapOf("$SUMMARY"))
	},

	// run "The {property} of {+noun} is the {[text]:: %lines}"
	// ex. The description of the nets is xxx
	"noun_assignment": func(r *Parser, item reader.Map) {
		val := r.namedStr(item, ephemera.PRIM_EXPR, "$LINES")
		prop := r.namedStr(item, ephemera.NAMED_FIELD, "$PROPERTY")
		defer r.on(ephemera.NAMED_NOUN, func(noun ephemera.Named) {
			r.NewValue(noun, prop, val)
		})()
		r.parseSlice(item.SliceOf("$NOUN"))
	},

	// "{relation} {+noun} {are_being} {+noun}."
	// ex. On the beach are shells.
	"relative_to_noun": func(r *Parser, item reader.Map) {
		relation := r.namedStr(item, ephemera.NAMED_VERB, "$RELATION")
		//
		r.parseSlice(item.SliceOf("$NOUN"))
		leadingNouns := r.nouns.Swap(nil)
		r.parseSlice(item.SliceOf("$NOUN1"))
		trailingNouns := r.nouns.Swap(leadingNouns)
		//
		for _, a := range leadingNouns {
			for _, b := range trailingNouns {
				r.NewRelative(a, relation, b)
			}
		}
	},

	// "{plural_kinds} {are_an} kind of {kind}."
	// ex. "cats are a kind of animal"
	"kinds_of_thing": func(r *Parser, item reader.Map) {
		kind := r.namedStr(item, ephemera.NAMED_KIND, "$PLURAL_KINDS")
		parent := r.namedStr(item, ephemera.NAMED_KIND, "$KIND")
		r.NewKind(kind, parent)
	},

	// "{qualities} {are_an} kind of value."
	// ex. colors are a kind of value
	"kinds_of_quality": func(r *Parser, item reader.Map) {
		aspect := r.namedStr(item, ephemera.NAMED_ASPECT, "$QUALITY")
		r.NewAspect(aspect)
	},

	// "{plural_kinds} {attribute_phrase}"
	// ex. animals are fast or slow.
	"class_attributes": func(r *Parser, item reader.Map) {
		kind := r.namedStr(item, ephemera.NAMED_KIND, "$PLURAL_KINDS")
		//
		var traits []ephemera.Named
		defer r.on(ephemera.NAMED_TRAIT, func(trait ephemera.Named) {
			traits = append(traits, trait)
		})()
		r.parse(item.MapOf("$ATTRIBUTE_PHRASE"))
		// create an implied aspect named after the first trait
		// fix? maybe we should include the columns of named in the returned struct so we can pick out the source better.
		aspect := r.Named(ephemera.NAMED_ASPECT, traits[0].String(), item.StrOf(itemId))
		r.NewPrimitive(ephemera.PRIM_ASPECT, kind, aspect)
	},

	// "{qualities} {attribute_phrase}"
	// (the) colors are red, blue, or green.
	"quality_attributes": func(r *Parser, item reader.Map) {
		aspect := r.namedStr(item, ephemera.NAMED_ASPECT, "$QUALITIES")
		rank := 0
		defer r.on(ephemera.NAMED_TRAIT, func(trait ephemera.Named) {
			r.NewTrait(trait, aspect, rank)
			rank += 1
		})()
		r.parse(item.MapOf("$ATTRIBUTE_PHRASE"))
	},

	//"{plural_kinds} {are_being} {certainty} {attribute}.");
	// horses are usually fast.
	"certainties": func(r *Parser, item reader.Map) {
		certainty := r.getStr(item, "$CERTAINTY")
		trait := r.namedStr(item, ephemera.NAMED_TRAIT, "$ATTRIBUTE")
		kind := r.namedStr(item, ephemera.NAMED_KIND, "$PLURAL_KINDS")
		r.NewCertainty(certainty, trait, kind)
	},

	// {plural_kinds} have {determiner} {primitive_type} called {property}.
	// ex. cats have some text called breed.
	// ex. horses have an aspect called speed.
	"kinds_possess_properties": func(r *Parser, item reader.Map) {
		kind := r.namedStr(item, ephemera.NAMED_KIND, "$PLURAL_KINDS")
		prop := r.namedStr(item, ephemera.NAMED_KIND, "$PROPERTY")
		primType := r.getStr(item, "$PRIMITIVE_TYPE")
		r.NewPrimitive(primType, kind, prop)
	},

	// run: "{+names} {noun_phrase}."
	"lede": func(r *Parser, item reader.Map) {
		r.parseSlice(item.SliceOf("$NOUN"))
		r.parse(item.MapOf("$NOUN_PHRASE"))
	},

	// run: "{pronoun} {noun_phrase}."
	"tail": func(r *Parser, item reader.Map) {
		r.parse(item.MapOf("$PRONOUN"))
		r.parse(item.MapOf("$NOUN_PHRASE"))
	},

	// opt: "{kind_of_noun}, {noun_attrs}, or {noun_relation}"
	"noun_phrase": func(r *Parser, item reader.Map) {
		r.parse(item)
	},

	// run: "{?are_being} {relation} {+noun}"
	// ex. (the cat and the hat) are in (the book)
	// ex. (Hector and Maria) are suspicious of (Santa and Santana).
	"noun_relation": func(r *Parser, item reader.Map) {
		relation := r.namedStr(item, ephemera.NAMED_VERB, "$RELATION")
		//
		leadingNouns := r.nouns.Swap(nil)
		r.parseSlice(item.SliceOf("$NOUN"))
		trailingNouns := r.nouns.Swap(leadingNouns)
		//
		for _, n := range leadingNouns {
			for _, d := range trailingNouns {
				r.NewRelative(n, relation, d)
			}
		}
	},

	// noun: "{proper_noun} or {common_noun}"
	"noun": func(r *Parser, item reader.Map) {
		once := "noun"
		if r.once(once) {
			things := r.Named(ephemera.NAMED_KIND, "things", once)
			nounType := r.Named(ephemera.NAMED_ASPECT, "nounType", once)
			common := r.Named(ephemera.NAMED_TRAIT, "common", once)
			proper := r.Named(ephemera.NAMED_TRAIT, "proper", once)
			r.NewPrimitive(ephemera.PRIM_ASPECT, things, nounType)
			r.NewTrait(common, nounType, 0)
			r.NewTrait(proper, nounType, 1)
		}
		r.parse(item)
	},

	// run: "{determiner} {common_name}"
	"common_noun": func(r *Parser, item reader.Map) {
		id := r.lastId
		det := r.getStr(item, "$DETERMINER")
		noun := r.namedStr(item, ephemera.NAMED_NOUN, "$COMMON_NAME")
		r.nouns.Add(noun)
		// set common nounType to true ( implicitly defined by "noun" )
		nounType := r.Named(ephemera.NAMED_TRAIT, "common", id)
		r.NewValue(noun, nounType, true)
		//
		if det[0] != '$' {
			article := r.Named(ephemera.NAMED_FIELD, "indefinite article", id)
			r.NewValue(noun, article, det)
			once := "common_noun"
			if r.once(once) {
				indefinite := r.Named(ephemera.NAMED_FIELD, "indefinite article", once)
				things := r.Named(ephemera.NAMED_KIND, "things", once)
				r.NewPrimitive(ephemera.PRIM_TEXT, things, indefinite)
			}
		}
	},

	// run: "{proper_name}"
	// common / proper setting
	"proper_noun": func(r *Parser, item reader.Map) {
		id := r.lastId
		noun := r.namedStr(item, ephemera.NAMED_NOUN, "$PROPER_NAME")
		r.nouns.Add(noun)
		// set proper nounType to true ( implicitly defined by "noun" )
		nounType := r.Named(ephemera.NAMED_TRAIT, "proper", id)
		r.NewValue(noun, nounType, true)
	},

	// run: "{are_an} {*attribute} {kind} {?noun_relation}"
	// ex. "(the box) is a closed container on the beach"
	"kind_of_noun": func(r *Parser, item reader.Map) {
		//
		kind := r.namedStr(item, ephemera.NAMED_KIND, "$KIND")
		for _, noun := range r.nouns.Named {
			r.NewNoun(noun, kind)
		}
		parseAttrs(r, item)
		// noun relation takes care of itself --
		// relating the new nouns to the existing nouns.
		r.parse(item.MapOf("$NOUN_RELATION"))
	},

	// run: "{are_being} {+attribute}"
	// ex. "(the box) is closed"
	"noun_attrs": parseAttrs,

	// run: "{are_either} {+attribute}."
	"attribute_phrase": parseAttrs,

	// run: "{The [summary] is:: %lines}"
	"summary": func(r *Parser, item reader.Map) {
		once := "summary"
		id := r.lastId
		if r.once(once) {
			things := r.Named(ephemera.NAMED_KIND, "things", once)
			appear := r.Named(ephemera.NAMED_FIELD, "appearance", once)
			r.NewPrimitive(ephemera.PRIM_EXPR, things, appear)
		}
		prop := r.Named(ephemera.NAMED_FIELD, "appearance", id)
		noun := r.nouns.Last()
		val := r.namedStr(item, ephemera.PRIM_EXPR, "$LINES")
		r.NewValue(noun, prop, val)
	},
}

// 	flag.Parse() for processing command line args
func main() {
	// const memory = "file:test.db?cache=shared&mode=memory"
	if user, e := user.Current(); e != nil {
		log.Fatalln(e)
	} else {
		fileName := path.Join(user.HomeDir, "iffyTest.db")
		if db, e := sql.Open("sqlite3", fileName); e != nil {
			log.Fatalln("db open", e)
		} else {
			var top reader.Map
			if e := json.Unmarshal([]byte(debug.Blob), &top); e != nil {
				log.Fatalln("read json", e)
			}
			genq := &ephemera.GenQueue{}
			dbq := ephemera.NewDBQueue(db)
			q := ephemera.NewStack(genq, dbq)
			rec := ephemera.NewRecorder("blob", q)
			r := Parser{Recorder: rec,
				table:      fns,
				oneTime:    make(map[string]bool),
				categories: make(map[string]CategoryEvent),
			}

			top.Expect(itemType, "story")
			r.parse(top)
			if b, e := json.MarshalIndent(genq.Tables, "", "  "); e != nil {
				log.Fatalln(e)
			} else {
				pretty.Println("tables", string(b))
			}

			defer db.Close()
		}
	}
}
