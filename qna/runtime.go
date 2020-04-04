package qna

import (
	"database/sql"

	"github.com/ionous/iffy/rt/print"
	"github.com/ionous/iffy/rt/scope"
)

func NewRuntime(db *sql.DB) (ret *Runner) {
	return &Runner{
		Fields: NewObjectValues(db),
	}
}

type Runner struct {
	print.Stack
	scope.ScopeStack
	Randomizer
	*Fields
}
