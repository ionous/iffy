package list

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

func getNewFloats(run rt.Runtime, assign core.Assignment) (ret []float64, err error) {
	if assign != nil {
		if v, e := assign.GetAssignedValue(run); e != nil {
			err = e
		} else {
			switch a := v.Affinity(); a {
			case affine.Number:
				if one, e := v.GetNumber(); e != nil {
					err = e
				} else {
					ret = []float64{one}
				}
			case affine.NumList:
				if many, e := v.GetNumList(); e != nil {
					err = e
				} else {
					ret = many
				}
			default:
				err = errutil.Fmt("cant add %q to a num list", a)
			}
		}
	}
	return
}

func getNewStrings(run rt.Runtime, assign core.Assignment) (ret []string, err error) {
	if assign != nil {
		if v, e := assign.GetAssignedValue(run); e != nil {
			err = e
		} else {
			switch a := v.Affinity(); a {
			case affine.Text:
				if one, e := v.GetText(); e != nil {
					err = e
				} else {
					ret = []string{one}
				}
			case affine.TextList:
				if many, e := v.GetTextList(); e != nil {
					err = e
				} else {
					ret = many
				}
			default:
				err = errutil.Fmt("cant add %q to a text list", a)
			}
		}
	}
	return
}

func getNewRecords(run rt.Runtime, oldType string,
	assign core.Assignment) (ret []*g.Record, err error) {
	if assign != nil {
		if v, e := assign.GetAssignedValue(run); e != nil {
			err = e
		} else if v.Type() != oldType {
			err = errutil.New("mismatched record types", oldType)
		} else {
			switch a := v.Affinity(); a {
			case affine.Record:
				if one, e := v.GetRecord(); e != nil {
					err = e
				} else {
					ret = []*g.Record{one}
				}

			case affine.RecordList:
				if many, e := v.GetRecordList(); e != nil {
					err = e
				} else {
					ret = many
				}
			default:
				err = errutil.Fmt("cant add %q to a record list", a)
			}
		}
	}
	return
}
