package assembly

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
	"github.com/kr/pretty"
)

// TestInitialFieldAssignment to verify initial values for fields can be assigned to instances.
func TestInitialFieldAssignment(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		if e := AddTestHierarchy(t.modeler, []TargetField{
			{"K", ""},
			{"L", "K"},
			{"M", "L,K"},
		}); e != nil {
			t.Fatal(e)
		} else if e := AddTestFields(t.modeler, []TargetValue{
			{"K", "t", tables.PRIM_TEXT},
			{"L", "d", tables.PRIM_DIGI},
		}); e != nil {
			t.Fatal(e)
		} else if e := AddTestNouns(t.modeler, []TargetField{
			{"apple", "K"},
			{"pear", "L"},
			{"toy boat", "M"},
			{"boat", "M"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addValues(t.rec, []TargetValue{
			{"apple", "t", "some text"},
			{"pear", "d", 123},
			{"toy", "d", 321},
			{"boat", "t", "more text"},
		}); e != nil {
			t.Fatal(e)
		} else if e := determineInitialFields(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchValues(t.db, []TargetValue{
			{"apple", "t", "some text"},
			{"boat", "t", "more text"},
			{"pear", "d", int64(123)}, // int64, re: go's default scanner.
			{"toy boat", "d", int64(321)},
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
		if e := AddTestHierarchy(t.modeler, []TargetField{
			{"K", ""},
			{"L", "K"},
			{"M", "L,K"},
		}); e != nil {
			t.Fatal(e)
		} else if e := AddTestFields(t.modeler, []TargetValue{
			{"K", "A", tables.PRIM_ASPECT},
			{"L", "B", tables.PRIM_ASPECT},
		}); e != nil {
			t.Fatal(e)
		} else if e := AddTestTraits(t.modeler, []TargetField{
			{"A", "w"}, {"A", "x"}, {"A", "y"},
			{"B", "z"},
		}); e != nil {
			t.Fatal(e)
		} else if e := AddTestNouns(t.modeler, []TargetField{
			{"apple", "K"},
			{"pear", "L"},
			{"toy boat", "M"},
			{"boat", "M"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addValues(t.rec, []TargetValue{
			{"apple", "A", "y"},
			{"pear", "x", true},
			{"toy", "w", true},
			{"boat", "z", true},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineValues(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchValues(t.db, []TargetValue{
			{"apple", "A", "y"},
			{"boat", "B", "z"},
			{"pear", "A", "x"},
			{"toy boat", "A", "w"},
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// match generated model defaults
func matchValues(db *sql.DB, want []TargetValue) (err error) {
	var curr TargetValue
	var have []TargetValue
	if e := tables.QueryAll(db,
		`select noun, field, value 
			from mdl_start
			order by noun, field, value`,
		func() (err error) {
			have = append(have, curr)
			return
		},
		&curr.Target, &curr.Field, &curr.Value); e != nil {
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
func addValues(rec *ephemera.Recorder, vals []TargetValue) (err error) {
	for _, v := range vals {
		noun := rec.Named(tables.NAMED_NOUN, v.Target, "test")
		prop := rec.Named(tables.NAMED_PROPERTY, v.Field, "test")
		value := v.Value
		rec.NewValue(noun, prop, value)
	}
	return
}
