package qna

import (
	"database/sql"
	"log"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/object"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/rt/writer"
	"github.com/ionous/iffy/tables"
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
			kinds: qnaKinds{
				typeOf:    fields.typeOf,
				fieldsFor: fields.fieldsFor,
				traitsFor: fields.traitsFor,
			},
			activeNouns:   activeNouns{q: fields.activeNouns},
			relativeKinds: relativeKinds{q: fields.relativeKinds},
			nounLocale:    nounLocale{q: fields.relativesOf},
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
	activeNouns
	relativeKinds
	nounLocale
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
		if _, e := run.fields.UpdatePairs(domain); e != nil {
			panic(e)
		} else {
			// log.Println("activate domain", domain, "affected", cnt, "noun pairs")
		}
	}
	run.activeNouns.reset()
	run.nounLocale.reset()
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

// assumes a and b are valid nouns
func (run *Runner) RelateTo(a, b, relation string) (err error) {
	// we validate inputs in go rather than sql b/c
	// a, the sql for validation gets big and ugly quick
	// b. we get better reporting this way.
	// -- perhaps there could be a standalone validation query that returns nice errors
	// but this is okay for now.
	if !run.isActive(a) {
		err = g.UnknownObject(a)
	} else if !run.isActive(b) {
		err = g.UnknownObject(b)
	} else if ak, e := run.GetField(object.Kinds, a); e != nil {
		err = e
	} else if bk, e := run.GetField(object.Kinds, b); e != nil {
		err = e
	} else if rel := run.relativeKind(relation); !compatibleKind(ak.String(), rel.kind) {
		err = errutil.Fmt("relation %s expects %s doesnt support %s ( a kind of %s )", relation, rel.kind, a, ak.String())
	} else if !compatibleKind(bk.String(), rel.otherKind) {
		err = errutil.Fmt("relation %s expects %s doesnt support %s ( a kind of %s )", relation, rel.otherKind, b, bk.String())
	} else if res, e := run.fields.relateTo.Exec(a, b, relation, rel.cardinality); e != nil {
		err = e
	} else {
		log.Println(tables.RowsAffected(res), "rows affected relating", a, "to", b, "via", relation)
	}
	return
}

// assumes a is a valid noun
func (run *Runner) RelativesOf(a, relation string) (ret []string, err error) {
	if !run.isActive(a) {
		err = g.UnknownObject(a)
	} else if rows, e := run.fields.relativesOf.Query(a, relation); e != nil {
		err = e
	} else {
		var otherNoun string
		err = tables.ScanAll(rows, func() (err error) {
			ret = append(ret, otherNoun)
			return
		}, &otherNoun)
	}
	return
}

// assumes b is a valid noun
func (run *Runner) ReciprocalOf(b, relation string) (ret []string, err error) {
	if !run.isActive(b) {
		err = g.UnknownObject(b)
	} else if rows, e := run.fields.reciprocalOf.Query(b, relation); e != nil {
		err = e
	} else {
		var noun string
		err = tables.ScanAll(rows, func() (err error) {
			ret = append(ret, noun)
			return
		}, &noun)
	}
	return
}
