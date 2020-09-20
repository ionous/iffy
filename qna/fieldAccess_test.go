package qna

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/tables"
)

//
func TestFieldAccess(t *testing.T) {
	db := newFieldAccessTest(t, memory)
	defer db.Close()
	q := NewRuntime(db)

	// ensure we can ask for object existence
	t.Run("object exists", func(t *testing.T) {
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
			tgt := v.name
			name := assembly.DomainNameOf("test", tgt)
			if p, e := q.GetField(name, object.Exists); e != nil {
				t.Fatal("existence", name, e)
			} else if exists, e := p.GetBool(nil); e != nil {
				t.Fatal("assign", e)
			} else if v.exists != exists {
				t.Fatal("existence", name, "wanted", v.exists)
			}
		}
	})
	// ensure queries for kinds work
	t.Run("object kind", func(t *testing.T) {
		els := FieldTest.kindsOfNoun
		for i, cnt := 0, len(els); i < cnt; i += 2 {
			tgt, field := els[i], els[i+1]
			name := assembly.DomainNameOf("test", tgt)
			if p, e := q.GetField(name, object.Kind); e != nil {
				t.Fatal(e)
			} else if kind, e := p.GetText(nil); e != nil {
				t.Fatal("assign", e)
			} else if kind != field {
				t.Fatal("mismatch", name, field, "got:", kind, "expected:", field)
			}
		}
		if k, e := q.GetField("speedboat", object.Kind); e == nil {
			t.Fatal("expected error; got", k)
		}
	})
	// ensure queries for paths work
	t.Run("object kinds", func(t *testing.T) {
		els := FieldTest.pathsOfNoun
		for i, cnt := 0, len(els); i < cnt; i += 2 {
			tgt, field := els[i], els[i+1]
			name := assembly.DomainNameOf("test", tgt)
			// asking for "Kinds" should get us the hierarchy
			if p, e := q.GetField(name, object.Kinds); e != nil {
				t.Fatal(e)
			} else if path, e := p.GetText(nil); e != nil {
				t.Fatal("assign", e)
			} else if path != field {
				t.Fatal("mismatch", name, field, "got:", name, "expected:", field)
			}
		}
		if path, e := q.GetField("speedboat", object.Kinds); e == nil {
			t.Fatal("expected error; got", path)
		}
	})
	t.Run("get_text", func(t *testing.T) {
		els := FieldTest.txtValues
		for i, cnt := 0, len(els); i < cnt; i += 3 {
			tgt, field, value := els[i].(string), els[i+1].(string), els[i+2]
			name := assembly.DomainNameOf("test", tgt)
			for i := 0; i < 2; i++ {
				p, e := q.GetField(name, field)
				switch e.(type) {
				default:
					t.Fatal(e)
				case rt.UnknownField:
					if value != nil {
						t.Fatal("got unknown field, but expecting a value")
					}
				case nil:
					if txt, e := p.GetText(nil); e != nil {
						t.Fatal("assign", e)
					} else if txt != value {
						t.Fatalf("mismatch %s.%s got:%q expected:%q", tgt, field, txt, value)
					}
				}
			}
		}
	})
	t.Run("get_numbers", func(t *testing.T) {
		els := FieldTest.numValues
		for i, cnt := 0, len(els); i < cnt; i += 3 {
			tgt, field, value := els[i].(string), els[i+1].(string), els[i+2].(float64)
			name := assembly.DomainNameOf("test", tgt)
			for i := 0; i < 2; i++ {
				if p, e := q.GetField(name, field); e != nil {
					t.Fatal(e)
				} else if num, e := p.GetNumber(nil); e != nil {
					t.Fatal("assign", e)
				} else if num != value {
					t.Fatal("mismatch", name, "have:", num, "want:", value)
				}
			}
		}
	})
	t.Run("get_traits", func(t *testing.T) {
		els := FieldTest.boolValues
		for i, cnt := 0, len(els); i < cnt; i += 2 {
			tgt, csv := els[i].(string), els[i+1].(string)
			name := assembly.DomainNameOf("test", tgt)
			if e := testTraits(q, name, csv); e != nil {
				t.Fatal(e)
			}
		}
	})
	t.Run("change traits", func(t *testing.T) {
		// apple.A had an implicit value of w; change it to "y"
		apple := assembly.DomainNameOf("test", "apple")
		boat := assembly.DomainNameOf("test", "boat")
		toyBoat := assembly.DomainNameOf("test", "toy boat")
		if e := q.SetField(apple, "a", &generic.String{Value: "y"}); e != nil {
			t.Fatal(e)
		} else if v, e := q.GetField(apple, "a"); e != nil {
			t.Fatal(e)
		} else if str, e := v.GetText(q); e != nil {
			t.Fatal(e)
		} else if str != "y" {
			t.Fatal("mismatch", str)
		} else if e := testTraits(q, apple, "y,w,x"); e != nil {
			t.Fatal(e)
		}
		// boat.B has a default value of zz
		if e := q.SetField(boat, "z", &generic.Bool{Value: true}); e != nil {
			t.Fatal(e)
		} else if v, e := q.GetField(boat, "b"); e != nil {
			t.Fatal(e)
		} else if str, e := v.GetText(q); e != nil {
			t.Fatal(e)
		} else if str != "z" {
			t.Fatal("mismatch", str)
		} else if e := testTraits(q, boat, "z, zz"); e != nil {
			t.Fatal(e)
		}
		// toy boat.A has an initial value of y
		if e := q.SetField(toyBoat, "w", &generic.Bool{Value: true}); e != nil {
			t.Fatal(e)
		} else if v, e := q.GetField(toyBoat, "a"); e != nil {
			t.Fatal(e)
		} else if str, e := v.GetText(q); e != nil {
			t.Fatal(e)
		} else if str != "w" {
			t.Fatal("mismatch", str)
		} else if e := testTraits(q, toyBoat, "w,x,y"); e != nil {
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

func testTraits(q *Runner, name, csv string) (err error) {
	traits := strings.Split(csv, ",")
	// the first value in the list of traits is supposed to be true
	for want := true; len(traits) > 0 && err == nil; want = false {
		trait := traits[0]
		traits = traits[1:]
		if p, e := q.GetField(name, trait); e != nil {
			err = errutil.New(e)
		} else if got, e := p.GetBool(nil); e != nil {
			err = errutil.New("assign", e)
		} else if got != want {
			err = errutil.New("mismatch", name, trait, "got:", got, "expected:", want)
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
