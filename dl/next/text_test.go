package next

import (
	"testing"

	"github.com/ionous/iffy/rt"
	"github.com/kr/pretty"
)

func TestText(t *testing.T) {
	var run baseRuntime

	testTrue := func(t *testing.T, eval rt.BoolEval) {
		if ok, e := rt.GetBool(&run, eval); e != nil {
			t.Fatal(e)
		} else if !ok {
			t.Fatal("expected true", pretty.Sprint(eval))
		}
	}
	t.Run("is", func(t *testing.T) {
		testTrue(t, &Is{&Bool{true}})
		testTrue(t, &IsNot{&Bool{false}})
	})

	t.Run("isEmpty", func(t *testing.T) {
		testTrue(t, &IsEmpty{&Text{}})
		testTrue(t, &IsNot{&IsEmpty{&Text{"xxx"}}})
	})

	t.Run("includes", func(t *testing.T) {
		testTrue(t, &Includes{
			&Text{"full"},
			&Text{"ll"},
		})
		testTrue(t, &IsNot{&Includes{
			&Text{"full"},
			&Text{"bull"},
		}})
	})
}
