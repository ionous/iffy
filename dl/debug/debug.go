package debug

import (
	"log"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
	"github.com/kr/pretty"
)

var Slats = []composer.Composer{
	(*Log)(nil),
}

type Log struct {
	Text string
	From core.Assignment
}

func (op *Log) Compose() composer.Spec {
	return composer.Spec{
		Name:  "log",
		Group: "debug",
	}
}
func (op *Log) Execute(run rt.Runtime) (err error) {
	if v, e := core.GetAssignedValue(run, op.From); e != nil {
		err = cmdError(op, e)
	} else {
		var i interface{}
		a := v.Affinity()
		switch a {
		case affine.Bool:
			i = v.Bool()
		case affine.Number:
			i = v.Float()
		case affine.NumList:
			i = v.Floats()
		case affine.Text:
			i = v.String()
		case affine.TextList:
			i = v.Strings()
		case affine.Object:
			i = v.String()
		case affine.Record:
			i = pretty.Sprint(generic.RecordToValue(v.Record()))
		case affine.RecordList:
			i = pretty.Sprint(generic.RecordsToValue(v.Records()))
		default:
			e := errutil.New("unknown affinity", a)
			err = cmdError(op, e)
		}
		if err == nil {
			log.Println(op.Text, a, i)
		}
	}
	return
}

func cmdError(op composer.Composer, e error) error {
	return errutil.Append(&core.CommandError{Cmd: op}, e)
}
