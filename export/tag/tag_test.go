package tag

import (
	r "reflect"
	"testing"
)

func TestMetaTags(t *testing.T) {
	t.Run("TestInvalidTag", func(t *testing.T) {
		tag := ReadTag(r.StructTag(`fi:"id"`))
		if cnt := len(tag); cnt != 0 {
			t.Fatal(tag)
		}
	})
	t.Run("TestKey", func(t *testing.T) {
		tag := ReadTag(r.StructTag(`if:"id"`))
		if id, ok := tag.Find("id"); !ok {
			t.Fatal("error")
		} else if cnt := len(id); cnt != 0 {
			t.Fatal(cnt)
		}
	})
	t.Run("TestSimilarTag", func(t *testing.T) {
		tag := ReadTag(r.StructTag(`if:"iddy=5"`))
		if _, ok := tag.Find("id"); ok {
			t.Fatal("error")
		}
	})
	t.Run("TestKeyValue", func(t *testing.T) {
		tag := ReadTag(r.StructTag(`if:"plural=tests"`))
		if _, ok := tag.Find("id"); ok {
			t.Fatal("error")
		} else if res, ok := tag.Find("plural"); !ok {
			t.Fatal("error")
		} else if str := res; str != "tests" {
			t.Fatal(str)
		}
	})
	t.Run("TestMultipleKeys", func(t *testing.T) {
		tag := ReadTag(r.StructTag(`if:"id,other=true,plural=tests"`))
		if id, ok := tag.Find("id"); !ok {
			t.Fatal("error")
		} else if cnt := len(id); cnt != 0 {
			t.Fatal(cnt)
		} else if other, ok := tag.Find("other"); !ok {
			t.Fatal("error")
		} else if str := other; str != "true" {
			t.Fatal(str)
		} else if res, ok := tag.Find("plural"); !ok {
			t.Fatal("error")
		} else if str := res; str != "tests" {
			t.Fatal(str)
		}
	})
	t.Run("TestMissingCommas", func(t *testing.T) {
		tag := ReadTag(r.StructTag(`if:"edge=test=tests"`))
		if res, ok := tag.Find("edge"); !ok {
			t.Fatal("error")
		} else if str := res; str != "test=tests" {
			t.Fatal(str)
		}
	})
	t.Run("TestEmptyOne", func(t *testing.T) {
		tag := ReadTag(r.StructTag(`if:"edge="`))
		if res, ok := tag.Find("edge"); !ok {
			t.Fatal("error")
		} else if len(res) > 0 {
			t.Fatal(res)
		}
	})
	t.Run("TestEmptyTwo", func(t *testing.T) {
		tag := ReadTag(r.StructTag(`if:"edge=,empty"`))
		if !tag.Exists("empty") {
			t.Fatal("error")
		} else if res, ok := tag.Find("edge"); !ok {
			t.Fatal("error")
		} else if len(res) > 0 {
			t.Fatal(res)
		}
	})
}
