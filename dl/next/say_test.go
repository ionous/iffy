package next

import (
	"bytes"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/qna"
	"github.com/ionous/iffy/rt"
	"github.com/kr/pretty"
)

func TestPrint(t *testing.T) {
	// functions that turn execute blocks into text
	t.Run("span", func(t *testing.T) {
		var run sayTester
		span := &Span{rt.Block{
			&Say{&Text{"hello"}},
			&Say{&Text{"there"}},
			&Say{&Text{"world"}},
		}}
		var buf bytes.Buffer
		if e := span.WriteText(&run, &buf); e != nil {
			t.Fatal(e)
		} else if e := matchLine(buf.String(), "hello there world"); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("bracket", func(t *testing.T) {
		var run sayTester
		bracket := &Bracket{rt.Block{
			&Say{&Text{"hello"}},
			&Say{&Text{"there"}},
			&Say{&Text{"world"}},
		}}
		var buf bytes.Buffer
		if e := bracket.WriteText(&run, &buf); e != nil {
			t.Fatal(e)
		} else if e := matchLine(buf.String(), "( hello there world )"); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("slash", func(t *testing.T) {
		var run sayTester
		slash := &Slash{rt.Block{
			&Say{&Text{"hello"}},
			&Say{&Text{"there"}},
			&Say{&Text{"world"}},
		}}
		var buf bytes.Buffer
		if e := slash.WriteText(&run, &buf); e != nil {
			t.Fatal(e)
		} else if e := matchLine(buf.String(), "hello / there / world"); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("commas", func(t *testing.T) {
		var run sayTester
		commas := &Commas{rt.Block{
			&Say{&Text{"hello"}},
			&Say{&Text{"there"}},
			&Say{&Text{"world"}},
		}}
		var buf bytes.Buffer
		if e := commas.WriteText(&run, &buf); e != nil {
			t.Fatal(e)
		} else if e := matchLine(buf.String(), "hello, there, and world"); e != nil {
			t.Fatal(e)
		}
	})
}

func matchLine(have, want string) (err error) {
	if diff := pretty.Diff(have, want); len(diff) > 0 {
		err = errutil.Fmt("mismatch; have %q, want: %q, diff: %q", have, want, diff)
	}
	return
}

type sayBase struct {
	rt.Panic
}
type sayTester struct {
	sayBase
	qna.WriterStack
}
