package story

import (
	"database/sql"
	"log"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/render"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/reader"
)

func ImportStory(src string, db *sql.DB, m reader.Map) (err error) {
	return ImportStories(src, db, []reader.Map{m})
}

func ImportStories(src string, db *sql.DB, ms []reader.Map) (err error) {
	iffy.RegisterGobs()
	var ds reader.Dilemmas
	dec := decode.NewDecoderReporter(src, ds.Report)
	k := NewImporterDecoder(src, db, dec)
	//
	for _, slats := range iffy.AllSlats {
		dec.AddDefaultCallbacks(slats)
	}
	// for _, slat := range Slats {
	// 	if _, ok := slat.(Imported); !ok {
	// 		dec.AddCallback(slat, nil)
	// 	} else {
	// 		dec.AddCallback(slat, func(m reader.Map) (ret interface{}, err error) {
	// 			slatElm := r.TypeOf(slat).Elem()
	// 			slatPtr := r.New(slatElm)
	// 			dec.ReadFields(reader.At(m), slatPtr.Elem(), m.MapOf(reader.ItemValue))
	// 			op := slatPtr.Interface().(Imported)
	// 			if e := op.Imported(k); e != nil {
	// 				err = e
	// 			} else {
	// 				ret = op
	// 			}
	// 			return
	// 		})
	// 	}
	// }
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
		{(*render.Template)(nil), k.BindRet(imp_render_template)},
	})

	for _, m := range ms {
		if e := imp_story(k, m); e != nil {
			err = e
			break
		}
	}
	// for _, m := range ms {
	// 	if i, e := dec.ReadSpec(m); e != nil {
	// 		err = e
	// 		break
	// 	} else {
	// 		pretty.Println(i)
	// 	}
	// }
	reader.PrintDilemmas(log.Writer(), ds)

	return
}

type Imported interface {
	Imported(*Importer) (err error)
}
