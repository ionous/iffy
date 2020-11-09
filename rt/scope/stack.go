package scope

import (
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type ScopeStack struct {
	stack []rt.Scope
}

func (k *ScopeStack) PushScope(scope rt.Scope) {
	if len(k.stack) > 25 {
		panic("stack overflow")
	}
	k.stack = append(k.stack, scope)
}

func (k *ScopeStack) PopScope() {
	if cnt := len(k.stack); cnt == 0 {
		panic("ScopeStack: popping an empty stack")
	} else {
		k.stack = k.stack[0 : cnt-1]
	}
}

// GetField returns the value at 'name'
func (k *ScopeStack) GetField(target, field string) (ret g.Value, err error) {
	norm := lang.Camelize(field)
	err = k.visit(target, field, func(scope rt.Scope) (err error) {
		if v, e := scope.GetField(target, norm); e != nil {
			err = e
		} else {
			ret = v
		}
		return err
	})
	return
}

// SetField writes the value of 'v' into the value at 'name'.
func (k *ScopeStack) SetField(target, field string, v g.Value) (err error) {
	norm := lang.Camelize(field)
	return k.visit(target, field, func(scope rt.Scope) error {
		return scope.SetField(target, norm, v)
	})
}
func (k *ScopeStack) visit(target, field string, visitor func(rt.Scope) error) (err error) {
	for i := len(k.stack) - 1; i >= 0; i-- {
		switch e := visitor(k.stack[i]); e.(type) {
		case nil:
			// no error? we're done.
			goto Done
		case rt.UnknownTarget, rt.UnknownField:
			// didn't find? keep looking...
		default:
			// other error? done.
			err = e
			goto Done
		}
	}
	err = rt.UnknownField{target, field}
Done:
	return
}
