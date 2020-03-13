package main

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"flag"
	"log"
	"path/filepath"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/check"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
	"github.com/ionous/iffy/tables"
)

func main() {
	var inFile string
	flag.StringVar(&inFile, "in", "", "input file name (sqlite3)")
	flag.Parse()
	if e := playFile(inFile); e != nil {
		log.Fatalln(e)
	}
}

// open db, select tests, de-gob and run them each in turn.
// print the results, only error on critical errors
func playFile(inFile string) (err error) {
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if db, e := sql.Open("sqlite3", inFile); e != nil {
		err = errutil.New("couldn't create output file", inFile, e)
	} else {
		defer db.Close()
		var prog []byte
		if e := tables.QueryAll(db,
			//select ck.name, pg.type, pg.bytes, ck.expect
			`select pg.bytes
		from mdl_check as ck
		join mdl_prog pg
			on (pg.rowid = ck.idProg)
		order by ck.name`,
			func() (err error) {
				var res check.Test
				dec := gob.NewDecoder(bytes.NewBuffer(prog))
				if e := dec.Decode(&res); e != nil {
					log.Println(e)
				} else if e := runTest(&res); e != nil {
					log.Println(e)
				}
				return
			}, &prog); e != nil {
			err = e
		}
	}
	return
}

func runTest(prog rt.Execute) (err error) {
	if run, e := rtm.New(nil).Rtm(); e != nil {
		err = e
	} else if e := prog.Execute(run); e != nil {
		err = e
	}
	return
}

func init() {
	registerGob()
}

// duplicated in import/internal --
// where should t live?
var registeredGob = false

func registerGob() {
	if !registeredGob {
		for _, t := range export.Runs {
			gob.Register(t.Type)
		}
		registeredGob = true
	}
}
