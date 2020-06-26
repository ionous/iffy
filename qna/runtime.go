package qna

import (
	"database/sql"

	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/tables"
)

func NewRuntime(db *sql.DB) *Runner {
	cache := tables.NewCache(db)
	fields := NewObjectValues(cache)
	run := &Runner{
		Fields: fields,
	}
	run.PushScope(&NounScope{fields: fields})
	return run
}

type Runner struct {
	print.Stack
	scope.ScopeStack
	Randomizer
	*Fields
}
