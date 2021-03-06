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

// LogLevel controls how much debugging to print
// The default level (0) means log everything but notes,
// otherwise it logs only at the named level and higher.
var LogLevel Level

type Log struct {
	Value core.Assignment `if:"selector"`
	Level Level           `if:"selector"`
}

func (op *Log) Compose() composer.Spec {
	return composer.Spec{
		Name:   "debug_log",
		Group:  "debug",
		Fluent: &composer.Fluid{Name: "log", Role: composer.Command},
	}
}

func (op *Log) Execute(run rt.Runtime) (err error) {
	// fix? at this time we cant guarantee a lack of side-effects
	// so we always eval even if we don't print.
	if v, e := core.GetAssignedValue(run, op.Value); e != nil {
		err = cmdError(op, e)
	} else {
		var i interface{}
		switch a := v.Affinity(); a {
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
		case affine.Object: // this probably never happens... but just in case.
			i = v.String()
		case affine.Record:
			i = pretty.Sprint(generic.RecordToValue(v.Record()))
		case affine.RecordList:
			i = pretty.Sprint(generic.RecordsToValue(v.Records()))
		default:
			e := errutil.New("unknown affinity", a)
			err = cmdError(op, e)
		}
		if err == nil && ((LogLevel != 0 && op.Level >= LogLevel) || (LogLevel == 0 && op.Level > Note)) {
			log.Println(op.Level.Header(), i)
		}
	}
	return
}
