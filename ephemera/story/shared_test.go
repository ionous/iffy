package story

import (
	"database/sql"
	"os/user"
	"path"
	"strings"
	"testing"

	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/export"
	"github.com/ionous/iffy/tables"
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
	db := newTestDB(t, memory)
	if e := tables.CreateEphemera(db); e != nil {
		t.Fatal("create ephemera", e)
	} else {
		ret = imp.NewImporterDecoder(t.Name(), db, dec)
		retDB = db
	}
	return
}

const memory = "file:test.db?cache=shared&mode=memory"

// if path is nil, it will use a file db.
func newTestDB(t *testing.T, where string) (ret *sql.DB) {
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
		ret = db
	}
	return
}
