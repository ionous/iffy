package internal

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/export"

	_ "github.com/mattn/go-sqlite3"
)

type Parse func(*Importer, reader.Map) (err error)

// fix: queue could handle this. stack an event queue.
// maybe then we could stack "last noun" in there?
type CategoryEvent func(ephemera.Named)

// Importer helps read Map(s) of unmarshalled json.
type Importer struct {
	*ephemera.Recorder
	oneTime    map[string]bool
	categories map[string]CategoryEvent // category
	nouns      reader.NameList
}

func NewImporter(srcURI string, db *sql.DB) *Importer {
	registerGob()
	rec := ephemera.NewRecorder(srcURI, db)
	return &Importer{
		Recorder:   rec,
		oneTime:    make(map[string]bool),
		categories: make(map[string]CategoryEvent),
	}
}

var registeredGob = false

func registerGob() {
	if !registeredGob {
		export.Register(gob.Register)
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

// return { m[key] } as a new Named entry
// named elements are considered unique within their category
func (k *Importer) namedStr(m reader.Map, cat, key string) ephemera.Named {
	return k.catStr(m.MapOf(key), cat)
}

func (k *Importer) namedType(m reader.Map, typeName string) (ret ephemera.Named, err error) {
	if id, e := reader.Type(m, typeName); e != nil {
		err = e
	} else {
		ret = k.Named(typeName, m.StrOf(reader.ItemValue), id)
	}
	return
}

// helper to read a value of type string ( and interpret it as a Named value of category cat ).
// triggers a callback for .on() events of the passed category.
func (k *Importer) catStr(m reader.Map, cat string) ephemera.Named {
	id, str := m.StrOf(reader.ItemId), m.StrOf(reader.ItemValue)
	named := k.Named(cat, str, id)
	if h, ok := k.categories[cat]; ok {
		h(named)
	}
	return named
}

func (k *Importer) bind(cb Parse) reader.ReadMap {
	return func(m reader.Map) error {
		return cb(k, m)
	}
}

//
func (k *Importer) expectProg(r reader.Map, slotType string) (ret interface{}, err error) {
	if slot, ok := export.Slots[slotType]; !ok {
		err = errutil.New("unknown slot", slotType, reader.At(r))
	} else if m, e := reader.Slat(r, slotType); e != nil {
		err = e
	} else if v, e := ImportSlot(slot.Type, reader.Unbox(m), cmds); e != nil {
		err = e
	} else {
		ret = v
	}
	return
}

func (k *Importer) newProg(t string, i interface{}) (ret ephemera.Prog, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if e := enc.Encode(i); e != nil {
		err = e
	} else {
		ret = k.NewProg(t, buf.Bytes())
	}
	return
}
