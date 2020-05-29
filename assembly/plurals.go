package assembly

import (
	"database/sql"

	"github.com/ionous/iffy/lang"
	"github.com/mattn/go-sqlite3"
)

func init() {
	// since we have to have app code for assembling plurals,
	// might as well use an extension to simplify the processing.
	sql.Register("iffy_asm", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) (err error) {
			pure := true // true means the function's return depends only on its inputs.
			if e := conn.RegisterFunc("pluralize", lang.Pluralize, pure); e != nil {
				err = e
			}
			return
		},
	})
}
