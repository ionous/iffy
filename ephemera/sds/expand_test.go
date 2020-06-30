package sds

import (
	"encoding/json"
	"testing"
)

var compactData = `[
    "::test@id-3",{
      "test_name::text@id-0": "hello, goodbye",
      "go::execute": [
        "id-1::choose@id-7",{
          "false::execute": [
            "id-4::say@id-15",{
              "text::text_eval@id-14": [
              "::text_value@id-20", {
                "text::text@id-19": "goodbye"
              }]}],
          "if::bool_eval@id-5": [
            "::bool_value@id-9",{
              "::bool@id-8": "$TRUE"
            }],
          "true::execute": [
            "id-6::say@id-11",{
              "text::text_eval@id-10": [
                "::text_value@id-13",{
                  "text::text@id-12": "hello"
                }]}]}],
      "::text@id-2": "hello"
  }]`

func TestExpand(t *testing.T) {
	var in []interface{}
	if e := json.Unmarshal([]byte(compactData), &in); e != nil {
		t.Fatal(e)
	} else {
		for s := NewSlice(in); s != nil; s = s.Next() {
			t.Log(s.Note())
			el := s.Elem()
			ps := el.Params()
			t.Log(ps)
			// ::text@id-0 -- fix param names should be by short name

			n, p := el.Param("test_name")
			t.Log(n, p.Value())

		}
	}
}
