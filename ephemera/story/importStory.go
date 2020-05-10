package story

import (
	"database/sql"
	"encoding/gob"

	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/export"
)

func ImportStory(src string, m reader.Map, db *sql.DB) error {
	registerGob()
	storyNouns.names = nil
	dec := decode.NewDecoder()
	k := imp.NewImporterDecoder(src, db, dec)
	return imp_story(k, m)
}

var registeredGob = false

// register imperative commands exposed by the exporter
func registerGob() {
	if !registeredGob {
		export.Register(gob.Register)
		registeredGob = true
	}
}
