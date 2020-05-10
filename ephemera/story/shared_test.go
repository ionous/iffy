package story

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/tables"

	_ "github.com/mattn/go-sqlite3"
)

func lines(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}

func newTestImporter(t *testing.T) (ret *imp.Porter, retDB *sql.DB) {
	return newTestImporterDecoder(t, nil)
}

func newTestDecoder(t *testing.T) (ret *imp.Porter, retDB *sql.DB) {
	dec := decode.NewDecoder()
	dec.AddDefaultCallbacks(export.Slats)
	return newTestImporterDecoder(t, dec)
}

func newTestImporterDecoder(t *testing.T, dec *decode.Decoder) (ret *imp.Porter, retDB *sql.DB) {
	const path = "file:test.db?cache=shared&mode=memory"
	// if path, e := getPath(t.Name() + ".db"); e != nil {
	// 	t.Fatal(e)
	// } else
	if db, e := sql.Open("sqlite3", path); e != nil {
		t.Fatal("db open", e)
	} else {
		if e := tables.CreateEphemera(db); e != nil {
			t.Fatal("create ephemera", e)
		} else {
			ret = imp.NewImporterDecoder(t.Name(), db, dec)
			retDB = db
		}
	}
	return
}
