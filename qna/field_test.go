package qna

import (
	"testing"
)

// test that pattern variables are accessible via objectValue's GetField
// this is no longer a true thing --
// we should only be able to get pattern variables via GetVariable...
// revisit to see what other pattern access should be tested?
func xTestPatternFields(t *testing.T) {
	// gob.Register((*core.Text)(nil))
	// gob.Register((*debug.MatchNumber)(nil))

	// db := newQnaDB(t, testdb.Memory)
	// defer db.Close()
	// if e := tables.CreateModel(db); e != nil {
	// 	t.Fatal(e)
	// } else if e := tables.CreateRun(db); e != nil {
	// 	t.Fatal(e)
	// } else if e := tables.CreateRunViews(db); e != nil {
	// 	t.Fatal(e)
	// }
	// m := assembly.NewAssembler(db)
	// if e := m.WritePat("pat", "param", "text_eval", 1); e != nil {
	// 	t.Fatal(e)
	// } else if e := m.WriteStart("pat", "param", "default"); e != nil {
	// 	t.Fatal(e)
	// }
	// //
	// run := NewRuntime(db)
	// if p, e := run.GetField("pat", "param"); e != nil {
	// 	t.Fatal(e)
	// } else if v, e := p.GetText(run); e != nil {
	// 	t.Fatal(e)
	// } else if v != "default" {
	// 	t.Fatal("mismatch", v)
	// } else /*if field, e := run.GetFieldByIndex("pat", 1); e != nil {
	// 	t.Fatal(e)
	// } else if field != "param" {
	// 	t.Fatal(e)
	// } else */if pairs := run.pairs; len(pairs) != 3 {
	// 	t.Fatal("unexpected cached values", pairs, len(pairs))
	// } else if val, ok := pairs[keyType{"pat", "param"}]; !ok {
	// 	t.Fatal("missing cached default; have", val)
	// } else if str, e := val.GetText(run); e != nil {
	// 	t.Fatal(e)
	// } else if str != "default" {
	// 	t.Fatal("expected cached default; have", val)
	// } /*else if field := pairs[keyType{"pat", "$1"}]; field != "param" {
	// 	t.Fatal("expected cached param name; have", field)
	// }*/
	// FIX -- test pattern prep ( what was GetFieldByIndex )
}
