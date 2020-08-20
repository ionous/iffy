package story

import (
	"strconv"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/rt"
)

func imp_cycle_text(k *Importer, r reader.Map) (ret interface{}, err error) {
	if m, e := reader.Unpack(r, "cycle_text"); e != nil {
		err = e
	} else if seq, e := fromSequence(k, m); e != nil {
		err = e
	} else {
		ret = &core.CycleText{Sequence: seq}
	}
	return
}

func imp_shuffle_text(k *Importer, r reader.Map) (ret interface{}, err error) {
	if m, e := reader.Unpack(r, "shuffle_text"); e != nil {
		err = e
	} else if seq, e := fromSequence(k, m); e != nil {
		err = e
	} else {
		ret = &core.ShuffleText{Sequence: seq}
	}
	return
}

func imp_stopping_text(k *Importer, r reader.Map) (ret interface{}, err error) {
	if m, e := reader.Unpack(r, "stopping_text"); e != nil {
		err = e
	} else if seq, e := fromSequence(k, m); e != nil {
		err = e
	} else {
		ret = &core.StoppingText{Sequence: seq}
	}
	return
}

func fromSequence(k *Importer, m reader.Map) (ret core.Sequence, err error) {
	var ps []rt.TextEval
	if e := reader.Repeats(m.SliceOf("$PARTS"),
		func(m reader.Map) (err error) {
			if p, e := k.DecodeSlot(m, "text_eval"); e != nil {
				err = e
			} else {
				ps = append(ps, p.(rt.TextEval))
			}
			return
		}); e != nil {
		err = e
	} else {
		counter := getCounter(k, m)
		ret = core.Sequence{counter, ps}
	}
	return
}

func getCounter(k *Importer, m reader.Map) (ret string) {
	if at := reader.At(m); len(at) > 0 {
		ret = at
	} else {
		k.autoCounter++
		ret = "autoimp" + strconv.Itoa(k.autoCounter)
	}
	return
}
