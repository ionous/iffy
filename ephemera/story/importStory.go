package story

import (
	"database/sql"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/express"
	"github.com/ionous/iffy/ephemera/reader"
)

func ImportStory(src string, db *sql.DB, m reader.Map) (err error) {
	return ImportStories(src, db, []reader.Map{m})
}

func ImportStories(src string, db *sql.DB, ms []reader.Map) (err error) {
	iffy.RegisterGobs()
	dec := decode.NewDecoder()
	k := NewImporterDecoder(src, db, dec)
	// fix: this is a bit ugly b/c it assumes "AssembleStory" uses the same base kind.
	// and "name" gets used by qna's runtime implementation.
	k.NewImplicitField("name", "kinds", "reserved")
	//
	for _, slats := range iffy.AllSlats {
		dec.AddDefaultCallbacks(slats)
	}
	dec.AddDefaultCallbacks(core.Slats)
	dec.AddCallbacks([]decode.Override{
		// {(*core.Activity)(nil), k.BindRet(func(i *Importer, m reader.Map) (interface{}, error) {
		// 	return imp_activity(i, m) // imp_activity returns *Activity, BindRet expects interface{}
		// })},
		{(*pattern.DetermineAct)(nil), k.BindRet(imp_determine_act)},
		{(*pattern.DetermineNum)(nil), k.BindRet(imp_determine_num)},
		{(*pattern.DetermineText)(nil), k.BindRet(imp_determine_text)},
		{(*pattern.DetermineBool)(nil), k.BindRet(imp_determine_bool)},
		{(*pattern.DetermineNumList)(nil), k.BindRet(imp_determine_num_list)},
		{(*pattern.DetermineTextList)(nil), k.BindRet(imp_determine_text_list)},
		//
		{(*core.Text)(nil), k.BindRet(imp_text_value)},
		{(*core.CycleText)(nil), k.BindRet(imp_cycle_text)},
		{(*core.ShuffleText)(nil), k.BindRet(imp_shuffle_text)},
		{(*core.StoppingText)(nil), k.BindRet(imp_stopping_text)},
		//
		{(*express.Render)(nil), k.BindRet(imp_render_template)},
	})
	//
	for _, m := range ms {
		if e := imp_story(k, m); e != nil {
			err = e
			break
		}
	}
	return
}
