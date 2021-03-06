// Generates ephemera from a story file.
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/ephemera/story"
	"github.com/ionous/iffy/tables"
	"github.com/kr/pretty"
)

// Import reads a json file (from the composer editor)
// and creates a new sqlite database of "ephemera".
// It uses package export's list of commands for parsing program statements.

// ex. go run import.go -in /Users/ionous/Documents/Iffy/stories/shared -out /Users/ionous/Documents/Iffy/scratch/shared/ephemera.db
func main() {
	var inFile, outFile string
	var printStories bool
	flag.StringVar(&inFile, "in", "", "input file or directory name (json)")
	flag.StringVar(&outFile, "out", "", "optional output filename (sqlite3)")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.BoolVar(&printStories, "log", false, "write imported stories to console")
	flag.Parse()
	if len(outFile) == 0 {
		dir, _ := filepath.Split(inFile)
		outFile = filepath.Join(dir, "ephemera.db")
	}
	if xs, e := distill(outFile, inFile); e != nil {
		log.Fatalln(e)
	} else {
		if printStories {
			for _, x := range xs {
				pretty.Println(x)
			}
		}
		log.Println("Imported", inFile, "into", outFile)
	}
}

func distill(outFile, inFile string) (ret []*story.Story, err error) {
	// fix: write to temp db file then copy the file on success?
	// currently stray files are left hanging around
	if inFile, e := filepath.Abs(inFile); e != nil {
		err = e
	} else if outFile, e := filepath.Abs(outFile); e != nil {
		err = e
	} else if inData, e := readPath(inFile); e != nil {
		err = errutil.New("couldn't read file", inFile, e)
	} else if e := os.Remove(outFile); e != nil && !os.IsNotExist(e) {
		err = errutil.New("couldn't clean output file", outFile, e)
	} else {
		// 0755 -> readable by all but only writable by the user
		// 0700 -> read/writable by user
		// 0777 -> ModePerm ... read/writable by all
		os.MkdirAll(path.Dir(outFile), os.ModePerm)
		//
		if outDB, e := sql.Open(tables.DefaultDriver, outFile); e != nil {
			err = errutil.New("couldn't create output file", outFile, e)
		} else {
			var ds reader.Dilemmas
			defer outDB.Close()
			if e := tables.CreateEphemera(outDB); e != nil {
				err = errutil.New("couldn't create tables", outFile, e)
			} else if xs, e := story.ImportStories(inFile, outDB, inData, ds.Report); e != nil {
				err = errutil.New("couldn't import story", e)
			} else {
				ret = xs
				reader.PrintDilemmas(log.Writer(), ds)
			}
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

// read a comma-separated list of files and directories
func readPath(filePaths string) (ret []reader.Map, err error) {
	split := strings.Split(filePaths, ",")
	for _, filePath := range split {
		if info, e := os.Stat(filePath); e != nil {
			err = e
		} else {
			if !info.IsDir() {
				if one, e := readJson(filePath); e != nil {
					err = e
				} else {
					ret = append(ret, one)
				}
			} else {
				if many, e := readDir(filePath); e != nil {
					err = e
				} else {
					ret = append(ret, many...)
				}
			}
		}
	}
	return
}

func readDir(path string) (ret []reader.Map, err error) {
	if !strings.HasSuffix(path, "/") {
		path += "/" // for opening symbolic directories
	}
	err = filepath.Walk(path, func(path string, info os.FileInfo, e error) (err error) {
		if e != nil {
			err = e
		} else if !info.IsDir() && filepath.Ext(path) == ".if" {
			if one, e := readJson(path); e != nil {
				err = errutil.New("error reading", path, e)
			} else {
				ret = append(ret, one)
			}
		}
		return
	})
	return
}
