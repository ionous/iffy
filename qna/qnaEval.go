package qna

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

type boolEval struct{ eval rt.BoolEval }

func (q boolEval) Affinity() affine.Affinity {
	return affine.Bool
}
func (q boolEval) Snapshot(run rt.Runtime) (g.Value, error) {
	return safe.GetBool(run, q.eval)
}

type numEval struct{ eval rt.NumberEval }

func (q numEval) Affinity() affine.Affinity {
	return affine.Number
}
func (q numEval) Snapshot(run rt.Runtime) (g.Value, error) {
	return safe.GetNumber(run, q.eval)
}

type textEval struct{ eval rt.TextEval }

func (q textEval) Affinity() affine.Affinity {
	return affine.Text
}
func (q textEval) Snapshot(run rt.Runtime) (g.Value, error) {
	return safe.GetText(run, q.eval)
}

type numListEval struct{ eval rt.NumListEval }

func (q numListEval) Affinity() affine.Affinity {
	return affine.NumList
}
func (q numListEval) Snapshot(run rt.Runtime) (g.Value, error) {
	return safe.GetNumList(run, q.eval)
}

type textListEval struct{ eval rt.TextListEval }

func (q textListEval) Affinity() affine.Affinity {
	return affine.TextList
}
func (q textListEval) Snapshot(run rt.Runtime) (g.Value, error) {
	return safe.GetTextList(run, q.eval)
}

func newEval(a affine.Affinity, buf []byte) (ret qnaValue, err error) {
	switch a {
	case affine.Bool:
		var eval rt.BoolEval
		if e := bytesToEval(buf, &eval); e != nil {
			err = e
		} else {
			ret = &boolEval{eval}
		}
	case affine.Number:
		var eval rt.NumberEval
		if e := bytesToEval(buf, &eval); e != nil {
			err = e
		} else {
			ret = &numEval{eval}
		}
	case affine.Text:
		var eval rt.TextEval
		if e := bytesToEval(buf, &eval); e != nil {
			err = e
		} else {
			ret = &textEval{eval}
		}
	case affine.NumList:
		var eval rt.NumListEval
		if e := bytesToEval(buf, &eval); e != nil {
			err = e
		} else {
			ret = &numListEval{eval}
		}
	case affine.TextList:
		var eval rt.TextListEval
		if e := bytesToEval(buf, &eval); e != nil {
			err = e
		} else {
			ret = &textListEval{eval}
		}
	}
	return
}
