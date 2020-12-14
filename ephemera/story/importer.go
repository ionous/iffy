package story

import (
	"database/sql"
	r "reflect"
	"strconv"

	"github.com/ionous/iffy/dl/composer"
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

//
func (k *Importer) AddModel(model []composer.Slat) {
	type stubImporter interface {
		ImportStub(k *Importer) (ret interface{}, err error)
	}
	dec := k.decoder
	for _, cmd := range model {
		if _, ok := cmd.(stubImporter); !ok {
			dec.AddCallback(cmd, nil)
		} else {
			// need to pin the loop variable for the callback
			// so pin the type. why not.
			rtype := r.TypeOf(cmd).Elem()
			dec.AddCallback(cmd, func(m reader.Map) (ret interface{}, err error) {
				// create an instance of the stub
				op := r.New(rtype)
				// read it in
				dec.ReadFields(reader.At(m), op.Elem(), m.MapOf(reader.ItemValue))
				// convert it
				stub := op.Interface().(stubImporter)
				return stub.ImportStub(k)
			})
		}
	}
}

//
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
