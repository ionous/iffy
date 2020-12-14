package story

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
)

func lines(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}

// for tests where we need a default decoder to read json
func newImporter(t *testing.T, where string) (ret *Importer, retDB *sql.DB) {
	db := newImportDB(t, where)
	if e := tables.CreateEphemera(db); e != nil {
		t.Fatal("create ephemera", e)
	} else {
		iffy.RegisterGobs()
		dec := decode.NewDecoderReporter(t.Name(), func(pos reader.Position, err error) {
			t.Errorf("%s at %s", err, pos)
		})
		k := NewImporterDecoder(t.Name(), db, dec)
		dec.AddDefaultCallbacks(core.Slats)
		k.AddModel(Model)
		ret, retDB = k, db
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
