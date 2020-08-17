package imp

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

// Porter helps read json.
type Porter struct {
	*ephemera.Recorder
	// sometimes the importer needs to define a singleton like type or instance
	oneTime     map[string]bool
	decoder     *decode.Decoder
	AutoCounter int // helper for making auto variables.
}

func NewImporter(srcURI string, db *sql.DB) *Porter {
	return NewImporterDecoder(srcURI, db, decode.NewDecoder())
}

func NewImporterDecoder(srcURI string, db *sql.DB, dec *decode.Decoder) *Porter {
	rec := ephemera.NewNormalizingRecorder(srcURI, db)
	return &Porter{
		Recorder: rec,
		oneTime:  make(map[string]bool),
		decoder:  dec,
	}
}

// return true if m is the first time once has been called with the specified string.
func (k *Porter) Once(s string) (ret bool) {
	if !k.oneTime[s] {
		k.oneTime[s] = true
		ret = true
	}
	return
}

// adapt an importer friendly function to the ephemera reader callback
func (k *Porter) Bind(cb func(*Porter, reader.Map) (err error)) reader.ReadMap {
	return func(m reader.Map) error {
		return cb(k, m)
	}
}

// adapt an importer friendly function to the ephemera reader callback
func (k *Porter) BindRet(cb func(*Porter, reader.Map) (ret interface{}, err error)) decode.ReadRet {
	return func(m reader.Map) (interface{}, error) {
		return cb(k, m)
	}
}

// read the passed map as if it contained a slot. ex bool_eval, etc.
func (k *Porter) DecodeSlot(m reader.Map, slotType string) (ret interface{}, err error) {
	if m, e := reader.Unpack(m, slotType); e != nil {
		err = e // reuses: "slat" to unpack the contained map.
	} else {
		ret, err = k.DecodeAny(m)
	}
	return
}
func (k *Porter) DecodeAny(m reader.Map) (ret interface{}, err error) {
	if k.decoder == nil {
		err = errutil.New("no decoder initialized")
	} else if m != nil {
		ret, err = k.decoder.ReadProg(m)
	}
	return
}

// NewProg add the passed cmd ephemera.
func (k *Porter) NewProg(typeName string, cmd interface{}) (ret ephemera.Prog, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if e := enc.Encode(cmd); e != nil {
		err = e
	} else {
		ret = k.Recorder.NewProg(typeName, buf.Bytes())
	}
	return
}

// NewImplicitField declares an assembler specified field
func (k *Porter) NewImplicitField(field, kind, fieldType string) {
	if src := "implicit " + kind + "." + field; k.Once(src) {
		kKind := k.NewName(kind, tables.NAMED_KINDS, src)
		kField := k.NewName(field, tables.NAMED_FIELD, src)
		k.NewField(kKind, kField, fieldType)
	}
}

// NewImplicitAspect declares an assembler specified aspect and its traits
func (k *Porter) NewImplicitAspect(aspect, kind string, traits ...string) {
	if src := "implicit " + kind + "." + aspect; k.Once(src) {
		kKind := k.NewName(kind, tables.NAMED_KINDS, src)
		kAspect := k.NewName(aspect, tables.NAMED_ASPECT, src)
		k.NewAspect(kAspect)
		k.NewField(kKind, kAspect, tables.PRIM_ASPECT)
		for i, trait := range traits {
			kTrait := k.NewName(trait, tables.NAMED_TRAIT, src)
			k.NewTrait(kTrait, kAspect, i)
		}
	}
}
