package qna

import (
	"encoding/gob"
	"testing"

	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/tables"
)

// manually add an assembled pattern to the database, test that it works as expected.
func TestSayMe(t *testing.T) {
	gob.Register((*core.Text)(nil))
	gob.Register((*debug.MatchNumber)(nil))

	db := newQnaDB(t, memory)
	defer db.Close()
	if e := tables.CreateModel(db); e != nil {
		t.Fatal(e)
	}
	m := assembly.NewAssembler(db)
	if _, e := m.WriteGob("sayMe", &debug.SayPattern); e != nil {
		t.Fatal(e)
	}
	//
	if e := tables.CreateRun(db); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateRunViews(db); e != nil {
		t.Fatal(e)
	}
	run := NewRuntime(db)
	for i, expect := range []string{"One!", "Two!", "Three!", "Not between 1 and 3."} {
		if text, e := debug.DetermineSay(i + 1).GetText(run); e != nil {
			t.Fatal(e)
		} else if expect != text {
			t.Fatal(i, text)
		} else {
			t.Log(text)
		}
	}
}
