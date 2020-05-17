package main

import (
	"database/sql"
	"encoding/gob"
	"flag"
	"log"
	"path/filepath"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/qna"
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
		err = qna.CheckAll(db)
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
		export.Register(gob.Register)
		registeredGob = true
	}
}
