package qna

import (
	"testing"

	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/assign"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/tables"
)

func TestGetFieldValues(t *testing.T) {
	db := newQnaDB(t, memory)
	defer db.Close()
	if e := tables.CreateModel(db); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateRunViews(db); e != nil {
		t.Fatal(e)
	} else {
		m := assembly.NewAssembler(db)
		pathsOfKind := []string{
			"Ks", "",
			"As", "Ks",
			"Ls", "Ks",
			"Fs", "Ls,Ks",
		}
		kindsOfNoun := []string{
			"apple", "Ks",
			"duck", "As",
			"toy boat", "Ls",
			"boat", "Fs",
		}
		pathsOfNoun := []string{
			"apple", "Ks",
			"duck", "As,Ks",
			"toy boat", "Ls,Ks",
			"boat", "Fs,Ls,Ks",
		}
		if e := assembly.AddTestHierarchy(m, pathsOfKind...); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestFields(m,
			"Ks", "d", tables.PRIM_DIGI,
			"Ks", "t", tables.PRIM_TEXT,
		); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestNouns(m, kindsOfNoun...); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestStarts(m,
			"apple", "d", 5,
			"duck", "d", 1,
			"toy boat", "t", "boboat",
			"boat", "t", "xyzzy",
		); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestDefaults(m,
			"Ks", "d", 42,
			"As", "t", "chippo",
			"Ls", "t", "weazy",
			"Fs", "d", 13,
		); e != nil {
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
					if p, e := q.GetField(v.name, object.Exists); e != nil {
						t.Fatal("existence", v.name, e)
					} else if exists, e := assign.ToBool(p); e != nil {
						t.Fatal("assign", e)
					} else if v.exists != exists {
						t.Fatal("existence", v.name, "wanted", v.exists)
					}
				}
			})
			// ensure queries for kinds work
			t.Run("object kind", func(t *testing.T) {
				els := kindsOfNoun
				for i, cnt := 0, len(els); i < cnt; i += 2 {
					tgt, field := els[i], els[i+1]
					if p, e := q.GetField(tgt, object.Kind); e != nil {
						t.Fatal(e)
					} else if kind, e := assign.ToString(p); e != nil {
						t.Fatal("assign", e)
					} else if kind != field {
						t.Fatal("mismatch", tgt, "got:", kind, "expected:", field)
					}
				}
				if k, e := q.GetField("speedboat", object.Kind); e == nil {
					t.Fatal("expected error; got", k)
				}
			})
			// ensure queries for paths work
			t.Run("object kinds", func(t *testing.T) {
				els := pathsOfNoun
				for i, cnt := 0, len(els); i < cnt; i += 2 {
					tgt, field := els[i], els[i+1]
					if p, e := q.GetField(tgt, object.Kinds); e != nil {
						t.Fatal(e)
					} else if path, e := assign.ToString(p); e != nil {
						t.Fatal("assign", e)
					} else if path != field {
						t.Fatal("mismatch", tgt, "got:", tgt, "expected:", field)
					}
				}
				if path, e := q.GetField("speedboat", object.Kinds); e == nil {
					t.Fatal("expected error; got", path)
				}
			})
			t.Run("get numbers", func(t *testing.T) {
				for _, v := range numValues {
					for i := 0; i < 2; i++ {
						if p, e := q.GetField(v.name, "d"); e != nil {
							t.Fatal(e)
						} else if num, e := assign.ToFloat(p); e != nil {
							t.Fatal("assign", e)
						} else if num != v.value {
							t.Fatal("mismatch", v.name, num, v.value)
						}
					}
				}
			})
			t.Run("get text", func(t *testing.T) {
				for _, v := range txtValues {
					for i := 0; i < 2; i++ {
						if p, e := q.GetField(v.name, "t"); e != nil {
							t.Fatal(e)
						} else if txt, e := assign.ToString(p); e != nil {
							t.Fatal("assign", e)
						} else if txt != v.value {
							t.Fatal("mismatch", v.name, "got:", txt, "expected:", v.value)
						}
					}
				}
			})
		}
	}
}
