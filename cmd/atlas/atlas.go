package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/context"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/web"
	"github.com/ionous/iffy/web/support"
	_ "github.com/mattn/go-sqlite3"
)

//go:generate templify -p main -o atlas.gen.go atlas.sql
func CreateAtlas(db *sql.DB) (err error) {
	if _, e := db.Exec(atlasTemplate()); e != nil {
		err = errutil.New("CreateAtlas:", e)
	}
	return
}

func main() {
	var testData bool
	flag.BoolVar(&testData, "test", false, "use testdata")
	flag.Parse()
	fileName := flag.Arg(0)
	if len(fileName) == 0 || fileName == "memory" {
		fileName = "file:test.db?cache=shared&mode=memory"
	}
	if db, e := sql.Open(tables.DefaultDriver, fileName); e != nil {
		log.Fatalln("db open", e)
	} else {
		if !testData {
			panic("unsupported")
		} else if e := createTestData(db); e != nil {
			log.Fatal(e)
		}

		m := http.NewServeMux()
		m.HandleFunc("/atlas/", web.HandleResource(Atlas(db)))
		go support.OpenBrowser("http://localhost:8080/atlas/")
		log.Fatal(http.ListenAndServe(":8080", m))
	}
}

func Atlas(db *sql.DB) web.Resource {
	return &web.Wrapper{
		Finds: func(name string) (ret web.Resource) {
			switch name {
			case "atlas":
				ret = &web.Wrapper{
					Finds: func(name string) (ret web.Resource) {
						switch name {
						case "nouns":
							ret = Empty(name)
						case "kinds":
							return &web.Wrapper{
								Gets: func(ctx context.Context, w http.ResponseWriter) error {
									return listOfKinds(w, db)
								},
							}
						}
						return
					},
					Gets: func(ctx context.Context, w http.ResponseWriter) error {
						return templates.ExecuteTemplate(w, "links", []struct{ Link, Text string }{
							{"/atlas/kinds/", "kinds"},
							{"/atlas/nouns/", "nouns"},
						})
					},
				}
			}
			return
		},
	}
}

func Empty(name string) web.Resource {
	return &web.Wrapper{
		Gets: func(ctx context.Context, w http.ResponseWriter) error {
			_, e := fmt.Fprintf(w, "No %s to see here", name)
			return e
		},
	}
}
