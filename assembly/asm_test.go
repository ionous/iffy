package assembly

import (
	"database/sql"
	"os/user"
	"path"
	"testing"

	"github.com/ionous/errutil"
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

func newAssemblyTest(t *testing.T, useMemory bool) (ret *assemblyTest) {
	var source string
	if useMemory {
		source = memory
	} else if path, e := getPath(t.Name() + ".db"); e != nil {
		t.Fatal(e)
	} else {
		source = path
	}
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		dbq := ephemera.NewDBQueue(db)
		rec := ephemera.NewRecorder(t.Name(), dbq)
		mdl := NewModelerDB(db)
		ret = &assemblyTest{
			T:       t,
			db:      db,
			dbq:     dbq,
			rec:     rec,
			modeler: mdl,
		}
	}
	return
}

type assemblyTest struct {
	*testing.T
	db      *sql.DB
	dbq     *ephemera.DbQueue
	rec     *ephemera.Recorder
	modeler *Modeler
}

func (t *assemblyTest) Close() {
	t.db.Close()
}

type kfp struct{ kind, field, fieldType string }
type pair struct{ key, value string }

// create some fake hierarchy
func fakeHierarchy(w *Modeler, kinds []pair) (err error) {
	for _, p := range kinds {
		if e := w.WriteAncestor(p.key, p.value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// create some fake hierarchy; mdl_field: field, kind, type.
func fakeFields(w *Modeler, kinds []kfp) (err error) {
	for _, p := range kinds {
		if e := w.WriteField(p.field, p.kind, p.fieldType); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
