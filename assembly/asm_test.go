package assembly

import (
	"database/sql"
	"os/user"
	"path"
	"testing"

	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

func getPath(file string) (ret string, err error) {
	if user, e := user.Current(); e != nil {
		err = e
	} else {
		ret = path.Join(user.HomeDir, file)
	}
	return
}

const memory = "file:test.db?cache=shared&mode=memory"

func newAssemblyTest(t *testing.T, path string) (ret *assemblyTest, err error) {
	ret = &assemblyTest{T: t}
	var source string
	if len(path) > 0 {
		source = path
	} else if p, e := getPath(t.Name() + ".db"); e != nil {
		err = e
	} else {
		source = p
	}
	//
	if err == nil {
		if db, e := sql.Open("sqlite3", source); e != nil {
			err = e
		} else if e := tables.CreateEphemera(db); e != nil {
			err = e
		} else if e := tables.CreateAssembly(db); e != nil {
			err = e
		} else if e := tables.CreateModel(db); e != nil {
			err = e
		} else {
			rec := ephemera.NewRecorder(t.Name(), db)
			mdl := NewModeler(db)
			ret = &assemblyTest{
				T:       t,
				db:      db,
				rec:     rec,
				modeler: mdl,
			}
		}
	}
	return
}

type assemblyTest struct {
	*testing.T
	db      *sql.DB
	rec     *ephemera.Recorder
	modeler *Modeler
}

type kfp struct{ kind, field, fieldType string }
type pair struct{ key, value string }

type triplet struct {
	target, prop string
	value        interface{}
}
