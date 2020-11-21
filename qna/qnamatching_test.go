package qna

import (
	"encoding/gob"
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/tables"
)

// manually add an assembled pattern to the database, test that it works as expected.
func TestLike(t *testing.T) {
	gob.Register((*core.MatchLike)(nil))

	db := newQnaDB(t, memory)
	defer db.Close()
	if e := tables.CreateAll(db); e != nil {
		t.Fatal(e)
	}
	run := NewRuntime(db)
	tests := []struct {
		text, pattern string
		ok            bool
	}{
		{
			"neon light", "neon%", true,
		},
		{
			"neon light", "%neon", false,
		},
		{
			"neon light", "neon_light", true,
		},
		{
			"neonlight", "neon_light", false,
		},
	}

	for i, test := range tests {
		c := core.MatchLike{
			Text:    &core.Text{test.text},
			Pattern: &core.Text{test.pattern},
		}
		if ok, e := c.GetBool(run); e != nil {
			t.Fatal(e)
		} else if ok := ok.Bool(); ok != test.ok {
			t.Fatalf("test %v (%q like %q) != %v", i, test.text, test.pattern, test.ok)
		} else {
			t.Log(test.text, "like", test.pattern, ok)
		}
	}
}
