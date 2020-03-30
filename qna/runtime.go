package qna

import (
	"database/sql"
)

func NewRuntime(db *sql.DB) (ret *Runner) {
	return &Runner{
		ObjectValues: NewObjectValues(db),
	}
}

type Runner struct {
	WriterStack
	ScopeStack
	Randomizer
	*ObjectValues
}
