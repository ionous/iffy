package assembly

import (
	"database/sql"

	"github.com/ionous/errutil"
)

func AssembleTests(asm *Assembler) (err error) {
	db := asm.cache.DB()
	return copyTests(db)
}

func copyTests(db *sql.DB) (err error) {
	if _, e := db.Exec(
		`insert into mdl_prog 
		select name, type, prog as bytes
		from asm_check`); e != nil {
		err = errutil.New("copyTests:", e)
	}
	return
}
