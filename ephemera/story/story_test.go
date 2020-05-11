package story

import (
	"encoding/json"
	"testing"

	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/ephemera/reader"
)

func TestImportStory(t *testing.T) {
	db := newTestDB(t, memory)
	defer db.Close()
	//
	var in reader.Map
	if e := json.Unmarshal([]byte(debug.Blob), &in); e != nil {
		t.Fatal("read json", e)
	} else if e := ImportStory(t.Name(), in, db); e != nil {
		t.Fatal("import", e)
	} else {
		t.Log("ok")
	}
}
