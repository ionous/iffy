package qna

import (
	"database/sql"

	"github.com/ionous/errutil"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/rt/writer"
)

func NewRuntime(db *sql.DB) *Runner {
	var run *Runner
	if plurals, e := NewPlurals(db); e != nil {
		panic(e) // report?
	} else if fields, e := NewFields(db); e != nil {
		panic(e)
	} else {
		run = &Runner{
			db:      db,
			fields:  fields,
			plurals: plurals,
			pairs:   make(valueMap),
			kinds:   qnaKinds{fieldsFor: fields.fieldsFor},
		}
		run.SetWriter(print.NewAutoWriter(writer.NewStdout()))
	}
	return run
}

type Runner struct {
	db *sql.DB
	scope.ScopeStack
	Randomizer
	writer.Sink
	fields  *Fields
	plurals *Plurals
	pairs   valueMap
	kinds   qnaKinds
}

func (run *Runner) ActivateDomain(domain string, active bool) {
	e := ActivateDomain(run.db, domain, active)
	if e != nil {
		panic(e)
	}
}

// notary interface
func (run *Runner) Make(kind string) (ret g.Value, err error) {
	if k, e := run.kinds.KindByName(kind); e != nil {
		err = errutil.New("can't make unknown kind", kind, e)
	} else {
		ret = k.NewRecord()
	}
	return
}

// notary interface
func (run *Runner) Copy(val g.Value) (ret g.Value, err error) {
	return g.CopyValue(run, val)
}

func (run *Runner) KindByName(n string) (*g.Kind, error) {
	return run.kinds.KindByName(n)
}

func (run *Runner) SingularOf(str string) (ret string) {
	if n, e := run.plurals.Singular(str); e != nil {
		ret = str // fix: report e
	} else {
		ret = n
	}
	return
}

func (run *Runner) PluralOf(str string) (ret string) {
	if n, e := run.plurals.Plural(str); e != nil {
		ret = str // fix: report e
	} else {
		ret = n
	}
	return
}
