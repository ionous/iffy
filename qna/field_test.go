package qna

import (
	"encoding/gob"
	"testing"

	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/pattern"
	"github.com/ionous/iffy/tables"
)

// test that pattern variables are accessible via objectValue's GetField
func TestPatternFields(t *testing.T) {
	gob.Register((*pattern.TextRule)(nil))
	gob.Register((*core.Text)(nil))
	gob.Register((*MatchNumber)(nil))

	db := newTestDB(t, memory)
	defer db.Close()
	if e := tables.CreateModel(db); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateRun(db); e != nil {
		t.Fatal(e)
	} else if e := tables.CreateRunViews(db); e != nil {
		t.Fatal(e)
	}
	m := assembly.NewModeler(db)
	if e := m.WritePat("pat", "param", "text_eval", 1); e != nil {
		t.Fatal(e)
	} else if e := m.WriteStart("pat", "param", "default"); e != nil {
		t.Fatal(e)
	}
	//
	run := NewRuntime(db)
	if p, e := run.GetField("pat", "param"); e != nil {
		t.Fatal(e)
	} else if v, e := core.GetText(run, p); e != nil {
		t.Fatal(e)
	} else if v != "default" {
		t.Fatal("mismatch", v)
	} else if p, e := run.GetFieldByIndex("pat", 1); e != nil {
		t.Fatal(e)
	} else if v, e := core.GetText(run, p); e != nil {
		t.Fatal(e)
	} else if v != "default" {
		t.Fatal("mismatch", v)
	} else if pairs := run.Fields.pairs; len(pairs) != 2 {
		t.Fatal("unexpected cached values", pairs)
	} else if val := pairs[keyType{"pat", "param"}]; val != "default" {
		t.Fatal("expected cached default; have", val)
	} else if field := pairs[keyType{"pat", "$1"}]; field != "param" {
		t.Fatal("expected cached param name; have", field)
	}
}
