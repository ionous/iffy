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

// ex. go run check.go -in /Users/ionous/Documents/Iffy/3ruwyfdnk4umh/play.db
func main() {
	var inFile, testName string
	flag.StringVar(&inFile, "in", "", "input file name (sqlite3)")
	flag.StringVar(&testName, "run", "", "optional specific test ( in camelcase )")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if cnt, e := checkFile(inFile, testName); e != nil {
		log.Fatalln(e)
	} else {
		log.Println("Checked", cnt, inFile)
	}
}

// open db, select tests, de-gob and run them each in turn.
// print the results, only error on critical errors
func checkFile(inFile, testName string) (ret int, err error) {
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
			ret, err = qna.CheckAll(db, testName)
		}
	}
	return
}

func init() {
	iffy.RegisterGobs()
}
