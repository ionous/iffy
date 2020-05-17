package story

import (
	"database/sql"
	"encoding/gob"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/export"
)

func ImportStory(src string, db *sql.DB, m reader.Map) (err error) {
	registerGob()
	storyNouns.names = nil
	dec := decode.NewDecoder()
	k := imp.NewImporterDecoder(src, db, dec)
	//
	dec.AddDefaultCallbacks(export.Slats)
	dec.AddCallbacks([]decode.Override{
		{(*core.DetermineAct)(nil), k.BindRet(imp_determine_act)},
		{(*core.DetermineNum)(nil), k.BindRet(imp_determine_num)},
		{(*core.DetermineText)(nil), k.BindRet(imp_determine_text)},
		{(*core.DetermineBool)(nil), k.BindRet(imp_determine_bool)},
		{(*core.DetermineNumList)(nil), k.BindRet(imp_determine_num_list)},
		{(*core.DetermineTextList)(nil), k.BindRet(imp_determine_text_list)},
	})
	//
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
