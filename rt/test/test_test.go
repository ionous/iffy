package test

import (
	"testing"

	"github.com/kr/pretty"
)

func TestKindsForType(t *testing.T) {
	var ks Kinds
	ks.Add((*GroupCollation)(nil))
	if diff := pretty.Diff(ks.fields, fieldMap{
		"Innumerable": {
			{"Is Innumerable", "bool", "trait"},
			{"Not Innumerable", "bool", "trait"},
		},
		"GroupOptions": {
			{"WithoutObjects", "bool", "trait"},
			{"WithoutArticles", "bool", "trait"},
			{"WithArticles", "bool", "trait"},
		},
		"GroupSettings": {
			{"Name", "text", "string"},
			{"Label", "text", "string"},
			{"Innumerable", "text", "aspect"},
			{"Options", "text", "aspect"},
		},
		"GroupObjects": {
			{"Settings", "record", "GroupSettings"},
			{"Objects", "text_list", "string"},
		},
		"GroupCollation": {
			{"Groups", "record_list", "GroupObjects"},
		},
	}); len(diff) > 0 {
		t.Fatal(pretty.Println(ks.fields))
	}
}
