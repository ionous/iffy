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
	var tests []checkTest
	var prog []byte
	var name string
	if e := tables.QueryAll(db,
		`select ck.name, pg.bytes
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
				tests = append(tests, checkTest{name, res})
			}
			return
		}, &name, &prog); e != nil {
		err = e
	} else {
		// FIX: we have to cache the statements b/c we cant use them during QueryAll
		for _, test := range tests {
			if e := test.runTest(db, run); e != nil {
				log.Println(e)
			}
		}
	}
	return
}

type checkTest struct {
	name string
	prog check.Testing
}

func (t *checkTest) runTest(db *sql.DB, run rt.Runtime) (err error) {
	name, prog := t.name, t.prog
	if e := ActivateDomain(db, name, true); e != nil {
		err = e
	} else {
		if e := prog.RunTest(run); e != nil {
			err = e
		}
		if e := ActivateDomain(db, name, false); e != nil {
			err = errutil.Append(err, e)
		}
	}

	if err != nil {
		err = errutil.New("unexpected failure", name, prog, err)
	}
	return
}
