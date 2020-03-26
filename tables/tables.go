package tables

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

//go:generate templify -p tables -o ephemera.gen.go ephemera.sql
func CreateEphemera(db *sql.DB) error {
	_, e := db.Exec(ephemeraTemplate())
	return e
}

//go:generate templify -p tables -o assembly.gen.go assembly.sql
func CreateAssembly(db *sql.DB) error {
	_, e := db.Exec(assemblyTemplate())
	return e
}

//go:generate templify -p tables -o model.gen.go model.sql
func CreateModel(db *sql.DB) error {
	_, e := db.Exec(modelTemplate())
	return e
}

//go:generate templify -p tables -o run.gen.go run.sql
func CreateRun(db *sql.DB) error {
	_, e := db.Exec(runTemplate())
	return e
}

//go:generate templify -p tables -o runViews.gen.go runViews.sql
func CreateRunViews(db *sql.DB) error {
	_, e := db.Exec(runViewsTemplate())
	return e
}
