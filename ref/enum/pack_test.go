package enum

import (
	"github.com/ionous/iffy/tests"
	r "reflect"
	"testing"
)

//
func TestPacking(t *testing.T) {
	t.Run("state to string", func(t *testing.T) {
		dst, src := "", tests.Maybe
		if ok := Unpack(r.ValueOf(&dst).Elem(), r.ValueOf(src)); !ok {
			t.Fatal("couldnt unpack")
		} else if dst != src.String() {
			t.Fatal("failed to copy", src, dst, len(dst))
		}
	})
	t.Run("string to state", func(t *testing.T) {
		dst, src := tests.No, "Maybe"
		if ok := Pack(r.ValueOf(&dst).Elem(), r.ValueOf(src)); !ok {
			t.Fatal("couldnt pack")
		} else if dst.String() != src {
			t.Fatal("failed to copy", src, dst)
		}
	})
}
