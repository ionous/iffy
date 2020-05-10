package decode

import (
  "testing"

  "github.com/ionous/iffy/ephemera/debug"
  "github.com/ionous/iffy/export"
  "github.com/kr/pretty"
)

// read simple unit test story into memory as a golang struct
func TestDecode(t *testing.T) {
  dec := NewDecoder()
  // register creation functions for all the slats.
  dec.AddDefaultCallbacks(export.Slats)
  // read say story data
  if prog, e := dec.ReadProg(debug.SayStory); e != nil {
    t.Fatal(e)
  } else if diff := pretty.Diff(&debug.SayTest, prog); len(diff) > 0 {
    t.Fatal(diff)
  } else {
    t.Log(pretty.Sprint(prog))
  }
}
