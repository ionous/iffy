package qna

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type boolEval struct{ eval rt.BoolEval }

func (q boolEval) Snapshot(run rt.Runtime) (ret g.Value, err error) {
	if v, e := rt.GetBool(run, q.eval); e != nil {
		err = e
	} else {
		ret, err = g.ValueOf(v)
	}
	return
}

type numEval struct{ eval rt.NumberEval }

func (q numEval) Snapshot(run rt.Runtime) (ret g.Value, err error) {
	if v, e := rt.GetNumber(run, q.eval); e != nil {
		err = e
	} else {
		ret, err = g.ValueOf(v)
	}
	return
}

type textEval struct{ eval rt.TextEval }

func (q textEval) Snapshot(run rt.Runtime) (ret g.Value, err error) {
	if v, e := rt.GetText(run, q.eval); e != nil {
		err = e
	} else {
		ret, err = g.ValueOf(v)
	}
	return
}

type numListEval struct{ eval rt.NumListEval }

func (q numListEval) Snapshot(run rt.Runtime) (ret g.Value, err error) {
	if v, e := rt.GetNumList(run, q.eval); e != nil {
		err = e
	} else {
		ret, err = g.ValueOf(v)
	}
	return
}

type textListEval struct{ eval rt.TextListEval }

func (q textListEval) Snapshot(run rt.Runtime) (ret g.Value, err error) {
	if v, e := rt.GetTextList(run, q.eval); e != nil {
		err = e
	} else {
		ret, err = g.ValueOf(v)
	}
	return
}

func newEval(a affine.Affinity, buf []byte) (ret snapper, err error) {
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
