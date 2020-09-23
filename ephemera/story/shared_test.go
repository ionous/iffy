package story

import (
	"database/sql"
	"os/user"
	"path"
	"strings"
	"testing"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/tables"
)

func lines(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}

func newTestImporter(t *testing.T, where string) (ret *Importer, retDB *sql.DB) {
	return newTestImporterDecoder(t, nil, where)
}

func newTestDecoder(t *testing.T, where string) (ret *Importer, retDB *sql.DB) {
	iffy.RegisterGobs()
	//
	dec := decode.NewDecoder()
	dec.AddDefaultCallbacks(core.Slats)
	return newTestImporterDecoder(t, dec, where)
}

func newTestImporterDecoder(t *testing.T, dec *decode.Decoder, where string) (ret *Importer, retDB *sql.DB) {
	db := newImportDB(t, where)
	if e := tables.CreateEphemera(db); e != nil {
		t.Fatal("create ephemera", e)
	} else {
		ret = NewImporterDecoder(t.Name(), db, dec)
		retDB = db
	}
	return
}

const memory = "file:test.db?cache=shared&mode=memory"

// if path is nil, it will use a file db.
func newImportDB(t *testing.T, where string) (ret *sql.DB) {
	var source string
	if len(where) > 0 {
		source = where
	} else if user, e := user.Current(); e != nil {
		t.Fatal(e)
	} else {
		source = path.Join(user.HomeDir, t.Name()+".db")
	}
	//
	if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else {
		t.Log("opened db", source)
		ret = db
	}
	return
}
