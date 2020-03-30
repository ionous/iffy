package qna

type ScopeStack struct {
	stack []VariableScope
}

func (k *ScopeStack) PushScope(scope VariableScope) {
	k.stack = append(k.stack, scope)
}

func (k *ScopeStack) PopScope() {
	if cnt := len(k.stack); cnt == 0 {
		panic("ScopeStack: popping an empty stack")
	} else {
		k.stack = k.stack[0 : cnt-1]
	}
}

func (k *ScopeStack) Scope() VariableScope {
	return k
}

// GetVariable writes the value at 'name' into the value pointed to by 'pv'.
func (k *ScopeStack) GetVariable(name string, pv interface{}) (err error) {
	return k.visit(name, func(scope VariableScope) error {
		return scope.GetVariable(name, pv)
	})
}

// SetVariable writes the value of 'v' into the value at 'name'.
func (k *ScopeStack) SetVariable(name string, v interface{}) (err error) {
	return k.visit(name, func(scope VariableScope) error {
		return scope.SetVariable(name, v)
	})
}

func (k *ScopeStack) visit(name string, visitor func(VariableScope) error) (err error) {
	for i := len(k.stack) - 1; i >= 0; i-- {
		switch e := visitor(k.stack[i]); e.(type) {
		case nil:
			// no error? we're done.
			goto Done
		case UnknownVariable:
			// didn't find? keep looking...
		default:
			// other error? done.
			err = e
			goto Done
		}
	}
	err = UnknownVariable(name)
Done:
	return
}
