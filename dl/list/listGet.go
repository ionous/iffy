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
				one := v.Float()
				ret = []float64{one}
			case affine.NumList:
				ret = v.Floats()
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
				one := v.String()
				ret = []string{one}
			case affine.TextList:
				ret = v.Strings()
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
				one := v.Record()
				ret = []*g.Record{one}
			case affine.RecordList:
				ret = v.Records()
			default:
				err = errutil.Fmt("cant add %q to a record list", a)
			}
		}
	}
	return
}
