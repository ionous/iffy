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
	oneTime map[string]bool
	nouns   nounList
}

func NewImporter(srcURI string, db *sql.DB) *Importer {
	registerGob()
	rec := ephemera.NewRecorder(srcURI, db)
	return &Importer{
		Recorder: rec,
		oneTime:  make(map[string]bool),
	}
}

var registeredGob = false

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
