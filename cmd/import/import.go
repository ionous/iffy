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
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/ephemera/story"
	"github.com/ionous/iffy/tables"
)

// Import reads a json file (from the composer editor)
// and creates a new sqlite database of "ephemera".
// It uses package export's list of commands for parsing program statements.
func main() {
	var inFile, outFile string
	flag.StringVar(&inFile, "in", "", "input file name (json)")
	flag.StringVar(&outFile, "out", "", "output file name (sqlite3)")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if len(outFile) == 0 {
		dir, _ := filepath.Split(inFile)
		outFile = filepath.Join(dir, "ephemera.db")
	}
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
	} else if e := os.Remove(outFile); e != nil && !os.IsNotExist(e) {
		err = errutil.New("couldn't clean output file", outFile, e)
	} else if outDB, e := sql.Open("sqlite3", outFile); e != nil {
		err = errutil.New("couldn't create output file", outFile, e)
	} else {
		defer outDB.Close()
		if e := tables.CreateEphemera(outDB); e != nil {
			err = errutil.New("couldn't create tables", outFile, e)
		} else if e := story.ImportStory(inFile, outDB, inData); e != nil {
			err = errutil.New("couldn't import story", e)
		}
	}
	return
}

func readJson(filePath string) (ret reader.Map, err error) {
	if f, e := os.Open(filePath); e != nil {
		err = e
	} else {
		defer f.Close()
		dec := json.NewDecoder(f)
		if e := dec.Decode(&ret); e != nil && e != io.EOF {
			err = e
		}
	}
	return
}
