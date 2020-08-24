package main

import (
	"database/sql"
	"flag"
	"log"
	"path/filepath"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy"
	"github.com/ionous/iffy/qna"
	"github.com/ionous/iffy/tables"
)

func main() {
	var inFile string
	flag.StringVar(&inFile, "in", "", "input file name (sqlite3)")
	flag.Parse()
	if e := checkFile(inFile); e != nil {
		log.Fatalln(e)
	} else {
		log.Println("Checked", inFile)
	}
}

// open db, select tests, de-gob and run them each in turn.
// print the results, only error on critical errors
func checkFile(inFile string) (err error) {
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if db, e := sql.Open("sqlite3", inFile); e != nil {
		err = errutil.New("couldn't create output file", inFile, e)
	} else {
		defer db.Close()
		if tables.CreateRun(db); e != nil {
			err = e
		} else if tables.CreateRunViews(db); e != nil {
			err = e
		} else if e := qna.ActivateDomain(db, "entireGame", true); e != nil {
			err = e
		} else {
			err = qna.CheckAll(db)
		}
	}
	return
}

func init() {
	iffy.RegisterGobs()
}
