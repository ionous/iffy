package main

import (
	"database/sql"
	"log"
	"os/user"
	"path"

	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/ephemera"
	_ "github.com/mattn/go-sqlite3"
)

func getPath(file string) (ret string, err error) {
	if user, e := user.Current(); e != nil {
		err = e
	} else {
		ret = path.Join(user.HomeDir, file)
	}
	return
}

func main() {
	if path, e := getPath("iffyTest.db"); e != nil {
		log.Fatalln(e)
	} else if db, e := sql.Open("sqlite3", path); e != nil {
		log.Fatalln("db open", e)
	} else {
		defer db.Close()
		q := ephemera.NewDBQueue(db)
		w := assembly.NewModeler(q)
		if e := assembly.DetermineAncestry(w, db, "things"); e != nil {
			log.Fatalln(e)
		} else if e := assembly.DetermineFields(w, db); e != nil {
			log.Fatalln(e)
		} else if e := assembly.DetermineAspects(w, db); e != nil {
			log.Fatalln(e)
		}
		// [-] adds relations between kinds
		// [-] creates instances
		// [-] sets instance properties
		// [-] relates instances
		// [] makes action handlers
		// [] makes event listeners
		// [] computes aliases
		// [] sets up printed name property
		// - backtracing to source:
		// ex. each "important" table entry gets an separate entry pointing back to original source
	}
}
