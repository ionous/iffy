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

type Parse func(*Importer, reader.Map) (err error)

// fix: queue could handle this. stack an event queue.
// maybe then we could stack "last noun" in there?
type CategoryEvent func(ephemera.Named)

// Importer helps read Map(s) of unmarshalled json.
// Its generators are *not* called automatically;
// instead individual functions may choose to call parseItem on specific values,
// and so on recursively.
type Importer struct {
	*ephemera.Recorder
	namedGen   map[string]Parse
	oneTime    map[string]bool
	categories map[string]CategoryEvent // category
	nouns      reader.NameList
	currId     string
	currType   string
	parentItem reader.Map
}

func NewImporter(srcURI string, db *sql.DB, namedGen map[string]Parse) *Importer {
	registerGob()
	rec := ephemera.NewRecorder(srcURI, db)
	return &Importer{
		Recorder:   rec,
		namedGen:   namedGen,
		oneTime:    make(map[string]bool),
		categories: make(map[string]CategoryEvent),
	}
}

var registeredGob = false

func registerGob() {
	if !registeredGob {
		for _, cmd := range export.Runs {
			gob.Register(cmd)
		}
		registeredGob = true
	}
}

// listen for "events" of parsing the passed category.
// return deferred removal
func (k *Importer) on(cat string, handler CategoryEvent) func() {
	was := k.categories[cat]
	k.categories[cat] = handler
	return func() {
		if was != nil {
			k.categories[cat] = was
		} else {
			delete(k.categories, cat)
		}
	}
}

// return true if m is the first time once has been called with the specified string.
func (k *Importer) once(s string) (ret bool) {
	if !k.oneTime[s] {
		k.oneTime[s] = true
		ret = true
	}
	return
}

// return m[key]["value"] as a string
func (*Importer) getStr(m reader.Map, param string) string {
	return m.MapOf(param).StrOf(itemValue)
}

// return { m[key] } as a new Named entry
// named elements are considered unique within their category
func (k *Importer) namedStr(m reader.Map, cat, key string) ephemera.Named {
	return k.catStr(m.MapOf(key), cat)
}

func (k *Importer) namedType(m reader.Map, typeName string) (ret ephemera.Named, err error) {
	if e := k.expectedType(m, typeName); e != nil {
		err = e
	} else {
		id, str := m.StrOf(itemId), m.StrOf(itemValue)
		ret = k.Named(typeName, str, id)
	}
	return
}

// helper to read a value of type string ( and interpret it as a Named value of category cat ).
// triggers a callback for .on() events of the passed category.
func (k *Importer) catStr(m reader.Map, cat string) ephemera.Named {
	id, str := m.StrOf(itemId), m.StrOf(itemValue)
	named := k.Named(cat, str, id)
	if h, ok := k.categories[cat]; ok {
		h(named)
	}
	return named
}

// helper to process m["type"]
func (k *Importer) parseItem(m reader.Map) (err error) {
	if len(m) > 0 {
		currId := m.StrOf(itemId)
		currType := m.StrOf(itemType)
		currValue := m.MapOf(itemValue)
		log.Println("parseItem", currId, currType)
		if fn, ok := k.namedGen[currType]; !ok {
			err = wrongType("", currType, currId)
		} else {
			// record into current context
			k.parentItem, k.currId, k.currType = m, currId, currType
			// call parsing function
			err = fn(k, currValue)
		}
	}
	return
}

func (k *Importer) parseSlice(ms []interface{}) (err error) {
	for _, it := range ms {
		if e := k.parseItem(reader.Box(it)); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (k *Importer) repeats(ms []interface{}, cb func(reader.Map) error) (err error) {
	for _, it := range ms {
		if e := cb(reader.Box(it)); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// we expect to see one, and only one, of the sub keys in the itemValue of m.
func (k *Importer) expectOpt(m reader.Map, expectedType string, sub map[string]Parse) (err error) {
	if currType := m.StrOf(itemType); currType != expectedType {
		err = wrongType(expectedType, currType, at(m))
	} else if currValue := m.MapOf(itemValue); len(currValue) != 1 {
		err = wrongType(expectedType, currType, at(m))
	} else {
		// only one in the list.
		for i, v := range currValue {
			if fn, ok := sub[i]; !ok {
				err = wrongValue(expectedType, k, at(m))
			} else {
				err = fn(k, reader.Box(v))
			}
			break
		}
	}
	return
}

// expect a string constant
func (k *Importer) expectConst(m reader.Map, expectedType, expectedValue string) (err error) {
	if currType := m.StrOf(itemType); currType != expectedType {
		err = wrongType(expectedType, currType, at(m))
	} else if currValue := m.StrOf(itemValue); currValue != expectedValue {
		err = wrongValue(currType, currValue, at(m))
	}
	return
}

// expect a map value
func (k *Importer) expectSlat(m reader.Map, expectedType string) (ret reader.Map, err error) {
	if e := k.expectedType(m, expectedType); e != nil {
		err = e
	} else {
		ret = m.MapOf(itemValue)
	}
	return
}

// expect a string variable
func (k *Importer) expectStr(m reader.Map, expectedType string) (ret string, err error) {
	if currType := m.StrOf(itemType); currType != expectedType {
		err = wrongType(expectedType, currType, at(m))
	} else if currValue := m.StrOf(itemValue); len(currValue) == 0 {
		err = wrongValue(currType, currValue, at(m))
	} else {
		ret = currValue
	}
	return
}

// expect a string variable
func (k *Importer) expectName(m reader.Map, expectedType, cat string) (err error) {
	if currType := m.StrOf(itemType); currType != expectedType {
		err = wrongType(expectedType, currType, at(m))
	} else if currValue := m.StrOf(itemValue); len(currValue) == 0 {
		err = wrongValue(currType, currValue, at(m))
	}
	return
}

// helper: check the type of the passed m map
func (k *Importer) expectedType(m reader.Map, expectedType string) (err error) {
	if currType := m.StrOf(itemType); currType != expectedType {
		err = wrongType(expectedType, currType, at(m))
	}
	return
}

func parseAttrs(k *Importer, m reader.Map) (err error) {
	defer k.on(tables.NAMED_TRAIT, func(trait ephemera.Named) {
		for _, noun := range k.nouns.Named {
			k.NewValue(noun, trait, true)
		}
	})()
	for _, it := range m.SliceOf("$ATTRIBUTE") {
		k.catStr(reader.Box(it), tables.NAMED_TRAIT)
	}
	return
}

func at(m reader.Map) string {
	return m.StrOf(itemId)
}

func wrongType(wanted, got, at string) error {
	return errutil.New("unexpected type", got, "wanted", wanted, "at", at)
}

func wrongValue(
	currType string,
	got interface{},
	at string,
) error {
	return errutil.New(currType, "has unexpected value", got, "at", at)
}
