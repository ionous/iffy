package coerce_test

import (
	r "reflect"
	"testing"

	"github.com/ionous/iffy/ref/coerce"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/sliceOf"
	"github.com/kr/pretty"
)

type TestText struct{ Text string }

func (t *TestText) GetText(rt.Runtime) (string, error) { return t.Text, nil }

//
func TestCoercion(t *testing.T) {
	// . bool <=> bool, string<=>string, eval<=>eval.
	// . NumericType <=> NumericType
	// . string <=> enumerated value
	t.Run("bool value", func(t *testing.T) {
		var dst, src bool
		src = true
		if e := coerce.Value(r.ValueOf(&dst).Elem(), r.ValueOf(src)); e != nil {
			t.Fatal(e)
		} else if dst != src {
			t.Fatal("failed to copy", src, dst)
		}
	})
	t.Run("string", func(t *testing.T) {
		var dst, src string
		src = "text"
		if e := coerce.Value(r.ValueOf(&dst).Elem(), r.ValueOf(src)); e != nil {
			t.Fatal(e)
		} else if dst != src {
			t.Fatal("failed to copy", src, dst)
		}
	})
	t.Run("eval", func(t *testing.T) {
		var dst rt.TextEval
		src := &TestText{"text"}
		if e := coerce.Value(r.ValueOf(&dst).Elem(), r.ValueOf(src)); e != nil {
			t.Fatal(e)
		} else if text, _ := dst.GetText(nil); text != src.Text {
			t.Fatal("failed to copy", src, dst)
		}
	})
	t.Run("int to float", func(t *testing.T) {
		var dst float64
		var src int = 23
		if e := coerce.Value(r.ValueOf(&dst).Elem(), r.ValueOf(src)); e != nil {
			t.Fatal(e)
		} else if dst != float64(src) {
			t.Fatal("failed to copy", src, dst)
		}
	})
	t.Run("float to int", func(t *testing.T) {
		var dst int
		var src float64 = 2.3
		if e := coerce.Value(r.ValueOf(&dst).Elem(), r.ValueOf(src)); e != nil {
			t.Fatal(e)
		} else if dst != int(src) {
			t.Fatal("failed to copy", src, dst)
		}
	})
	t.Run("strings", func(t *testing.T) {
		var dst, src []string
		src = sliceOf.String("ready", "set", "go")
		if e := coerce.Value(r.ValueOf(&dst).Elem(), r.ValueOf(src)); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(dst, src); len(diff) > 0 {
			t.Fatal("failed to copy", diff)
		} else {
			// FIX: what *is* right? maybe we want clone instead of assign?
			src[1] = "steady"
			if diff := pretty.Diff(dst, src); len(diff) > 0 {
				t.Fatal("should have assigned, not cloned", dst, src)
			}
		}
	})
	t.Run("ints to floats", func(t *testing.T) {
		dst := sliceOf.Float()
		src := sliceOf.Int(1, 1, 2, 3)
		if e := coerce.Value(r.ValueOf(&dst).Elem(), r.ValueOf(src)); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(dst, sliceOf.Float64(1, 1, 2, 3)); len(diff) > 0 {
			t.Fatal("failed to copy", diff)
		}
	})
	t.Run("floats to ints", func(t *testing.T) {
		dst := sliceOf.Int()
		src := sliceOf.Float(1.1, 2.3, 5.8)
		if e := coerce.Value(r.ValueOf(&dst).Elem(), r.ValueOf(src)); e != nil {
			t.Fatal(e)
		} else if diff := pretty.Diff(dst, sliceOf.Int(1, 2, 5)); len(diff) > 0 {
			t.Fatal("failed to copy", diff)
		}
	})

}
