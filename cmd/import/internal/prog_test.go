package internal

import (
  "bytes"
  "database/sql"
  "encoding/gob"
  "encoding/json"
  "testing"

  "github.com/ionous/iffy/dl/check"
  "github.com/ionous/iffy/dl/core"
  "github.com/ionous/iffy/ephemera/reader"
  "github.com/ionous/iffy/export"
  "github.com/ionous/iffy/rt"
  "github.com/ionous/iffy/tables"
  "github.com/kr/pretty"
)

// read simple unit test story into memory as a golang struct
func TestImportProg(t *testing.T) {
  var in export.Dict
  if e := json.Unmarshal([]byte(sayStory), &in); e != nil {
    t.Fatal(e)
  } else {
    var prog check.Test
    if e := ImportProg(&prog, in, export.Runs); e != nil {
      t.Fatal(e)
    } else if diff := pretty.Diff(sayTest, prog); len(diff) > 0 {
      t.Fatal(diff)
    }
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
    var in reader.Map
    if e := json.Unmarshal([]byte(sayStory), &in); e != nil {
      t.Fatal("read json", e)
    } else {
      r := NewParser(t.Name(), db, fns)
      //
      if e := r.parseItem(in); e != nil {
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
            var res check.Test
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
}

var sayTest = check.Test{TestName: "hello, goodbye",
  Go: []rt.Execute{
    &core.Choose{
      If: &core.BoolValue{Bool: true},
      True: []rt.Execute{
        &core.Say{
          Text: &core.TextValue{Text: "hello"},
        },
      },
      False: []rt.Execute{
        &core.Say{
          Text: &core.TextValue{Text: "goodbye"},
        },
      },
    },
  },
  Lines: "hello",
}

var sayStory = `{
    "id": "id-1709ef632af-3",
    "type": "test",
    "value": {
      "$TEST_NAME": {
        "id": "id-1709ef632af-0",
        "type": "text",
        "value": "hello, goodbye"
      },
      "$GO": [
        {
          "id": "id-1709ef632af-1",
          "type": "execute",
          "value": {
            "id": "id-1709ef632af-7",
            "type": "choose",
            "value": {
              "$FALSE": [
                {
                  "id": "id-1709ef632af-4",
                  "type": "execute",
                  "value": {
                    "id": "id-1709ef632af-15",
                    "type": "say",
                    "value": {
                      "$TEXT": {
                        "id": "id-1709ef632af-14",
                        "type": "text_eval",
                        "value": {
                          "id": "id-1709ef632af-20",
                          "type": "text_value",
                          "value": {
                            "$TEXT": {
                              "id": "id-1709ef632af-19",
                              "type": "lines",
                              "value": "goodbye"
                            }
                          }
                        }
                      }
                    }
                  }
                }
              ],
              "$IF": {
                "id": "id-1709ef632af-5",
                "type": "bool_eval",
                "value": {
                  "id": "id-1709ef632af-9",
                  "type": "bool_value",
                  "value": {
                    "$BOOL": {
                      "id": "id-1709ef632af-8",
                      "type": "bool",
                      "value": "$TRUE"
                    }
                  }
                }
              },
              "$TRUE": [
                {
                  "id": "id-1709ef632af-6",
                  "type": "execute",
                  "value": {
                    "id": "id-1709ef632af-11",
                    "type": "say",
                    "value": {
                      "$TEXT": {
                        "id": "id-1709ef632af-10",
                        "type": "text_eval",
                        "value": {
                          "id": "id-1709ef632af-13",
                          "type": "text_value",
                          "value": {
                            "$TEXT": {
                              "id": "id-1709ef632af-12",
                              "type": "lines",
                              "value": "hello"
                            }
                          }
                        }
                      }
                    }
                  }
                }
              ]
            }
          }
        }
      ],
      "$LINES": {
        "id": "id-1709ef632af-2",
        "type": "lines",
        "value": "hello"
      }
    }
  }`
