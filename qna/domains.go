package qna

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/tables"
)

func ActivateDomain(db *sql.DB, domain string, active bool) (err error) {
	name := lang.Camelize(domain)
	if _, e := db.Exec(run_domain, name, active); e != nil {
		err = errutil.New("ActivateDomain", e)
	}
	return
}

var run_domain = tables.InsertWith("run_domain",
	"on conflict(domain) do update set active=excluded.active;",
	"domain", "active")
