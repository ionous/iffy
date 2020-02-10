package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dbutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/kr/pretty"
)

// TestInitialFieldAssignment to verify initial values for fields can be assigned to instances.
func TestInitialFieldAssignment(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		if e := fakeHierarchy(t.modeler, []pair{
			{"K", ""},
			{"L", "K"},
			{"M", "L,K"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeFields(t.modeler, []kfp{
			{"K", "t", ephemera.PRIM_TEXT},
			{"L", "d", ephemera.PRIM_DIGI},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeNouns(t.modeler, []pair{
			{"apple", "K"},
			{"pear", "L"},
			{"machine gun", "M"},
			{"gun", "M"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addValues(t.rec, []triplet{
			{"apple", "t", "some text"},
			{"pear", "d", 123},
			{"machine", "d", 321},
			{"gun", "t", "more text"},
		}); e != nil {
			t.Fatal(e)
		} else if e := determineInitialFields(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchValues(t.db, []triplet{
			{"apple", "t", "some text"},
			{"gun", "t", "more text"},
			{"machine gun", "d", int64(321)},
			{"pear", "d", int64(123)}, // int64, re: go's default scanner.
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// TestInitialTraitAssignments to verify default traits can be assigned to kinds.
func TestInitialTraitAssignment(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		//
		if e := fakeHierarchy(t.modeler, []pair{
			{"K", ""},
			{"L", "K"},
			{"M", "L,K"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeFields(t.modeler, []kfp{
			{"K", "A", ephemera.PRIM_ASPECT},
			{"L", "B", ephemera.PRIM_ASPECT},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeTraits(t.modeler, []pair{
			{"A", "w"}, {"A", "x"}, {"A", "y"},
			{"B", "z"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeNouns(t.modeler, []pair{
			{"apple", "K"},
			{"pear", "L"},
			{"machine gun", "M"},
			{"gun", "M"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addValues(t.rec, []triplet{
			{"apple", "A", "y"},
			{"pear", "x", true},
			{"machine", "w", true},
			{"gun", "z", true},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineValues(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchValues(t.db, []triplet{
			{"apple", "A", "y"},
			{"gun", "B", "z"},
			{"machine gun", "A", "w"},
			{"pear", "A", "x"},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// match generated model defaults
func matchValues(db *sql.DB, want []triplet) (err error) {
	var curr triplet
	var have []triplet
	if e := dbutil.QueryAll(db,
		`select noun, field, value 
			from mdl_start
			order by noun, field, value`,
		func() (err error) {
			have = append(have, curr)
			return
		},
		&curr.target, &curr.prop, &curr.value); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch",
			"have:", pretty.Sprint(have),
			"want:", pretty.Sprint(want))
	}
	return
}

// eph_value: fake noun, prop, value
// prop: k, f, v
func addValues(rec *ephemera.Recorder, vals []triplet) (err error) {
	for _, v := range vals {
		noun := rec.Named(ephemera.NAMED_NOUN, v.target, "test")
		prop := rec.Named(ephemera.NAMED_PROPERTY, v.prop, "test")
		value := v.value
		rec.NewValue(noun, prop, value)
	}
	return
}

// mdl_noun:  kind, ret id noun
// mdl_name: noun, name, rank
func fakeNouns(m *Modeler, els []pair) (err error) {
	for _, el := range els {
		noun, kind := el.key, el.value
		if e := m.WriteNounWithNames(noun, kind); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
