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

// fix: queue could handle this. stack an event queue.
// maybe then we could stack "last noun" in there?
type CategoryEvent func(ephemera.Named)

type imperativeHandler func(Importer, reader.Map) (ret interface{}, err error)
type CommandMap map[string]imperativeHandler

// Importer helps read json.
type Importer struct {
	eph *ephemera.Recorder
	// sometimes the importer needs to define a singleton like type or instance
	oneTime map[string]bool
	// a list of recently referenced nouns
	// helps simplify some aspects of importing
	nouns nounList
	// some commands can add information to the world model
	// ( ex. variable references, type casting, etc. )
	// those need special handling by the importer.
	commands CommandMap
}

func NewImporter(srcURI string, db *sql.DB, cmds CommandMap) *Importer {
	registerGob()
	rec := ephemera.NewRecorder(srcURI, db)
	return &Importer{
		eph:      rec,
		oneTime:  make(map[string]bool),
		commands: cmds,
	}
}

var registeredGob = false

// register imperative commands exposed by the exporter
func registerGob() {
	if !registeredGob {
		export.Register(gob.Register)
		registeredGob = true
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

// adapt an importer friendly function to the ephemera reader callback
func (k *Importer) bind(cb func(*Importer, reader.Map) (err error)) reader.ReadMap {
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
		ret = k.eph.NewProg(t, buf.Bytes())
	}
	return
}
