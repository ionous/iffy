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

type triplet struct {
	target, prop string
	value        interface{}
}

// create some fake model hierarchy
func fakeHierarchy(m *Modeler, kinds []pair) (err error) {
	for _, p := range kinds {
		if e := m.WriteAncestor(p.key, p.value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// create some fake model hierarchy; mdl_field: field, kind, type.
func fakeFields(m *Modeler, kinds []kfp) (err error) {
	for _, p := range kinds {
		if e := m.WriteField(p.field, p.kind, p.fieldType); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// write aspect, trait pairs
func fakeTraits(m *Modeler, traits []pair) (err error) {
	for _, t := range traits {
		// rank is not set yet, see DetermineAspects
		if e := m.WriteTrait(t.key, t.value, 0); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// write kind, aspect pairs
func fakeAspects(m *Modeler, kindAspects []pair) (err error) {
	for _, t := range kindAspects {
		if e := m.WriteAspect(t.key, t.value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
