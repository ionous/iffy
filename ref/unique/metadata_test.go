package unique

import (
	r "reflect"
	"testing"
)

func makeMetadata(tag string) (out Metadata) {
	s := Tag(r.StructTag(tag))
	if len(s) > 0 {
		out = make(Metadata)
		out.AddString(s, "fill")
	}
	return
}

func TestMetadata(t *testing.T) {
	t.Run("TestKey", func(t *testing.T) {
		m := makeMetadata(`if:"id"`)
		if cnt := len(m); cnt != 1 {
			t.Fatal(cnt)
		} else if str := m["id"]; str != "fill" {
			t.Fatal(str)
		}
	})
	t.Run("TestKeyValue", func(t *testing.T) {
		m := makeMetadata(`if:"plural:tests"`)
		if cnt := len(m); cnt != 1 {
			t.Fatal(cnt)
		} else if str := m["plural"]; str != "tests" {
			t.Fatal(str)
		}
	})
	t.Run("TestMultipleKeys", func(t *testing.T) {
		m := makeMetadata(`if:"id,plural:tests"`)
		if cnt := len(m); cnt != 2 {
			t.Fatal(cnt)
		} else if str := m["id"]; str != "fill" {
			t.Fatal(str)
		} else if str := m["plural"]; str != "tests" {
			t.Fatal(str)
		}
	})
	t.Run("TestMissingCommas", func(t *testing.T) {
		m := makeMetadata(`if:"edge:test:tests"`)
		if cnt := len(m); cnt != 1 {
			t.Fatal(cnt)
		} else if str := m["edge"]; str != "test:tests" {
			t.Fatal(str)
		}
	})
}

func TestMetaTags(t *testing.T) {

	t.Run("TestInvalidTag", func(t *testing.T) {
		tag := Tag(r.StructTag(`fi:"id"`))
		if cnt := len(tag); cnt != 0 {
			t.Fatal(tag)
		}
	})
	t.Run("TestKey", func(t *testing.T) {
		tag := Tag(r.StructTag(`if:"id"`))
		if id, ok := tag.Find("id"); !ok {
			t.Fatal("error")
		} else if cnt := len(id); cnt != 0 {
			t.Fatal(cnt)
		}
	})
	t.Run("TestSimilarTag", func(t *testing.T) {
		tag := Tag(r.StructTag(`if:"iddy:5"`))
		if _, ok := tag.Find("id"); ok {
			t.Fatal("error")
		}
	})
	t.Run("TestKeyValue", func(t *testing.T) {
		tag := Tag(r.StructTag(`if:"plural:tests"`))
		if _, ok := tag.Find("id"); ok {
			t.Fatal("error")
		} else if res, ok := tag.Find("plural"); !ok {
			t.Fatal("error")
		} else if str := res; str != "tests" {
			t.Fatal(str)
		}
	})
	t.Run("TestMultipleKeys", func(t *testing.T) {
		tag := Tag(r.StructTag(`if:"id,other:true,plural:tests"`))
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
		tag := Tag(r.StructTag(`if:"edge:test:tests"`))
		if res, ok := tag.Find("edge"); !ok {
			t.Fatal("error")
		} else if str := res; str != "test:tests" {
			t.Fatal(str)
		}
	})
}
