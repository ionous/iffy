package scope

import (
	"testing"

	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

func TestScopeStack(t *testing.T) {
	names := []string{"inner", "outer", "top"}
	mocks := make(map[string]*mockScope)
	for _, n := range names {
		mocks[n] = &mockScope{name: n}
	}
	var stack ScopeStack

	// push and pop scopes onto the stack
	// we expect to hear these counts back
	counts := [][]int{
		{-1, -1, -1},
		{+0, -1, -1},
		{+1, +0, -1},
		{+2, +1, +0},
		{+3, +2, -1},
		{+4, -1, -1},
		{-1, -1, -1},
	}
	step := 0
	check := func(reason string) {
		count := counts[step]
		for i, name := range names {
			var have int
			switch p, e := stack.GetField(object.Variables, name); e.(type) {
			default:
				t.Fatal("fatal", e)
			case rt.UnknownTarget, rt.UnknownField:
				// t.Log(reason, "loop", i, "asking for", name, "... unknown")
				have = -1
			case nil:
				if v, e := p.GetNumber(nil); e != nil {
					t.Fatal("fatal", e)
				} else {
					have = int(v)
					// t.Log(reason, "loop", i, name, "got", have)
				}
			}

			//
			if want := count[i]; want != have {
				t.Fatal("fatal", reason, "step", step, name, "have:", have, "want:", want)
			} else {
				next := &generic.Int{Value: have + 1}
				switch e := stack.SetField(object.Variables, name, next); e.(type) {
				default:
					t.Fatal("fatal", reason, "step", step, name, "set failed", e)
				case rt.UnknownField:
					if have != -1 {
						t.Fatal("fatal", "step", step, name, "set failed", e)
					}
				case nil:
					if have == -1 {
						t.Fatal("fatal", reason, "step", step, name, "set unexpected success")
					} else {
						t.Log(reason, name, "set", next)
					}
				}
			}
		}
		step++

	}
	check("startup")
	for _, name := range names {
		m := mocks[name]
		stack.PushScope(m)
		check("pushed " + name)
	}
	for _, name := range names {
		stack.PopScope()
		check("popped " + name)
	}

	access := []int{5, 3, 1}
	for i, name := range names {
		m := mocks[name]
		cnt := access[i]
		if m.gets != cnt || m.sets != cnt {
			t.Fatal("fatal", name, "expected", cnt, "got", m.gets, m.sets)
		} else {
			t.Log(name, "accessed", cnt, "times")
		}
	}
}

type mockScope struct {
	name       string
	gets, sets int
	val        int
}

func (k *mockScope) GetField(target, field string) (ret rt.Value, err error) {
	if target != object.Variables {
		err = rt.UnknownTarget{target}
	} else if field != k.name {
		err = rt.UnknownField{target, field}
	} else {
		k.gets++
		ret = &generic.Int{Value: k.val}
	}
	return
}

func (k *mockScope) SetField(target, field string, v rt.Value) (err error) {
	if target != object.Variables {
		err = rt.UnknownTarget{target}
	} else if field != k.name {
		err = rt.UnknownField{target, field}
	} else if n, e := v.GetNumber(nil); e != nil {
		err = e
	} else {
		k.val = int(n)
		k.sets++
	}
	return
}
