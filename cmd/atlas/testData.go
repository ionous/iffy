package main

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

//go:generate templify -p main -o testData.gen.go testData.sql
func createTestData(db *sql.DB) (err error) {
	if e := tables.CreateModel(db); e != nil {
		err = e
	} else if _, e := db.Exec(testDataTemplate()); e != nil {
		err = errutil.New("createTestData", e)
	}
	return
}
