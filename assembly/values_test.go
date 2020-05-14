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
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		if e := AddTestHierarchy(asm.modeler, []TargetField{
			{"K", ""},
			{"L", "K"},
			{"M", "L,K"},
		}); e != nil {
			t.Fatal(e)
		} else if e := AddTestFields(asm.modeler, []TargetValue{
			{"K", "t", tables.PRIM_TEXT},
			{"L", "d", tables.PRIM_DIGI},
		}); e != nil {
			t.Fatal(e)
		} else if e := AddTestNouns(asm.modeler, []TargetField{
			{"apple", "K"},
			{"pear", "L"},
			{"toy boat", "M"},
			{"boat", "M"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addValues(asm.rec, []TargetValue{
			{"apple", "t", "some text"},
			{"pear", "d", 123},
			{"toy", "d", 321},
			{"boat", "t", "more text"},
		}); e != nil {
			t.Fatal(e)
		} else if e := determineInitialFields(asm.modeler, asm.db); e != nil {
			t.Fatal(e)
		} else if e := matchValues(asm.db, []TargetValue{
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
	if asm, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer asm.db.Close()
		//
		if e := AddTestHierarchy(asm.modeler, []TargetField{
			{"K", ""},
			{"L", "K"},
			{"M", "L,K"},
		}); e != nil {
			t.Fatal(e)
		} else if e := AddTestFields(asm.modeler, []TargetValue{
			{"K", "A", tables.PRIM_ASPECT},
			{"L", "B", tables.PRIM_ASPECT},
		}); e != nil {
			t.Fatal(e)
		} else if e := AddTestTraits(asm.modeler, []TargetField{
			{"A", "w"}, {"A", "x"}, {"A", "y"},
			{"B", "z"},
		}); e != nil {
			t.Fatal(e)
		} else if e := AddTestNouns(asm.modeler, []TargetField{
			{"apple", "K"},
			{"pear", "L"},
			{"toy boat", "M"},
			{"boat", "M"},
		}); e != nil {
			t.Fatal(e)
		} else if e := addValues(asm.rec, []TargetValue{
			{"apple", "A", "y"},
			{"pear", "x", true},
			{"toy", "w", true},
			{"boat", "z", true},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineValues(asm.modeler, asm.db); e != nil {
			t.Fatal(e)
		} else if e := matchValues(asm.db, []TargetValue{
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
		noun := rec.NewName(v.Target, tables.NAMED_NOUN, "test")
		prop := rec.NewName(v.Field, tables.NAMED_PROPERTY, "test")
		value := v.Value
		rec.NewValue(noun, prop, value)
	}
	return
}
