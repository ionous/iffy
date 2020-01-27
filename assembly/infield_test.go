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

// TestFieldAssignment to verify initial values for fields can be assigned to instances.
func TestFieldAssignment(t *testing.T) {
	if t, e := newAssemblyTest(t, memory); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		m, rec := t.modeler, t.rec
		//
		if e := fakeHierarchy(m, []pair{
			{"T", ""},
			{"K", "T"},
			{"L", "K,T"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeNouns(m, []pair{
			{"apple", "T"},
			{"pear", "K"},
			{"machine gun", "L"},
		}); e != nil {
			t.Fatal(e)
		} else if e := fakeFields(m, []kfp{
			{"T", "t", ephemera.PRIM_TEXT},
			{"K", "d", ephemera.PRIM_DIGI},
		}); e != nil {
			t.Fatal(e)
		} else if e := addValues(rec, []prop{
			{"apple", "t", "some text"},
			{"pear", "d", 123},
			{"machine", "d", 321},
			{"gun", "t", "more text"},
		}); e != nil {
			t.Fatal(e)
		} else if e := DetermineValues(t.modeler, t.db); e != nil {
			t.Fatal(e)
		} else if e := matchValues(t.db, []prop{
			{"apple", "t", "some text"},
			{"machine gun", "d", int64(321)},
			{"machine gun", "t", "more text"},
			{"pear", "d", int64(123)}, // int64, re: go's default scanner.
		}); e != nil {
			t.Fatal(e)
		}
	}
}

// match generated model defaults
func matchValues(db *sql.DB, want []prop) (err error) {
	var curr prop
	var have []prop
	if e := dbutil.QueryAll(db,
		`select noun, field, value 
			from mdl_value
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
func addValues(rec *ephemera.Recorder, vals []prop) (err error) {
	for _, v := range vals {
		noun := rec.Named(ephemera.NAMED_NOUN, v.target, "test")
		prop := rec.Named(ephemera.NAMED_FIELD, v.prop, "test")
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
