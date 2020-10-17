package story

import (
	"database/sql"
	"strconv"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

// Importer helps read story specific json.
type Importer struct {
	*ephemera.Recorder
	// sometimes the importer needs to define a singleton like type or instance
	oneTime     map[string]bool
	decoder     *decode.Decoder
	autoCounter autoCounter
	entireGame  ephemera.Named
	StoryEnv
}

// helper for making auto variables.
type autoCounter map[string]uint64

func (m *autoCounter) next(name string) string {
	c := (*m)[name] + 1
	(*m)[name] = c
	return name + "#" + strconv.FormatUint(c, 36)
}

func NewImporter(srcURI string, db *sql.DB) *Importer {
	return NewImporterDecoder(srcURI, db, decode.NewDecoder())
}

func NewImporterDecoder(srcURI string, db *sql.DB, dec *decode.Decoder) *Importer {
	rec := ephemera.NewRecorder(srcURI, db)
	return &Importer{
		Recorder:    rec,
		oneTime:     make(map[string]bool),
		decoder:     dec,
		autoCounter: make(autoCounter),
	}
}

func (k *Importer) NewName(name, category, ofs string) ephemera.Named {
	domain := k.Current.Domain
	if !domain.IsValid() {
		domain = k.gameDomain()
	}
	return k.NewDomainName(domain, name, category, ofs)
}

func (k *Importer) gameDomain() ephemera.Named {
	if !k.entireGame.IsValid() {
		k.entireGame = k.Recorder.NewName("Entire Game", tables.NAMED_SCENE, "internal")
	}
	return k.entireGame
}

// return true if m is the first time once has been called with the specified string.
func (k *Importer) Once(s string) (ret bool) {
	if !k.oneTime[s] {
		k.oneTime[s] = true
		ret = true
	}
	return
}

// adapt an importer friendly function to the ephemera reader callback
func (k *Importer) Bind(cb func(*Importer, reader.Map) (err error)) reader.ReadMap {
	return func(m reader.Map) error {
		return cb(k, m)
	}
}

// adapt an importer friendly function to the ephemera reader callback
func (k *Importer) BindRet(cb func(*Importer, reader.Map) (ret interface{}, err error)) decode.ReadRet {
	return func(m reader.Map) (interface{}, error) {
		return cb(k, m)
	}
}

// read the passed map as if it contained a slot. ex bool_eval, etc.
func (k *Importer) DecodeSlot(m reader.Map, slotType string, outPtr interface{}) (err error) {
	if m, e := reader.Unpack(m, slotType); e != nil {
		err = e // reuses: "slat" to unpack the contained map.
	} else {
		err = k.DecodeAny(m, outPtr)
	}
	return
}

func (k *Importer) DecodeAny(m reader.Map, outPtr interface{}) (err error) {
	if k.decoder == nil {
		err = errutil.New("no decoder initialized")
	} else if m != nil {
		err = k.decoder.ReadProg(m, outPtr)
	}
	return
}

// // NewImplicitField declares an assembler specified field
// func (k *Importer) NewImplicitField(field, kind, fieldType string) {
// 	if src := "implicit " + kind + "." + field; k.Once(src) {
// 		domain := k.gameDomain()
// 		kKind := k.NewDomainName(domain, kind, tables.NAMED_KINDS, src)
// 		kField := k.NewDomainName(domain, field, tables.NAMED_FIELD, src)
// 		k.NewField(kKind, kField, fieldType)
// 	}
// }

// NewImplicitAspect declares an assembler specified aspect and its traits
func (k *Importer) NewImplicitAspect(aspect, kind string, traits ...string) {
	if src := "implicit " + kind + "." + aspect; k.Once(src) {
		domain := k.gameDomain()
		kKind := k.NewDomainName(domain, kind, tables.NAMED_KINDS, src)
		kAspect := k.NewDomainName(domain, aspect, tables.NAMED_ASPECT, src)
		k.NewAspect(kAspect)
		k.NewField(kKind, kAspect, tables.PRIM_ASPECT)
		for i, trait := range traits {
			kTrait := k.NewDomainName(domain, trait, tables.NAMED_TRAIT, src)
			k.NewTrait(kTrait, kAspect, i)
		}
	}
}
