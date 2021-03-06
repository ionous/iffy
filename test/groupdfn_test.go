package test

import (
	"testing"

	"github.com/ionous/iffy/test/testutil"
	"github.com/kr/pretty"
)

func TestKindsForType(t *testing.T) {
	var ks testutil.Kinds
	ks.AddKinds((*GroupCollation)(nil))
	if diff := pretty.Diff(ks.Fields, testutil.FieldMap{
		"Innumerable": {
			{"NotInnumerable", "bool", "trait"},
			{"IsInnumerable", "bool", "trait"},
		},
		"GroupOptions": {
			{"WithoutObjects", "bool", "trait"},
			{"ObjectsWithArticles", "bool", "trait"},
			{"ObjectsWithoutArticles", "bool", "trait"},
		},
		"GroupSettings": {
			{"Name", "text", "string"},
			{"Label", "text", "string"},
			{"Innumerable", "text", "aspect"},
			{"GroupOptions", "text", "aspect"},
		},
		"GroupedObjects": {
			{"Settings", "record", "GroupSettings"},
			{"Objects", "text_list", "string"},
		},
		"GroupCollation": {
			{"Groups", "record_list", "GroupedObjects"},
		},
	}); len(diff) > 0 {
		t.Fatal(pretty.Println(ks.Fields))
	}
}
