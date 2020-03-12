package internal

import (
	"database/sql"
	"encoding/gob"
	"log"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/tables"

	_ "github.com/mattn/go-sqlite3"
)

type Parse func(*Parser, reader.Map) (err error)

// fix: queue could handle this. stack an event queue.
// maybe then we could stack "last noun" in there?
type CategoryEvent func(ephemera.Named)

type Parser struct {
	*ephemera.Recorder
	table      map[string]Parse
	oneTime    map[string]bool
	categories map[string]CategoryEvent // category
	nouns      reader.NameList
	currId     string
	currType   string
	parentItem reader.Map
}

func NewParser(srcURI string, db *sql.DB, fnds map[string]Parse) *Parser {
	registerGob()
	rec := ephemera.NewRecorder(srcURI, db)
	return &Parser{
		Recorder:   rec,
		table:      generators,
		oneTime:    make(map[string]bool),
		categories: make(map[string]CategoryEvent),
	}
}

var registeredGob = false

func registerGob() {
	if !registeredGob {
		for _, t := range export.Runs {
			gob.Register(t.Type)
		}
		registeredGob = true
	}
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
// named elements are considered unique within their category
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
func (r *Parser) parseItem(m reader.Map) (err error) {
	if len(m) > 0 {
		currId := m.StrOf(itemId)
		currType := m.StrOf(itemType)
		currValue := m.MapOf(itemValue)
		log.Println("parsing", currId, currType)
		if fn, ok := r.table[currType]; !ok {
			err = errutil.New("unknown type", currType)
		} else {
			// record into current context
			r.parentItem, r.currId, r.currType = m, currId, currType
			// call parsing function
			err = fn(r, currValue)
		}
	}
	return
}

func (r *Parser) parseSlice(ms []interface{}) (err error) {
	for _, it := range ms {
		if e := r.parseItem(reader.Box(it)); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func parseAttrs(r *Parser, item reader.Map) (err error) {
	defer r.on(tables.NAMED_TRAIT, func(trait ephemera.Named) {
		for _, noun := range r.nouns.Named {
			r.NewValue(noun, trait, true)
		}
	})()
	for _, it := range item.SliceOf("$ATTRIBUTE") {
		r.catStr(reader.Box(it), tables.NAMED_TRAIT)
	}
	return
}
