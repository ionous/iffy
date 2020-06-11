package qna

import (
	"database/sql"

	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/tables"
)

func NewRuntime(db *sql.DB) *Runner {
	cache := tables.NewCache(db)
	run := &Runner{
		Fields: NewObjectValues(cache),
	}
	run.PushScope(&NounScope{db: cache})
	return run
}

type Runner struct {
	print.Stack
	scope.ScopeStack
	Randomizer
	*Fields
}
