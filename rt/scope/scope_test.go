package scope

import (
	"testing"

	"github.com/ionous/iffy/assign"
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
			switch p, e := stack.GetVariable(name); e.(type) {
			default:
				t.Fatal("fatal", e)
			case UnknownVariable:
				// t.Log(reason, "loop", i, "asking for", name, "... unknown")
				have = -1
			case nil:
				if v, e := assign.ToInt(p); e != nil {
					t.Fatal("fatal", e)
				} else {
					have = v
					// t.Log(reason, "loop", i, name, "got", have)
				}
			}

			//
			if want := count[i]; want != have {
				t.Fatal("fatal", reason, "step", step, name, "have:", have, "want:", want)
			} else {
				next := have + 1
				switch e := stack.SetVariable(name, next); e.(type) {
				default:
					t.Fatal("fatal", reason, "step", step, name, "set failed", e)
				case UnknownVariable:
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

func (k *mockScope) GetVariable(name string) (ret interface{}, err error) {
	if name != k.name {
		err = UnknownVariable(name)
	} else {
		k.gets++
		ret = k.val
	}
	return
}

func (k *mockScope) SetVariable(name string, v interface{}) (err error) {
	if name != k.name {
		err = UnknownVariable(name)
	} else {
		k.sets++
		k.val = v.(int)
	}
	return
}
