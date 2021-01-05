package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ionous/errutil"
	. "github.com/ionous/iffy/cmd/migrate/internal"
)

func main() {
	paths := "/Users/ionous/Dev/go/src/github.com/ionous/iffy/stories"
	patchPath := "/Users/ionous/Dev/go/src/github.com/ionous/iffy/cmd/migrate/push.patch.js"

	flag.StringVar(&paths, "in", paths, "comma separated input files or directory names")
	flag.StringVar(&patchPath, "patch", patchPath, "patch file")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	//
	var patch Patch
	if e := readJson(patchPath, &patch); e != nil {
		panic(e)
	} else if e := migratePaths(paths, patch); e != nil {
		panic(e)
	}
}

func migratePaths(paths string, patch Migration) (err error) {
	return readPaths(paths, func(path string, doc interface{}) (err error) {
		log.Printf("migrating %q...", path)
		if cnt, e := patch.Migrate(doc); e != nil {
			err = e
		} else if cnt == 0 {
			log.Println("unchanged.")
		} else {
			if f, e := os.Create(path); e != nil {
				err = e
			} else {
				defer f.Close()
				j := json.NewEncoder(f)
				j.SetIndent("", "  ")
				err = j.Encode(doc)
				log.Println("migrated.")
			}
		}
		return //
	})
}

func marshall(i interface{}) (ret string) {
	if b, e := json.MarshalIndent(i, "", " "); e != nil {
		panic(e)
	} else {
		ret = string(b)
	}
	return
}

// read a comma-separated list of files and directories
func readPaths(filePaths string, cb func(path string, data interface{}) error) (err error) {
	split := strings.Split(filePaths, ",")
	for _, path := range split {
		if info, e := os.Stat(path); e != nil {
			err = e
		} else {
			if !info.IsDir() {
				var one interface{}
				if e := readJson(path, &one); e != nil {
					err = e
					break
				} else if e := cb(path, one); e != nil {
					err = e
					break
				}
			} else {
				if !strings.HasSuffix(path, "/") {
					path += "/" // for opening symbolic directories
				}
				if e := filepath.Walk(path, func(path string, info os.FileInfo, e error) (err error) {
					if e != nil {
						err = e
					} else if !info.IsDir() && filepath.Ext(path) == ".if" {
						var one interface{}
						if e := readJson(path, &one); e != nil {
							err = errutil.New("error reading", path, e)
						} else if e := cb(path, one); e != nil {
							err = e
						}
					}
					return // walk
				}); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

func readJson(path string, out interface{}) (err error) {
	// log.Printf("readJson %q", path)
	if f, e := os.Open(path); e != nil {
		err = e
	} else {
		defer f.Close()
		if e := json.NewDecoder(f).Decode(out); e != nil && e != io.EOF {
			err = e
		}
	}
	return
}
