// Package main for 'asm'.
// Generates a model database from ephemera data.
package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/tables"
)

func main() {
	var inFile, outFile string
	flag.StringVar(&inFile, "in", "", "input file name (sqlite3)")
	flag.StringVar(&outFile, "out", "", "output file name (sqlite3)")
	flag.Parse()
	if e := assemble(outFile, inFile); e != nil {
		log.Fatalln(e)
	} else {
		log.Println("Assembled", inFile, "into", outFile)
	}
}

func assemble(outFile, inFile string) (err error) {
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if outFile, e := filepath.Abs(outFile); e != nil {
		err = e
	} else if e := os.Remove(outFile); e != nil && !os.IsNotExist(e) {
		err = errutil.New("couldn't clean output file", outFile, e)
	} else if db, e := sql.Open("iffy_asm", outFile); e != nil {
		err = errutil.New("couldn't create output file", outFile, e)
	} else {
		defer db.Close()
		//
		if e := tables.CreateModel(db); e != nil {
			err = e // create this in our output db
		} else if e := tables.CreateAssembly(db); e != nil {
			err = e // assembly are temporary tables used for computing the model
		} else if e := func() (err error) {
			// stat fails if there's no such file :(
			ai, _ := os.Stat(inFile)
			bi, _ := os.Stat(outFile)
			if !os.SameFile(ai, bi) {
				s := "attach database '" + inFile + "' as indb;"
				_, err = db.Exec(s)
			}
			return
		}(); e != nil {
			err = errutil.New("error attaching", e, inFile)
		} else {
			var ds assembly.Dilemmas
			if e := assembly.AssembleStory(db, "things", ds.Add); e != nil {
				err = errutil.New("error assembling", e, inFile)
			} else if len(ds) > 0 {
				err = errutil.New("issues assembling", ds.Err())
			}
		}
	}
	return
}
