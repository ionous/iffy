package tables

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // queries are specific to sqlite, so force the sqlite driver.
)

// CreateEphemera creates the tables listed in ephemera.sql
//go:generate templify -p tables -o ephemera.gen.go ephemera.sql
func CreateEphemera(db *sql.DB) error {
	_, e := db.Exec(ephemeraTemplate())
	return e
}

// CreateAssembly creates the tables listed in assembly.sql
//go:generate templify -p tables -o assembly.gen.go assembly.sql
func CreateAssembly(db *sql.DB) error {
	_, e := db.Exec(assemblyTemplate())
	return e
}

// CreateModel creates the tables listed in model.sql
//go:generate templify -p tables -o model.gen.go model.sql
func CreateModel(db *sql.DB) error {
	_, e := db.Exec(modelTemplate())
	return e
}

// CreateRun creates the tables listed in run.sql
//go:generate templify -p tables -o run.gen.go run.sql
func CreateRun(db *sql.DB) error {
	_, e := db.Exec(runTemplate())
	return e
}

// CreateRunViews creates the tables listed in runViews.sql
//go:generate templify -p tables -o runViews.gen.go runViews.sql
func CreateRunViews(db *sql.DB) error {
	_, e := db.Exec(runViewsTemplate())
	return e
}

// CreateAll tables listed in the various .sql files.
// fix? long term, iffy should throw out the ephemera and assembly after we are done with them.
func CreateAll(db *sql.DB) (err error) {
	if e := CreateEphemera(db); e != nil {
		err = e
	} else if e := CreateAssembly(db); e != nil {
		err = e
	} else if e := CreateModel(db); e != nil {
		err = e
	} else if e := CreateRun(db); e != nil {
		err = e
	} else if e := CreateRunViews(db); e != nil {
		err = e
	}
	return
}
