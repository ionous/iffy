package qna

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"log"

	"github.com/ionous/errutil"

	"github.com/ionous/iffy/check"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/tables"
)

// CheckAll tests stored in the passed db.
// It logs the results of running the tests, and only returns error on critical errors.
func CheckAll(db *sql.DB) (err error) {
	run := NewRuntime(db)
	var prog []byte
	var name, expect string
	var tests []check.CheckOutput
	if e := tables.QueryAll(db,
		`select ck.name, pg.bytes, ck.expect
		from mdl_check as ck
		join mdl_prog pg
			on ((pg.type = 'execute') and 
				(pg.type = ck.type) and 
				(pg.name = ck.name))
		order by ck.name, ck.type, pg.rowid`,
		func() (err error) {
			var exe rt.Execute
			dec := gob.NewDecoder(bytes.NewBuffer(prog))
			if e := dec.Decode(&exe); e != nil {
				log.Println(e)
			} else {
				tests = append(tests, check.CheckOutput{name, expect, exe})
			}
			return
		}, &name, &prog, &expect); e != nil {
		err = e
	} else {
		// FIX: we have to cache the statements b/c we cant use them during QueryAll
		for _, t := range tests {
			if e := t.RunTest(run); e != nil {
				err = errutil.New("unexpected failure", t.Name, t.Prog, err)
				break
			}
		}
	}
	return
}
