package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ionous/iffy/web"
	"github.com/ionous/iffy/web/support"
	_ "github.com/mattn/go-sqlite3"
)

//go:generate templify -p main -o atlas.gen.go atlas.sql
func CreateAtlas(db *sql.DB) error {
	_, e := db.Exec(atlasTemplate())
	return e
}

func main() {
	var testData bool
	flag.BoolVar(&testData, "test", false, "use testdata")
	flag.Parse()
	fileName := flag.Arg(0)
	if len(fileName) == 0 || fileName == "memory" {
		fileName = "file:test.db?cache=shared&mode=memory"
	}
	if db, e := sql.Open("sqlite3", fileName); e != nil {
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
								Gets: func(w http.ResponseWriter) error {
									return kinds(w, db)
								},
							}
						}
						return
					},
					Gets: func(w http.ResponseWriter) error {
						return links.Execute(w, []struct{ Link, Text string }{
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

var links = template.Must(template.New("links").Parse(`
<ul>
{{range .}}
<li><a href="{{.Link}}">{{.Text}}</a></li>
{{end}}
</ul>
`))

func Empty(name string) web.Resource {
	return &web.Wrapper{
		Gets: func(w http.ResponseWriter) error {
			_, e := fmt.Fprintf(w, "No %s to see here", name)
			return e
		},
	}
}
