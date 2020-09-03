package qna

import (
	"database/sql"

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
			Fields:  fields,
			plurals: plurals,
		}
		run.SetWriter(print.NewAutoWriter(writer.NewStdout()))
		run.PushScope(&NounScope{fields: fields})
	}
	return run
}

type Runner struct {
	db *sql.DB
	scope.ScopeStack
	Randomizer
	writer.Sink
	*Fields
	plurals *Plurals
}

func (run *Runner) ActivateDomain(domain string, active bool) {
	e := ActivateDomain(run.db, domain, active)
	if e != nil {
		panic(e)
	}
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
