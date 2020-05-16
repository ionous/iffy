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

// todo: check the wrong nouns using a verb, using the wrong verb, etc.
// todo: ensure that the same stem can be used in multiple relations ( so long as the kinds differ, ex. in room, vs in box. )
func TestRelativeFormation(t *testing.T) {
	if asm, e := newRelativesTest(t, memory, [][3]string{
		{"a", "v1", "a"},
		{"b", "v1", "c"},
		{"c", "v1", "b"},
		{"z", "v1", "e"},
		{"b", "v1", "d"},
		{"c", "v1", "a"},
		{"z", "v1", "f"},
		//
		{"z", "v1x", "f"},
		{"c", "v1x", "a"},
		{"b", "v1x", "e"},
		{"c", "v1x", "c"},
		{"c", "v1x", "b"},
		{"z", "v1x", "d"},
		{"z", "v1x", "f"},
		//
		{"z", "vx1", "b"},
		{"f", "vx1", "f"},
		{"l", "vx1", "b"},
		{"b", "vx1", "a"},
		{"d", "vx1", "b"},
		{"c", "vx1", "d"},
		{"f", "vx1", "f"},
		{"e", "vx1", "f"},
		//
		{"a", "vx", "a"},
		{"e", "vx", "d"},
		{"a", "vx", "b"},
		{"a", "vx", "c"},
		{"f", "vx", "d"},
		{"l", "vx", "d"},
		{"a", "vx", "b"},
	}); e != nil {
		t.Fatal(e)
		return
	} else {
		defer asm.db.Close()
		if e := AssembleRelatives(asm.modeler, asm.db); e != nil {
			t.Fatal(e)
		} else if e := matchRelatives(asm.db, [][3]string{
			{"Rel1", "a", "a"},
			{"Rel1", "b", "c"},
			{"Rel1", "e", "z"},
			//
			{"Rel1x", "b", "e"},
			{"Rel1x", "c", "a"},
			{"Rel1x", "c", "b"},
			{"Rel1x", "c", "c"},
			{"Rel1x", "z", "d"},
			{"Rel1x", "z", "f"},
			//
			{"Relx1", "b", "a"},
			{"Relx1", "c", "d"},
			{"Relx1", "d", "b"},
			{"Relx1", "e", "f"},
			{"Relx1", "f", "f"},
			{"Relx1", "l", "b"},
			{"Relx1", "z", "b"},
			//
			{"Relxx", "a", "a"},
			{"Relxx", "a", "b"},
			{"Relxx", "a", "c"},
			{"Relxx", "e", "d"},
			{"Relxx", "f", "d"},
			{"Relxx", "l", "d"},
		}); e != nil {
			t.Fatal(e)
		} else {
			t.Log("okay")
		}
	}
}

func matchRelatives(db *sql.DB, want [][3]string) (err error) {
	var curr [3]string
	var have [][3]string
	if e := tables.QueryAll(db,
		`select relation, noun, otherNoun
			from mdl_pair
			order by relation, noun, otherNoun`,
		func() (err error) {
			have = append(have, curr)
			return
		},
		&curr[0], &curr[1], &curr[2]); e != nil {
		err = e
	} else if !reflect.DeepEqual(have, want) {
		err = errutil.New("mismatch",
			"have:", pretty.Sprint(have),
			"want:", pretty.Sprint(want))
	}
	return
}

func TestOneToOneViolations(t *testing.T) {
	test := func(add, want [][3]string) (err error) {
		if asm, e := newRelativesTest(t, memory, add); e != nil {
			err = e
		} else {
			defer asm.db.Close()
			if e := AssembleRelatives(asm.modeler, asm.db); e != nil {
				err = e
			} else {
				var have [][3]string
				var got [3]string
				if e := tables.QueryAll(asm.db,
					`select distinct coalesce(noun, ''), 
									 coalesce(stem, ''), 
									 coalesce(otherNoun, '')
					from asm_mismatch`,
					func() (err error) {
						have = append(have, got)
						return
					},
					&got[0], &got[1], &got[2]); e != nil {
					err = e
				} else if !reflect.DeepEqual(have, want) {
					e := errutil.New("mismatch",
						"want", pretty.Sprint(want),
						"have", pretty.Sprint(have))
					err = e
				}
			}
		}
		return
	}
	if e := test([][3]string{
		{"a", "v1", "a"},
		{"a", "v1", "b"},
		{"a", "v1", "c"},
		{"d", "v1", "a"},
		//
	}, [][3]string{
		{"a", "v1", "b"},
		{"a", "v1", "c"},
		{"a", "v1", "d"}, // nouns are sorted
		//
	}); e != nil {
		t.Fatal(e)
	}
	if e := test([][3]string{
		{"z", "v1x", "f"},
		{"c", "v1x", "f"},
		{"z", "v1x", "a"},
	}, [][3]string{
		{"c", "v1x", "f"},
	}); e != nil {
		t.Fatal(e)
	}
	if e := test([][3]string{
		{"f", "v1x", "c"},
		{"b", "v1x", "c"},
		{"b", "v1x", "d"},
	}, [][3]string{
		{"b", "v1x", "c"},
	}); e != nil {
		t.Fatal(e)
	}

}

func newRelativesTest(t *testing.T, path string, relatives [][3]string) (ret *assemblyTest, err error) {
	ret = &assemblyTest{T: t}
	if asm, e := newAssemblyTest(t, path); e != nil {
		err = e
	} else {
		if e := AddTestHierarchy(asm.modeler, []TargetField{
			{"K", ""},
			{"L", "K"},
			{"N", "K"},
		}); e != nil {
			err = e
		} else if e := AddTestNouns(asm.modeler, []TargetField{
			{"a", "K"},
			{"b", "K"},
			{"c", "K"},
			{"d", "K"},
			{"e", "K"},
			{"f", "K"},
			{"l", "L"},
			{"n", "N"},
			{"z", "K"},
		}); e != nil {
			err = e
		} else if e := AddTestRelations(asm.modeler, [][4]string{
			// relation, kind, cardinality, otherKind
			{"Rel1", "K", tables.ONE_TO_ONE, "K"},
			{"Rel1x", "K", tables.ONE_TO_MANY, "K"},
			{"Relx1", "K", tables.MANY_TO_ONE, "K"},
			{"Relxx", "K", tables.MANY_TO_MANY, "K"},
		}); e != nil {
			err = e
		} else if e := AddTestVerbs(asm.modeler, [][2]string{
			// rel, verb
			{"Rel1", "v1"},
			{"Rel1x", "v1x"},
			{"Relx1", "vx1"},
			{"Relxx", "vx"},
		}); e != nil {
			err = e
		} else if e := addRelatives(asm.rec, relatives); e != nil {
			err = e
		}
		//
		if err != nil {
			asm.db.Close()
		} else {
			ret = asm
		}
	}
	return
}

// add noun, stem, otherNoun ephemera
func addRelatives(rec *ephemera.Recorder, els [][3]string) (err error) {
	for _, el := range els {
		noun, stem, otherNoun := el[0], el[1], el[2]
		if e := addRelative(rec, noun, stem, otherNoun); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// add ephemera
func addRelative(rec *ephemera.Recorder, noun, stem, otherNoun string) (err error) {
	name := rec.NewName(noun, tables.NAMED_NOUN, "test")
	namedStem := rec.NewName(stem, tables.NAMED_VERB, "test")
	otherName := rec.NewName(otherNoun, tables.NAMED_NOUN, "test")
	rec.NewRelative(name, namedStem, otherName)
	return
}
