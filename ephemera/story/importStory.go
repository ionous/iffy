package story

import (
	"database/sql"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/ephemera/reader"
)

func ImportStory(src string, db *sql.DB, m reader.Map) (err error) {
	iffy.RegisterGobs()
	storyNouns.names = nil
	dec := decode.NewDecoder()
	k := imp.NewImporterDecoder(src, db, dec)
	//
	for _, slats := range iffy.AllSlats {
		dec.AddDefaultCallbacks(slats)
	}
	dec.AddDefaultCallbacks(core.Slats)
	dec.AddCallbacks([]decode.Override{
		{(*core.DetermineAct)(nil), k.BindRet(imp_determine_act)},
		{(*core.DetermineNum)(nil), k.BindRet(imp_determine_num)},
		{(*core.DetermineText)(nil), k.BindRet(imp_determine_text)},
		{(*core.DetermineBool)(nil), k.BindRet(imp_determine_bool)},
		{(*core.DetermineNumList)(nil), k.BindRet(imp_determine_num_list)},
		{(*core.DetermineTextList)(nil), k.BindRet(imp_determine_text_list)},

		{(*core.CycleText)(nil), k.BindRet(imp_cycle_text)},
		{(*core.ShuffleText)(nil), k.BindRet(imp_shuffle_text)},
		{(*core.StoppingText)(nil), k.BindRet(imp_stopping_text)},
	})
	//
	return imp_story(k, m)
}
