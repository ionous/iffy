package qna

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"log"

	"github.com/ionous/errutil"

	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/tables"
)

// CheckAll tests stored in the passed db.
// It logs the results of running the tests, and only returns error on critical errors.
func CheckAll(db *sql.DB, actuallyJustThisOne string) (ret int, err error) {
	run := NewRuntime(db)
	var name string
	var prog []byte
	var tests []check.CheckOutput
	if e := tables.QueryAll(db,
		`select name, bytes 
		from mdl_prog pg 
		where type='CheckOutput'
		order by name`,
		func() (err error) {
			if len(actuallyJustThisOne) == 0 || actuallyJustThisOne == name {
				var curr check.CheckOutput
				dec := gob.NewDecoder(bytes.NewBuffer(prog))
				if e := dec.Decode(&curr); e != nil {
					log.Println(e)
				} else {
					tests = append(tests, curr)
				}
			}
			return nil // we log rather than return error
		}, &name, &prog); e != nil {
		err = e
	} else if len(tests) == 0 {
		err = errutil.New("no matching tests found")
	} else {
		// FIX: we have to cache the statements b/c we cant use them during QueryAll
		for _, t := range tests {
			if e := t.RunTest(run); e != nil {
				err = e
				break
			}
			ret++
		}
	}
	return
}
