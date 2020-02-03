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

// check forming valid rel1, 1x, x1, xx
// check violations for rel1, 1x, x1, xx
// check the wrong nouns using a verb

func TestOneToOneFormation(t *testing.T) {
	if t, e := newRelativesTest(t, memory, [][3]string{
		{"a", "v1", "a"},
		{"b", "v1", "c"},
		{"c", "v1", "b"},
		{"z", "v1", "e"},
	}); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		if e := DetermineRelatives(t.db); e != nil {
			t.Fatal(e)
		} else if e := matchRelatives(t.db, [][3]string{
			{"a", "Rel1", "a"},
			{"b", "Rel1", "c"},
			{"e", "Rel1", "z"},
		}); e != nil {
			t.Fatal(e)
		} else {
			t.Log("okay")
		}
	}
}

func TestOneToManyFormation(t *testing.T) {
	if t, e := newRelativesTest(t, memory, [][3]string{
		{"z", "v1x", "f"},
		{"c", "v1x", "a"},
		{"b", "v1x", "e"},
		{"c", "v1x", "c"},
		{"c", "v1x", "b"},
		{"z", "v1x", "d"},
		{"z", "v1x", "f"},
	}); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		if e := DetermineRelatives(t.db); e != nil {
			t.Fatal(e)
		} else if e := matchRelatives(t.db, [][3]string{
			{"b", "Rel1x", "e"},
			{"c", "Rel1x", "a"},
			{"c", "Rel1x", "b"},
			{"c", "Rel1x", "c"},
			{"z", "Rel1x", "d"},
			{"z", "Rel1x", "f"},
		}); e != nil {
			t.Fatal(e)
		} else {
			t.Log("okay")
		}
	}
}

func TestManyToOneFormation(t *testing.T) {
	if t, e := newRelativesTest(t, memory, [][3]string{
		{"z", "vx1", "b"},
		{"f", "vx1", "f"},
		{"l", "vx1", "b"},
		{"b", "vx1", "a"},
		{"d", "vx1", "b"},
		{"c", "vx1", "d"},
		{"f", "vx1", "f"},
		{"e", "vx1", "f"},
	}); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		if e := DetermineRelatives(t.db); e != nil {
			t.Fatal(e)
		} else if e := matchRelatives(t.db, [][3]string{
			{"b", "Relx1", "a"},
			{"c", "Relx1", "d"},
			{"d", "Relx1", "b"},
			{"e", "Relx1", "f"},
			{"f", "Relx1", "f"},
			{"l", "Relx1", "b"},
			{"z", "Relx1", "b"},
		}); e != nil {
			t.Fatal(e)
		} else {
			t.Log("okay")
		}
	}
}

func TestManyToManyFormation(t *testing.T) {
	if t, e := newRelativesTest(t, memory, [][3]string{
		{"a", "vx", "a"},
		{"e", "vx", "d"},
		{"a", "vx", "b"},
		{"a", "vx", "c"},
		{"f", "vx", "d"},
		{"l", "vx", "d"},
		{"a", "vx", "b"},
	}); e != nil {
		t.Fatal(e)
	} else {
		defer t.Close()
		if e := DetermineRelatives(t.db); e != nil {
			t.Fatal(e)
		} else if e := matchRelatives(t.db, [][3]string{
			{"a", "Relx", "a"},
			{"a", "Relx", "b"},
			{"a", "Relx", "c"},
			{"e", "Relx", "d"},
			{"f", "Relx", "d"},
			{"l", "Relx", "d"},
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
	if e := dbutil.QueryAll(db,
		`select * 
			from start_rel
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

func xTestOneToOneViolations(t *testing.T) {
	test := func(add, want [][3]string) (err error) {
		if t, e := newRelativesTest(t, memory, add); e != nil {
			err = e
		} else {
			defer t.Close()
			var have [][3]string
			if e := oneToOneViolations(t.db, func(n, m, r string) (err error) {
				have = append(have, [3]string{n, m, r})
				return
			}); e != nil {
				err = e
			} else if !reflect.DeepEqual(have, want) {
				e := errutil.New("mismatch",
					"want", pretty.Sprint(want),
					"have", pretty.Sprint(have))
				err = e
			} else {
				t.Log("okay")
			}
		}
		return
	}
	if e := test([][3]string{
		{"a", "v1", "a"},
		{"a", "v1", "b"},
		{"a", "v1", "c"},
		{"d", "v1", "a"},
	}, [][3]string{
		{"a", "a", "Rel1"},
		{"a", "b", "Rel1"},
		{"a", "c", "Rel1"},
		{"d", "a", "Rel1"},
	}); e != nil {
		t.Fatal(e)
	}
	// if e := test([][3]string{
	// 	{"a", "v1", "a"},
	// 	{"b", "v1", "b"},
	// 	{"c", "v1", "c"},
	// }, nil); e != nil {
	// 	t.Fatal(e)
	// }
	// if e := test([][3]string{
	// 	{"a", "v1", "b"},
	// 	{"b", "v1", "a"},
	// 	{"c", "v1", "d"},
	// }, nil); e != nil {
	// 	t.Fatal(e)
	// }
}

func oneToOneViolations(db *sql.DB, cb func(n, m, r string) error) (err error) {
	var n, m, r string
	if e := dbutil.QueryAll(db,
		`select distinct asm1.noun, asm1.otherNoun, asm1.relation
	from asm_relation asm1
	join asm_relation asm2
		using (relation, cardinality)
	where (asm1.idEphRel != asm2.idEphRel)
		and (cardinality = 'one_one')
		and (asm1.noun = asm2.otherNoun
		or asm1.otherNoun = asm2.noun)
	order by relation, asm1.idEphRel`,
		func() (err error) {
			cb(n, m, r)
			return
		},
		&n, &m, &r); e != nil {
		err = e
	}
	return
}

func newRelativesTest(t *testing.T, path string, relatives [][3]string) (ret *assemblyTest, err error) {
	if t, e := newAssemblyTest(t, path); e != nil {
		err = e
	} else {
		if e := fakeHierarchy(t.modeler, []pair{
			{"K", ""},
			{"L", "K"},
			{"N", "K"},
		}); e != nil {
			err = e
		} else if e := fakeNouns(t.modeler, []pair{
			{"a", "K"},
			{"b", "K"},
			{"c", "K"},
			{"d", "K"},
			{"e", "K"},
			{"f", "K"},
			{"l", "L"},
			{"n", "N"},
		}); e != nil {
			err = e
		} else if e := fakeRelations(t.modeler, [][4]string{
			// relation, kind, cardinality, otherKind
			{"Rel1", "K", ephemera.ONE_TO_ONE, "K"},
			{"Rel1x", "K", ephemera.ONE_TO_MANY, "K"},
			{"Relx1", "K", ephemera.MANY_TO_ONE, "K"},
			{"Relx", "K", ephemera.MANY_TO_MANY, "K"},
		}); e != nil {
			err = e
		} else if e := fakeVerbs(t.modeler, [][2]string{
			// rel, verb
			{"Rel1", "v1"},
			{"Rel1x", "v1x"},
			{"Relx1", "vx1"},
			{"Relx", "vx"},
		}); e != nil {
			err = e
		} else if e := addRelatives(t.rec, relatives); e != nil {
			err = e
		}
		//
		if err != nil {
			t.Close()
		} else {
			ret = t
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
	name := rec.Named(ephemera.NAMED_NOUN, noun, "test")
	namedStem := rec.Named(ephemera.NAMED_NOUN, stem, "test")
	otherName := rec.Named(ephemera.NAMED_NOUN, otherNoun, "test")
	rec.NewRelative(name, namedStem, otherName)
	return
}

// add modeled data: relation, kind, cardinality, otherKind
func fakeRelations(m *Modeler, relations [][4]string) (err error) {
	for _, el := range relations {
		relation, kind, cardinality, otherKind := el[0], el[1], el[2], el[3]
		if e := m.WriteRelation(relation, kind, cardinality, otherKind); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func fakeVerbs(m *Modeler, relationStem [][2]string) (err error) {
	for _, el := range relationStem {
		rel, verb := el[0], el[1]
		if e := m.WriteVerb(rel, verb); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
