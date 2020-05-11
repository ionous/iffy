package story

import (
	"database/sql"
	"encoding/gob"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/tables"
)

func ImportStory(src string, m reader.Map, db *sql.DB) (err error) {
	registerGob()
	storyNouns.names = nil
	if e := tables.CreateEphemera(db); e != nil {
		err = e
	} else {
		dec := decode.NewDecoder()
		k := imp.NewImporterDecoder(src, db, dec)
		//
		dec.AddDefaultCallbacks(export.Slats)
		dec.AddCallbacks([]decode.Override{
			{(*core.DetermineNum)(nil), k.BindRet(imp_determine_num)},
		})
		//
		err = imp_story(k, m)
	}
	return
}

var registeredGob = false

// register imperative commands exposed by the exporter
func registerGob() {
	if !registeredGob {
		export.Register(gob.Register)
		registeredGob = true
	}
}
