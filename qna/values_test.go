package qna

import (
	"database/sql"
	"os/user"
	"path"
	"testing"

	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/object"
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

func TestGetObjectValues(t *testing.T) {
	if source, e := sqlFile(t, memory); e != nil {
		t.Fatal(e)
	} else if db, e := sql.Open("sqlite3", source); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateModel(db); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateRunViews(db); e != nil {
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
			{"A", "t", "chippo"},
			{"L", "t", "weazy"},
			{"F", "d", 13},
		}); e != nil {
			t.Fatal(e)
		} else {
			numValues := []struct {
				name  string
				value float64
			}{
				{"apple", 5},
				{"boat", 13},
				{"duck", 1},
				{"toy boat", 42},
			}
			txtValues := []struct {
				name  string
				value string
			}{
				{"apple", ""},
				{"boat", "xyzzy"},
				{"duck", "chippo"},
				{"toy boat", "boboat"},
			}
			existence := []struct {
				name   string
				exists bool
			}{
				{"apple", true},
				{"boat", true},
				{"duck", true},
				{"toy boat", true},
				{"speedboat", false},
			}

			q := NewObjectValues(db)
			for _, v := range existence {
				var exists bool
				if e := q.GetObject(v.name, object.Exists, &exists); e != nil {
					t.Fatal("existence", v.name, e)
				} else if v.exists != exists {
					t.Fatal("existence", v.name, "wanted", v.exists)
				}
			}
			for _, v := range numValues {
				for i := 0; i < 2; i++ {
					var num float64
					if e := q.GetObject(v.name, "d", &num); e != nil {
						t.Fatal(e)
					} else if num != v.value {
						t.Fatal("mismatch", v.name, num, v.value)
					}
				}
			}
			for _, v := range txtValues {
				for i := 0; i < 2; i++ {
					var txt string
					if e := q.GetObject(v.name, "t", &txt); e != nil {
						t.Fatal(e)
					} else if txt != v.value {
						t.Fatal("mismatch", v.name, "got:", txt, "expected:", v.value)
					}
				}
			}
		}
	}
}
