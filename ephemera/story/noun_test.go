package story

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
	"github.com/kr/pretty"
)

func TestImportNamedNouns(t *testing.T) {
	k, db := newImporter(t, testdb.Memory)
	defer db.Close()
	//
	nouns := []string{
		"our", "Trevor",
		"an", "apple",
		"3", "triangles",
		"one", "square",
		"a gaggle of", "robot sheep",
	}
	for i := 0; i < len(nouns); i += 2 {
		n := makeNoun(nouns[i], nouns[i+1])
		if e := n.Import(k); e != nil {
			t.Fatal(e, "at", i)
		}
	}
	var buf strings.Builder
	tables.WriteCsv(db, &buf, "select count() from eph_noun", 1)
	tables.WriteCsv(db, &buf, `select en.name as noun, ep.name as trait
	from eph_value ev
	join eph_named en 
		on( ev.idNamedNoun == en.rowid )
	join eph_named ep 
		on( ev.idNamedProp == ep.rowid )
order by noun collate nocase, trait`, 2)
	tables.WriteCsv(db, &buf, "select name from eph_named where category='plural_kinds' order by name collate nocase", 1)
	tables.WriteCsv(db, &buf, "select name from eph_named where category='singular_kind' order by name collate nocase", 1)
	tables.WriteCsv(db, &buf, "select value from eph_value where value != 1 order by value", 1)
	tables.WriteCsv(db, &buf,
		`select nt.name, na.name
	from eph_trait t 
	join eph_named nt
		on (t.idNamedTrait = nt.rowid)
	left join eph_named na
		on (t.idNamedAspect = na.rowid)
	order by na.name, t.rank, nt.name`, 2)
	// tables.WriteCsv(db, &buf, "select name, category from eph_named where category != 'scene'", 2)
	have, want := buf.String(), lines(
		// note: we dont expect actual noun ephemera because we're only parsing names
		// we're not actually doing anything with those names.
		"0",
		//
		"apple,common_named",
		"apple,indefinite_article",
		//
		"robot_sheep,common_named",
		"robot_sheep,indefinite_article",
		"square_1,counted", // COUNTER:#
		"Trevor,indefinite_article",
		"Trevor,proper_named", // COUNTER:#
		"triangles_1,counted",
		"triangles_2,counted",
		"triangles_3,counted",
		//
		// "triangles", // plural -- disabled in ReadCountedNoun
		// "square",    // singular -- disabled in ReadCountedNoun
		// indefinite articles
		"a gaggle of",
		"an",
		"our",
		// implicitly generated aspects
		// listed in rank order (default first)
		"common_named,noun_types",
		"proper_named,noun_types",
		"counted,noun_types",
		"publicly_named,private_names",
		"privately_named,private_names",
		//
	)
	if diff := pretty.Diff(have, want); len(diff) > 0 {
		t.Fatal(have)
	}
}

func makeNoun(det, name string) NamedNoun {
	return NamedNoun{
		Determiner: Determiner{Str: det},
		Name:       NounName{Str: name},
	}
}
