package qna

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"log"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/tables"
)

// CheckAll tests stored in the passed db.
// It logs the results of running the tests, and only returns error on critical errors.
func CheckAll(db *sql.DB) (err error) {
	run := NewRuntime(db)
	var prog []byte
	var tests []check.Testing
	if e := tables.QueryAll(db,
		`select pg.bytes
		from mdl_check as ck
		join mdl_prog pg
			on (pg.rowid = ck.idProg)
		order by ck.name`,
		func() (err error) {
			var res check.Testing
			dec := gob.NewDecoder(bytes.NewBuffer(prog))
			if e := dec.Decode(&res); e != nil {
				log.Println(e)
			} else {
				tests = append(tests, res)
			}
			return
		}, &prog); e != nil {
		err = e
	} else {
		// FIX: we have to pre-prepare the statements we want otherwise we cant use them during the select loop
		for _, test := range tests {
			if e := runTest(run, test); e != nil {
				log.Println(e)
			}
		}
	}
	return
}

func runTest(run rt.Runtime, prog check.Testing) (err error) {
	if e := prog.RunTest(run); e != nil {
		err = e
	} else if e != nil {
		err = errutil.New("unexpected failure", prog, e)
	}
	return
}
