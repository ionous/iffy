package story

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
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

// if path is nil, it will use a file db.
func newImportDB(t *testing.T, where string) (ret *sql.DB) {
	var source string
	if len(where) > 0 {
		source = where
	} else if p, e := testdb.PathFromName(t.Name()); e != nil {
		t.Fatal(e)
	} else {
		source = p
	}
	//
	if db, e := sql.Open(tables.DefaultDriver, source); e != nil {
		t.Fatal(e)
	} else {
		t.Log("opened db", source)
		ret = db
	}
	return
}
