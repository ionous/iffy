package qna

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt"
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
}

func (run *Runner) ActivateDomain(domain string, active bool) {
	e := ActivateDomain(run.db, domain, active)
	if e != nil {
		panic(e)
	}
}

func (run *Runner) MakeRecord(kind string) (ret rt.Value, err error) {
	err = errutil.Fmt("couldn't create record of kind %q", kind)
	// fields := make(map[string]rt.Value)
	// ret = generic.NewRecord(op.kind, fields)

	return
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
