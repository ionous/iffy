package qna

import (
	"database/sql"
	"log"

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
	// fix: we want activate to return a list of *newly* active domains;
	// this is a lot like a state transition list.
	// you might be able to do something clever? like "time activated" --
	// and select for current time (even possibly override the sql timer with game round )
	// or use the domain path.
	if active {
		if cnt, e := run.fields.UpdatePairs(domain); e != nil {
			panic(e)
		} else {
			log.Println("activate domain affected", cnt, "rows")
		}
	}
}

func (run *Runner) GetKindByName(n string) (*g.Kind, error) {
	return run.kinds.GetKindByName(n)
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

func (run *Runner) Relate(a, b, relation string) (err error) {
	return
}

func (run *Runner) Relatives(a, relation string) (ret []string, err error) {
	return
}
