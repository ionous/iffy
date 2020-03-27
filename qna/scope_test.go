package qna

import "testing"

func TestScopeStack(t *testing.T) {
	names := []string{"inner", "outer", "top"}
	mocks := make(map[string]*mockScope)
	for _, n := range names {
		mocks[n] = &mockScope{name: n}
	}
	var stack ScopeStack

	// push and pop scopes onto the stack
	// we expect to here these counts back
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
	check := func() {
		count := counts[step]
		for i, name := range names {
			var have int
			switch e := stack.GetVariable(name, &have); e.(type) {
			default:
				t.Fatal(e)
			case UnknownVariable:
				have = -1
			case nil:
				t.Log(name, "got", have)
				break
			}
			if want := count[i]; want != have {
				t.Fatal("step", step, name, "have:", have, "want:", want)
			} else {
				next := have + 1
				switch e := stack.SetVariable(name, next); e.(type) {
				default:
					t.Fatal("step", step, name, "set failed", e)
				case UnknownVariable:
					if have != -1 {
						t.Fatal("step", step, name, "set failed", e)
					}
				case nil:
					if have == -1 {
						t.Fatal("step", step, name, "set unexpected success")
					} else {
						t.Log(name, "set", next)
					}
				}
			}
		}
		step++

	}
	check()
	for _, name := range names {
		m := mocks[name]
		stack.PushScope(m)
		check()
	}
	for range names {
		stack.PopScope()
		check()
	}

	access := []int{5, 3, 1}
	for i, name := range names {
		m := mocks[name]
		cnt := access[i]
		if m.gets != cnt || m.sets != cnt {
			t.Fatal(name, "expected", cnt, "got", m.gets, m.sets)
		} else {
			t.Log(name, "accessed", cnt, "times")
		}
	}
}

type mockScope struct {
	name       string
	gets, sets int
	val        int64
}

func (k *mockScope) GetVariable(name string, pv interface{}) (err error) {
	if name != k.name {
		err = UnknownVariable(name)
	} else {
		k.gets++
		err = Assign(pv, k.val)
	}
	return
}

func (k *mockScope) SetVariable(name string, v interface{}) (err error) {
	if name != k.name {
		err = UnknownVariable(name)
	} else {
		k.sets++
		err = Assign(&k.val, v)
	}
	return
}
