package internal

import (
  "bytes"
  "database/sql"
  "encoding/gob"
  "testing"

  "github.com/ionous/iffy/dl/check"
  "github.com/ionous/iffy/dl/core"
  "github.com/ionous/iffy/export"
  "github.com/ionous/iffy/rt"
  "github.com/ionous/iffy/tables"
  "github.com/kr/pretty"
)

// read simple unit test story into memory as a golang struct
func TestImportProg(t *testing.T) {
  var prog check.TestOutput
  cmds := makeTypeMap(export.Slats)
  if e := ReadProg(&prog, sayStory, cmds); e != nil {
    t.Fatal(e)
  } else if diff := pretty.Diff(sayTest, prog); len(diff) > 0 {
    t.Fatal(diff)
  }
}

// read simple unit test story into sqlite, extract as a golang struct
func TestProcessProg(t *testing.T) {
  const memory = "file:test.db?cache=shared&mode=memory"
  if db, e := sql.Open("sqlite3", memory); e != nil {
    t.Fatal("db open", e)
  } else if e := tables.CreateEphemera(db); e != nil {
    t.Fatal("create eph", e)
  } else {
    defer db.Close()
    k := NewImporter(t.Name(), db, nil)
    if e := imp_test_output(k, k.eph.Named(t.Name(), "test", ""), sayStory); e != nil {
      t.Fatal(e)
    } else {
      var testName string
      var progid int
      var expect string
      if e := db.QueryRow("select * from eph_check").Scan(&testName, &progid, &expect); e != nil {
        t.Fatal(e)
      } else {
        t.Log(testName, progid, expect)
        //
        var id int
        var typeName string
        var prog []byte
        if e := db.QueryRow("select * from eph_prog").Scan(&id, &typeName, &prog); e != nil {
          t.Fatal(e)
        } else {
          var res check.TestOutput
          t.Log(id, typeName)
          dec := gob.NewDecoder(bytes.NewBuffer(prog))
          if e := dec.Decode(&res); e != nil {
            t.Fatal(e)
          } else if diff := pretty.Diff(sayTest, res); len(diff) > 0 {
            t.Fatal(diff)
          }
        }
      }
    }
  }
}

var sayTest = check.TestOutput{
  "hello", []rt.Execute{
    &core.Choose{
      If: &core.Bool{true},
      True: []rt.Execute{&core.Say{
        Text: &core.Text{"hello"},
      }},
      False: []rt.Execute{&core.Say{
        Text: &core.Text{"goodbye"},
      }},
    }},
}

var sayStory = map[string]interface{}{
  "type": "test_output",
  "value": map[string]interface{}{
    "$LINES": map[string]interface{}{
      "type":  "lines",
      "value": "hello",
    },
    "$GO": []interface{}{
      map[string]interface{}{
        "type": "execute",
        "value": map[string]interface{}{
          "type": "choose",
          "value": map[string]interface{}{
            "$FALSE": []interface{}{
              map[string]interface{}{
                "type": "execute",
                "value": map[string]interface{}{
                  "type": "say_text",
                  "value": map[string]interface{}{
                    "$TEXT": map[string]interface{}{
                      "type": "text_eval",
                      "value": map[string]interface{}{
                        "type": "text_value",
                        "value": map[string]interface{}{
                          "$TEXT": map[string]interface{}{
                            "type":  "lines",
                            "value": "goodbye",
                          }}}}}}}},
            "$IF": map[string]interface{}{
              "type": "bool_eval",
              "value": map[string]interface{}{
                "type": "bool_value",
                "value": map[string]interface{}{
                  "$BOOL": map[string]interface{}{
                    "type":  "bool",
                    "value": "$TRUE",
                  }}}},
            "$TRUE": []interface{}{
              map[string]interface{}{
                "type": "execute",
                "value": map[string]interface{}{
                  "type": "say_text",
                  "value": map[string]interface{}{
                    "$TEXT": map[string]interface{}{
                      "type": "text_eval",
                      "value": map[string]interface{}{
                        "type": "text_value",
                        "value": map[string]interface{}{
                          "$TEXT": map[string]interface{}{
                            "type":  "lines",
                            "value": "hello",
                          }}}}}}}}}}}}},
}
