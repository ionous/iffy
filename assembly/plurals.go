package assembly

import (
	"database/sql"

	"github.com/ionous/iffy/lang"
	"github.com/mattn/go-sqlite3"
)

const SqlCustomDriver = "iffy_asm"

// requires using sql.Open(SqlCustomDriver)
func init() {
	// since we have to have app code for assembling plurals,
	// might as well use an extension to simplify the processing.
	sql.Register(SqlCustomDriver, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) (err error) {
			pure := true // true means the function's return depends only on its inputs.
			if e := conn.RegisterFunc("pluralize", lang.Pluralize, pure); e != nil {
				err = e
			}
			return
		},
	})
}

// build list of plurals
// currently this is just from singular kinds
// eventually it would be from pluralization modeling statements
func AssemblePlurals(asm *Assembler) (err error) {
	_, e := asm.cache.DB().Exec(
		`insert into mdl_plural
	select distinct name, pluralize(name) 
	from eph_named
			where category='singular_kind'`)
	return e
}
