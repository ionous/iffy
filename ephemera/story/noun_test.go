package story

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/tables"
	"github.com/kr/pretty"
)

func TestImportNamedNouns(t *testing.T) {
	k, db := newTestDecoder(t, memory)
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
		det, name := nouns[i], nouns[i+1]
		if e := imp_named_noun(k, makeNoun(det, name)); e != nil {
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
		"apple,commonNamed",
		"apple,indefiniteArticle",
		"robot sheep,commonNamed",
		"robot sheep,indefiniteArticle",
		"square#1,counted",
		"square#1,privatelyNamed",
		"Trevor,indefiniteArticle",
		"Trevor,properNamed",
		"triangles#1,counted",
		"triangles#1,privatelyNamed",
		"triangles#2,counted",
		"triangles#2,privatelyNamed",
		"triangles#3,counted",
		"triangles#3,privatelyNamed",
		//
		"triangles", // plural
		"square",    // singular
		// indefinite articles
		"a gaggle of",
		"an",
		"our",
		// implicitly generated aspects
		// listed in rank order (default first)
		"commonNamed,nounTypes",
		"properNamed,nounTypes",
		"counted,nounTypes",
		"publiclyNamed,privateNames",
		"privatelyNamed,privateNames",
		//
	)
	if diff := pretty.Diff(have, want); len(diff) > 0 {
		t.Fatal(have)
	}

}

func makeNoun(det, name string) map[string]interface{} {
	// test all three named noun types
	return map[string]interface{}{
		"type": "named_noun",
		"value": map[string]interface{}{
			"$DETERMINER": map[string]interface{}{
				"type":  "determiner",
				"value": det,
			},
			"$NAME": map[string]interface{}{
				"type":  "noun_name",
				"value": name,
			},
		},
	}
}
