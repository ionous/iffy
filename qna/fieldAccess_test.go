package qna

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/object"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
)

//
func TestFieldAccess(t *testing.T) {
	db := newFieldAccessTest(t, testdb.Memory)
	defer db.Close()
	q := NewRuntime(db)

	// ensure we can ask for object existence
	t.Run("object_exists", func(t *testing.T) {
		// whether a name exists
		existence := []struct {
			name   string
			exists bool
		}{
			{"apple", true},
			{"boat", true},
			{"duck", true},
			{"toy boat", true},
			{"speedboat", false}, // no such noun
		}
		els := existence
		for _, v := range els {
			name := v.name
			switch _, e := q.GetField(object.Value, name); e.(type) {
			default:
				t.Fatal("assign", e)

			case g.UnknownObject:
				if v.exists {
					t.Fatal("wanted to exist and it doesnt", name)
				}
			case nil:
				if !v.exists {
					t.Fatal("didnt want to exist and it does", name)
				}
			}
		}
	})
	// ensure queries for kinds work
	t.Run("object_kind", func(t *testing.T) {
		els := FieldTest.kindsOfNoun
		for i, cnt := 0, len(els); i < cnt; i += 2 {
			name, field := els[i], els[i+1]
			if obj, e := q.GetField(object.Value, name); e != nil {
				t.Fatal(e)
			} else if p, e := obj.FieldByName(object.Kind); e != nil {
				t.Fatal(e)
			} else if kind := p.String(); kind != field {
				t.Fatal("mismatch", name, field, "got:", kind, "expected:", field)
			}
		}
	})
	// ensure queries for paths work
	t.Run("object_kinds", func(t *testing.T) {
		els := FieldTest.pathsOfNoun
		for i, cnt := 0, len(els); i < cnt; i += 2 {
			name, field := els[i], els[i+1]
			// asking for "Kinds" should get us the hierarchy
			if obj, e := q.GetField(object.Value, name); e != nil {
				t.Fatal(e)
			} else if p, e := obj.FieldByName(object.Kinds); e != nil {
				t.Fatal(e)
			} else if path := p.String(); path != field {
				t.Fatal("mismatch", name, field, "got:", name, "expected:", field)
			}
		}
	})
	t.Run("get_text", func(t *testing.T) {
		els := FieldTest.txtValues
		for i, cnt := 0, len(els); i < cnt; i += 3 {
			name, field, value := els[i].(string), els[i+1].(string), els[i+2]
			for n := 0; n < 2; n++ {
				if obj, e := q.GetField(object.Value, name); e != nil {
					t.Fatal(e)
				} else {
					p, e := obj.FieldByName(field)
					switch e.(type) {
					default:
						t.Fatal(e)
					case g.UnknownField:
						if value != nil {
							t.Fatal("got unknown field, but expecting a value")
						}
					case nil:
						if p == nil {
							t.Fatal("value and error are both nil for", name, field)
						} else if txt := p.String(); txt != value {
							t.Fatalf("mismatch %s.%s got:%q expected:%q", name, field, txt, value)
						}
					}
				}
			}
		}
	})
	t.Run("get_numbers", func(t *testing.T) {
		els := FieldTest.numValues
		for i, cnt := 0, len(els); i < cnt; i += 3 {
			name, field, value := els[i].(string), els[i+1].(string), els[i+2].(float64)
			for i := 0; i < 2; i++ {
				if obj, e := q.GetField(object.Value, name); e != nil {
					t.Fatal(e)
				} else if p, e := obj.FieldByName(field); e != nil {
					t.Fatal(e)
				} else if num := p.Float(); num != value {
					t.Fatal("mismatch", name, "have:", num, "want:", value)
				}
			}
		}
	})
	t.Run("get_traits", func(t *testing.T) {
		els := FieldTest.boolValues
		for i, cnt := 0, len(els); i < cnt; i += 2 {
			name, csv := els[i].(string), els[i+1].(string)
			if obj, e := q.GetField(object.Value, name); e != nil {
				t.Fatal(e)
			} else if e := testTraits(obj, csv); e != nil {
				t.Fatal(e)
			}
		}
	})
	t.Run("change_traits", func(t *testing.T) {
		// apple.A had an implicit value of w; change it to "y"
		if apple, e := q.GetField(object.Value, "apple"); e != nil {
			t.Fatal(e)
		} else if e := apple.SetFieldByName("a", g.StringOf("y")); e != nil {
			t.Fatal(e)
		} else if v, e := apple.FieldByName("a"); e != nil {
			t.Fatal(e)
		} else if str := v.String(); str != "y" {
			t.Fatal("mismatch", str)
		} else if e := testTraits(apple, "y,w,x"); e != nil {
			t.Fatal(e)
		}
		// b is an aspect with traits "z" and "zz"
		// boat.B has a default value of zz
		if boat, e := q.GetField(object.Value, "boat"); e != nil {
			t.Fatal(e)
		} else if e := boat.SetFieldByName("z", g.BoolOf(true)); e != nil {
			t.Fatal(e)
		} else if v, e := boat.FieldByName("b"); e != nil {
			t.Fatal(e)
		} else if str := v.String(); str != "z" {
			t.Fatal("mismatch", str)
		} else if e := testTraits(boat, "z, zz"); e != nil {
			t.Fatal(e)
		}
		// toy boat.A has an initial value of y
		if toyBoat, e := q.GetField(object.Value, "toy_boat"); e != nil {
			t.Fatal(e)
		} else if e := toyBoat.SetFieldByName("w", g.BoolOf(true)); e != nil {
			t.Fatal(e)
		} else if v, e := toyBoat.FieldByName("a"); e != nil {
			t.Fatal(e)
		} else if str := v.String(); str != "w" {
			t.Fatal("mismatch", str)
		} else if e := testTraits(toyBoat, "w,x,y"); e != nil {
			t.Fatal(e)
		}
	})
}

func newFieldAccessTest(t *testing.T, dbloc string) (ret *sql.DB) {
	db := newQnaDB(t, dbloc)
	if e := tables.CreateModel(db); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateRun(db); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateRunViews(db); e != nil {
		t.Fatal(e)
	} else {
		m := assembly.NewAssembler(db)
		if e := assembly.AddTestHierarchy(m, FieldTest.pathsOfKind...); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestFields(m, FieldTest.fields...); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestTraits(m, FieldTest.traits...); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestStarts(m, FieldTest.startingValues...); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestNouns(m, FieldTest.kindsOfNoun...); e != nil {
			t.Fatal(e)
		} else if e := assembly.AddTestDefaults(m, FieldTest.defaultValues...); e != nil {
			t.Fatal(e)
		} else if e := ActivateDomain(db, "test", true); e != nil {
			t.Fatal(e)
		} else {
			ret = db
		}
	}
	return
}

func testTraits(obj g.Value, csv string) (err error) {
	traits := strings.Split(csv, ",")
	// the first value in the list of traits is supposed to be true
	for want := true; len(traits) > 0 && err == nil; want = false {
		trait := traits[0]
		traits = traits[1:]
		if p, e := obj.FieldByName(trait); e != nil {
			err = errutil.New(e)
		} else if got := p.Bool(); got != want {
			err = errutil.New("mismatch", trait, "got:", got, "expected:", want)
		}
	}
	return
}

var FieldTest = struct {
	// kind hierarchy
	pathsOfKind,
	// parents of nouns
	kindsOfNoun,
	// noun hierarchy
	pathsOfNoun,
	// kind, field, type
	fields,
	// aspect, trait pairs
	traits []string
	// noun, field, value triplets
	defaultValues, startingValues,
	// computed noun, field, text value triplets
	txtValues,
	// computed noun, field, num value triplets
	numValues,
	boolValues []interface{}
}{
	/* pathsOfKind*/ []string{
		"Ks", "",
		"Js", "Ks",
		"Ls", "Ks",
		"Fs", "Ls,Ks",
	},
	/*kindsOfNoun*/ []string{
		"apple", "Ks",
		"duck", "Js",
		"toy boat", "Ls",
		"boat", "Fs",
	},
	/*pathsOfNoun*/ []string{
		"apple", "Ks",
		"duck", "Js,Ks",
		"toy boat", "Ls,Ks",
		"boat", "Fs,Ls,Ks",
	},
	/*fields*/ []string{
		"Ks", "d", tables.PRIM_DIGI,
		"Ks", "t", tables.PRIM_TEXT,
		"Ks", "a", tables.PRIM_ASPECT,
		"Ls", "b", tables.PRIM_ASPECT,
	},
	/*traits*/ []string{
		"a", "w",
		"a", "x",
		"a", "y",
		"b", "z",
		"b", "zz",
	},
	/*default values*/ []interface{}{
		"Ks", "d", 42,
		"Js", "t", "chippo",
		"Ls", "t", "weazy",
		"Fs", "d", 13,
		"Fs", "b", "zz",
		"Ls", "a", "x",
	},
	/*starting values*/ []interface{}{
		"apple", "d", 5,
		"duck", "d", 1,
		"toy boat", "t", "boboat",
		"boat", "t", "xyzzy",
		"toy boat", "a", "y",
	},
	/*txtValues*/ []interface{}{
		"apple", "t", "",
		"boat", "t", "xyzzy",
		"duck", "t", "chippo",
		"toy boat", "t", "boboat",
		//
		"apple" /*   */, "a", "w",
		"duck" /*    */, "a", "w",
		"toy boat" /**/, "a", "y",
		"boat" /* */, "a", "x",
		//
		"toy boat" /**/, "b", "z",
		"boat" /* */, "b", "zz",

		// asking for an improper or invalid aspect returns nothing
		// fix? should it return or log error instead?
		"apple" /*   */, "b", nil,
		"boat" /*   */, "G", nil,
	},
	/*numValues*/ []interface{}{
		"apple", "d", 5.0,
		"boat", "d", 13.0,
		"duck", "d", 1.0,
		"toy boat", "d", 42.0,
	},
	// noun, truth values. the first comma separated value is true, the rest false.
	/*boolValues*/ []interface{}{
		"apple", "w,x,y",
		"duck", "w,x,y",
		//
		"toy boat", "y,w,x",
		"toy boat", "z,zz",
		//
		"boat", "x,w,y",
		"boat", "zz,z",
	},
}
