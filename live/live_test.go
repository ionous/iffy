package live

import (
	"database/sql"
	"os/user"
	"path"
	"testing"

	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/tables"
)

const memory = "file:test.db?cache=shared&mode=memory"

func getPath(file string) (ret string, err error) {
	if user, e := user.Current(); e != nil {
		err = e
	} else {
		ret = path.Join(user.HomeDir, file)
	}
	return
}

func sqlFile(t *testing.T, path string) (ret string, err error) {
	if len(path) > 0 {
		ret = path
	} else if p, e := getPath(t.Name() + ".db"); e != nil {
		err = e
	} else {
		ret = p
	}
	return
}

func TestGetField(t *testing.T) {

	if source, e := sqlFile(t, ""); e != nil {
		t.Fatal(e)
	} else if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateModel(db); e != nil {
		t.Fatal(e)
	} else {
		m := assembly.NewModeler(db)

		if e := assembly.AddTestHierarchy(m, []assembly.TargetField{
			{"K", ""},
			{"A", "K"},
			{"L", "K"},
			{"F", "L,K"},
		}); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestFields(m, []assembly.TargetValue{
			{"K", "d", tables.PRIM_DIGI},
			{"K", "t", tables.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestNouns(m, []assembly.TargetField{
			{"apple", "K"},
			{"duck", "A"},
			{"toy boat", "L"},
			{"boat", "F"},
		}); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestStarts(m, []assembly.TargetValue{
			{"apple", "d", 5},
			{"duck", "d", 1},
			{"toy boat", "t", "boboat"},
			{"boat", "t", "xyzzy"},
		}); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestDefaults(m, []assembly.TargetValue{
			{"K", "d", 42},
			{"F", "d", 13},
			{"L", "t", "weazy"},
		}); e != nil {
			t.Fatal(e)
		}
	}
}
