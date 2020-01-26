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

const memory = "file:test.db?cache=shared&mode=memory"

func newAssemblyTest(t *testing.T, path string) (ret *assemblyTest, err error) {
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

type prop struct {
	owner, prop string
	value       interface{}
}

// create some fake hierarchy
func fakeHierarchy(m *Modeler, kinds []pair) (err error) {
	for _, p := range kinds {
		if e := m.WriteAncestor(p.key, p.value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// create some fake hierarchy; mdl_field: field, kind, type.
func fakeFields(m *Modeler, kinds []kfp) (err error) {
	for _, p := range kinds {
		if e := m.WriteField(p.field, p.kind, p.fieldType); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
