package story

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

func TestImportStory(t *testing.T) {
	const memory = "file:test.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal("db open", e)
	} else {
		defer db.Close()
		if e := tables.CreateEphemera(db); e != nil {
			t.Fatal("create ephemera", e)
		}
		var in reader.Map
		if e := json.Unmarshal([]byte(debug.Blob), &in); e != nil {
			t.Fatal("read json", e)
		} else if e := ImportStory(t.Name(), in, db); e != nil {
			t.Fatal("import", e)
		} else {
			t.Log("ok")
		}
	}
}
