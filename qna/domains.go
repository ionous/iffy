package qna

import (
	"database/sql"

	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/tables"
)

func ActivateDomain(db *sql.DB, domain string, active bool) error {
	name := lang.Camelize(domain)
	_, e := db.Exec(run_domain, name, active)
	return e
}

var run_domain = tables.InsertWith("run_domain",
	"on conflict(domain) do update set active=excluded.active;",
	"domain", "active")
