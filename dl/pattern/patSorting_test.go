package pattern

import (
	"testing"

	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

func TestRuleSorting(t *testing.T) {
	ps := []*Rule{
		{Flags: Terminal, Execute: Text("1")},
		{Flags: Postfix, Execute: Text("2")},
		{Flags: Prefix, Execute: Text("3")},
		{Filter: Skip, Execute: Text("0")},
		{Flags: Postfix, Execute: Text("4")},
	}
	if inds, e := splitRules(nil, ps); e != nil {
		t.Fatal(e)
	} else if cnt := len(inds); cnt != 4 {
		t.Fatal("expected 4 matching rules")
	} else {
		var got string
		for _, i := range inds {
			got += string(ps[i].Execute.(Text))
		}
		if got != "3124" {
			t.Fatal("got", got)
		}
	}
}

type Text string

func (Text) Execute(rt.Runtime) error { return nil }

type Bool bool

func (b Bool) GetBool(rt.Runtime) (g.Value, error) {
	return g.BoolOf(bool(b)), nil
}

var Skip = Bool(false)
