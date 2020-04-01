package qna

import (
	"database/sql"

	"github.com/ionous/iffy/scope"
)

func NewRuntime(db *sql.DB) (ret *Runner) {
	return &Runner{
		ObjectValues: NewObjectValues(db),
	}
}

type Runner struct {
	WriterStack
	scope.ScopeStack
	Randomizer
	*ObjectValues
}
