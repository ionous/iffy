// Generates ephemera from a story file.
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/cmd/import/internal"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var inFile, outFile string
	flag.StringVar(&inFile, "in", "", "input file name (json)")
	flag.StringVar(&outFile, "out", "", "output file name (sqlite3)")
	flag.Parse()
	if e := distill(outFile, inFile); e != nil {
		log.Fatalln(e)
	} else {
		log.Println("Imported", inFile, "into", outFile)
	}
}

func distill(outFile, inFile string) (err error) {
	// fix: write to temp db file then copy the file on success?
	// currently stray files are left hanging around
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if outFile, e := filepath.Abs(outFile); e != nil {
		err = e
	} else if inData, e := readJson(inFile); e != nil {
		err = errutil.New("couldn't read file", inFile, e)
	} else if outDB, e := sql.Open("sqlite3", outFile); e != nil {
		err = errutil.New("couldn't create output file", outFile, e)
	} else {
		defer outDB.Close()
		if e := tables.CreateEphemera(outDB); e != nil {
			err = errutil.New("couldn't create tables", outFile, e)
		} else if e := internal.ImportStory(inFile, inData, outDB); e != nil {
			err = errutil.New("couldn't import story", e)
		}
	}
	return
}

func readJson(filePath string) (ret reader.Map, err error) {
	if fp, e := os.Open(filePath); e != nil {
		err = e
	} else {
		dec := json.NewDecoder(fp)
		if e := dec.Decode(&ret); e != nil && e != io.EOF {
			err = e
		}
	}
	return
}
