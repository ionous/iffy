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
	if id, e := k.expectedType(m, typeName); e != nil {
		err = e
	} else {
		ret = k.Named(typeName, m.StrOf(itemValue), id)
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

func (k *Importer) repeats(ms []interface{}, cb Parse) (err error) {
	for _, it := range ms {
		if e := cb(k, reader.Box(it)); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// we expect to see one, and only one, of the sub keys in the itemValue of m.
func (k *Importer) expectOpt(r reader.Map, expectedType string, sub map[string]Parse) (err error) {
	if t := r.StrOf(itemType); t != expectedType {
		err = wrongType(expectedType, t, at(r))
	} else if m := r.MapOf(itemValue); len(m) != 1 {
		err = wrongValue(t, m, at(r))
	} else {
		// only one in the list.
		for key, value := range m {
			if fn, ok := sub[key]; !ok {
				err = wrongValue(t, key, at(r))
			} else if e := fn(k, reader.Box(value)); e != nil {
				err = e
			}
			break
		}
	}
	return
}

// we expect to see one, and only one, of the sub keys in the itemValue of m.
func (k *Importer) expectSlot(r reader.Map, expectedType string, sub map[string]Parse) (err error) {
	if t := r.StrOf(itemType); t != expectedType {
		err = wrongType(expectedType, t, at(r))
	} else {
		m := r.MapOf(itemValue)
		t := m.StrOf(itemType)
		if fn, ok := sub[t]; !ok {
			err = wrongType(expectedType, t, at(r))
		} else {
			err = fn(k, m)
		}
	}
	return
}

//
func (k *Importer) expectProg(r reader.Map, slotType string) (ret interface{}, err error) {
	if slot, ok := export.Slots[slotType]; !ok {
		err = errutil.New("unknown slot", slotType, at(r))
	} else if m, e := k.expectSlat(r, slotType); e != nil {
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

//
func (k *Importer) expectEnum(
	m reader.Map,
	expectedType string,
	sub map[string]interface{},
) (ret interface{}, err error) {
	if t := m.StrOf(itemType); t != expectedType {
		err = wrongType(expectedType, t, at(m))
	} else {
		n := m.StrOf(itemValue)
		if i, ok := sub[n]; !ok {
			err = errutil.New("unexpected", expectedType, n)
		} else {
			ret = i
		}
	}
	return
}

// expect a string constant
func (k *Importer) expectConst(m reader.Map, expectedType, expectedValue string) (err error) {
	if t := m.StrOf(itemType); t != expectedType {
		err = wrongType(expectedType, t, at(m))
	} else if v := m.StrOf(itemValue); v != expectedValue {
		err = wrongValue(t, v, at(m))
	}
	return
}

// expect a map value
func (k *Importer) expectSlat(m reader.Map, expectedType string) (ret reader.Map, err error) {
	if _, e := k.expectedType(m, expectedType); e != nil {
		err = e
	} else {
		ret = m.MapOf(itemValue)
	}
	return
}

// expect a string variable
func (k *Importer) expectStr(m reader.Map, expectedType string) (ret string, err error) {
	if t := m.StrOf(itemType); t != expectedType {
		err = wrongType(expectedType, t, at(m))
	} else if v := m.StrOf(itemValue); len(v) == 0 {
		err = wrongValue(t, v, at(m))
	} else {
		ret = v
	}
	return
}

// expect a string variable
func (k *Importer) expectName(m reader.Map, expectedType, cat string) (err error) {
	if t := m.StrOf(itemType); t != expectedType {
		err = wrongType(expectedType, t, at(m))
	} else if v := m.StrOf(itemValue); len(v) == 0 {
		err = wrongValue(t, v, at(m))
	}
	return
}

// helper: check the type of the passed m map
func (k *Importer) expectedType(m reader.Map, expectedType string) (ret string, err error) {
	if t := m.StrOf(itemType); t != expectedType {
		err = wrongType(expectedType, t, at(m))
	} else {
		ret = m.StrOf(itemId)
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
	t string,
	got interface{},
	at string,
) error {
	return errutil.New(t, "has unexpected value", got, "at", at)
}
