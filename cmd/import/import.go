// Generates ephemera from a story file.
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"github.com/ionous/iffy/cmd/import/internal"
	"github.com/ionous/iffy/ephemera/reader"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var inFile, outFile string
	flag.StringVar(&inFile, "in", "", "input file name (json)")
	flag.StringVar(&outFile, "out", "", "output file name (sqlite3)")
	flag.Parse()

	if in, e := readJson(inFile); e != nil {
		log.Fatalln("couldn't read file", inFile, e)
	} else if db, e := sql.Open("sqlite3", outFile); e != nil {
		log.Fatalln("couldn't create output file", outFile, e)
	} else {
		defer db.Close()
		if e := internal.ImportStory(in, db); e != nil {
			log.Fatalln("couldn't import story", e)
		}
		println("Imported", inFile, "into", outFile)
	}
}

func readJson(filePath string) (ret reader.Map, err error) {
	if fp, e := os.Open(filePath); e != nil {
		err = e
	} else {
		dec := json.NewDecoder(fp)
		if e := dec.Decode(&ret); e != nil && e == io.EOF {
			err = e
		}
	}
	return
}
