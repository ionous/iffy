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

func TestGetFieldValues(t *testing.T) {
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
		pathsOfKind := []assembly.TargetField{
			{"K", ""},
			{"A", "K"},
			{"L", "K"},
			{"F", "L,K"},
		}
		kindsOfNoun := []assembly.TargetField{
			{"apple", "K"},
			{"duck", "A"},
			{"toy boat", "L"},
			{"boat", "F"},
		}
		pathsOfNoun := []assembly.TargetField{
			{"apple", "K"},
			{"duck", "A,K"},
			{"toy boat", "L,K"},
			{"boat", "F,L,K"},
		}
		if e := assembly.AddTestHierarchy(m, pathsOfKind); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestFields(m, []assembly.TargetValue{
			{"K", "d", tables.PRIM_DIGI},
			{"K", "t", tables.PRIM_TEXT},
		}); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestNouns(m, kindsOfNoun); e != nil {
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

			// ensure we can ask for object existence
			t.Run("object exists", func(t *testing.T) {
				for _, v := range existence {
					var exists bool
					if e := q.GetField(v.name, object.Exists, &exists); e != nil {
						t.Fatal("existence", v.name, e)
					} else if v.exists != exists {
						t.Fatal("existence", v.name, "wanted", v.exists)
					}
				}
			})
			// ensure queries for kinds work
			t.Run("object kind", func(t *testing.T) {
				for _, v := range kindsOfNoun {
					for i := 0; i < 2; i++ {
						var kind string
						if e := q.GetField(v.Target, object.Kind, &kind); e != nil {
							t.Fatal(e)
						} else if kind != v.Field {
							t.Fatal("mismatch", v.Target, "got:", kind, "expected:", v.Field)
						}
					}
				}
				var kind string
				if e := q.GetField("speedboat", object.Kind, &kind); e == nil {
					t.Fatal("expected error; got", kind)
				}
			})
			// ensure queries for paths work
			t.Run("object kinds", func(t *testing.T) {
				for _, v := range pathsOfNoun {
					for i := 0; i < 2; i++ {
						var path string
						if e := q.GetField(v.Target, object.Kinds, &path); e != nil {
							t.Fatal(e)
						} else if path != v.Field {
							t.Fatal("mismatch", v.Target, "got:", path, "expected:", v.Field)
						}
					}
				}
				var path string
				if e := q.GetField("speedboat", object.Kinds, &path); e == nil {
					t.Fatal("expected error; got", path)
				}
			})
			t.Run("get numbers", func(t *testing.T) {
				for _, v := range numValues {
					for i := 0; i < 2; i++ {
						var num float64
						if e := q.GetField(v.name, "d", &num); e != nil {
							t.Fatal(e)
						} else if num != v.value {
							t.Fatal("mismatch", v.name, num, v.value)
						}
					}
				}
			})
			t.Run("get text", func(t *testing.T) {
				for _, v := range txtValues {
					for i := 0; i < 2; i++ {
						var txt string
						if e := q.GetField(v.name, "t", &txt); e != nil {
							t.Fatal(e)
						} else if txt != v.value {
							t.Fatal("mismatch", v.name, "got:", txt, "expected:", v.value)
						}
					}
				}
			})

		}
	}
}
