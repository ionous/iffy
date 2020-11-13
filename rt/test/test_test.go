package test

import (
	"testing"

	"github.com/kr/pretty"
)

func TestKindsForType(t *testing.T) {
	var ks Kinds
	ks.Add((*Groupings)(nil))
	if diff := pretty.Diff(ks.fields, fieldMap{
		"GroupedObjects": {
			{"Settings", "record", "GroupSettings"},
			{"Objects", "text_list", "string"},
		},
		"Groupings": {
			{"Ungrouped", "record_list", "GroupSettings"},
			{"Groups", "record_list", "GroupedObjects"},
		},
		"ObjectGrouping": {
			{"WithoutObjects", "bool", "trait"},
			{"WithoutArticles", "bool", "trait"},
			{"WithArticles", "bool", "trait"},
		},
		"GroupSettings": {
			{"Name", "text", "string"},
			{"Label", "text", "string"},
			{"Innumerable", "bool", "bool"},
			{"ObjectGrouping", "text", "aspect"},
		},
	}); len(diff) > 0 {
		t.Fatal(diff)
	}
}
