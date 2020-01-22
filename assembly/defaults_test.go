package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/kr/pretty"
	_ "github.com/mattn/go-sqlite3"
)

func TestDefaultsAssigment(test *testing.T) {
	t := newAssemblyTest(test, true)
	defer t.Close()
	//
	if e := fakeHierarchy(t.modeler, []pair{
		{"T", ""},
		{"P", "T"},
		{"C", "P,T"},
	}); e != nil {
		t.Fatal(e)
	} else if e := fakeFields(t.modeler, []kfp{
		{"T", "d", ephemera.PRIM_DIGI},
		{"T", "t", ephemera.PRIM_TEXT},
		{"T", "t2", ephemera.PRIM_TEXT},
		{"P", "p", ephemera.PRIM_TEXT},
		{"C", "c", ephemera.PRIM_TEXT},
	}); e != nil {
		t.Fatal(e)
	} else if e := writeDefaults(t.rec, []defaultValue{
		{"T", "t", "some text"},
		{"P", "t2", "other text"},
		{"C", "c", "c text"},
		{"C", "d", 123},
	}); e != nil {
		t.Fatal(e)
	} else if e := DetermineDefaults(t.modeler, t.db); e != nil {
		t.Fatal(e)
	} else if e := matchDefaults(t.db, []defaultValue{
		{"C", "c", "c text"},
		// default scanner uses https://golang.org/pkg/database/sql/#Scanner
		{"T", "d", int64(123)}, //
		{"T", "t", "some text"},
		{"T", "t2", "other text"},
	}); e != nil {
		t.Fatal(e)
	}
}
func TestDefaultsMissingField(t *testing.T) {
}

func TestDefaultsConflictingValues(t *testing.T) {
}

func TestDefaultsInvalidType(t *testing.T) {
	//- for now, we only allow text and number [ text and digi ]
	// - later we could add ambiguity for conversion [ 4 -> "4" ]
}

type defaultInfo struct {
	idModelField int64
	value        interface{}
}
type defaultValue struct {
	kind, field string
	value       interface{}
}

// match generated model defaults
func matchDefaults(db *sql.DB, want []defaultValue) (err error) {
	var curr defaultValue
	var have []defaultValue
	if e := dbutil.QueryAll(db, `select kind,field,value 
			from mdl_default d, mdl_field f 
			on d.idModelField = f.rowid
			order by kind, field, value`,
		func() (err error) {
			have = append(have, curr)
			return
		},
		&curr.kind, &curr.field, &curr.value); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch", "have:", pretty.Sprint(have), "want:", pretty.Sprint(want))
	}
	return
}

// write ephemera describing some initial values
func writeDefaults(rec *ephemera.Recorder, defaults []defaultValue) (err error) {
	for _, el := range defaults {
		namedKind := rec.Named(ephemera.NAMED_KIND, el.kind, "test")
		namedField := rec.Named(ephemera.NAMED_FIELD, el.field, "test")
		rec.NewDefault(namedKind, namedField, el.value)
	}
	return
}
