package scope

import (
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
)

type ScopeStack struct {
	Scopes         []rt.Scope
	NormalizeNames bool
}

func (k *ScopeStack) PushScope(scope rt.Scope) {
	if len(k.Scopes) > 25 {
		panic("Scopes overflow")
	}
	k.Scopes = append(k.Scopes, scope)
}

func (k *ScopeStack) PopScope() {
	if cnt := len(k.Scopes); cnt == 0 {
		panic("ScopeStack: popping an empty Scopes")
	} else {
		k.Scopes = k.Scopes[0 : cnt-1]
	}
}

// GetField returns the value at 'name'
func (k *ScopeStack) GetField(target, field string) (ret g.Value, err error) {
	norm := field
	if k.NormalizeNames {
		norm = lang.Camelize(field)
	}
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
	norm := field
	if k.NormalizeNames {
		norm = lang.Camelize(field)
	}
	return k.visit(target, field, func(scope rt.Scope) error {
		return scope.SetField(target, norm, v)
	})
}
func (k *ScopeStack) visit(target, field string, visitor func(rt.Scope) error) (err error) {
	for i := len(k.Scopes) - 1; i >= 0; i-- {
		switch e := visitor(k.Scopes[i]); e.(type) {
		case nil:
			// no error? we're done.
			goto Done
		case g.UnknownTarget, g.UnknownField:
			// didn't find? keep looking...
		default:
			// other error? done.
			err = e
			goto Done
		}
	}
	err = g.UnknownField{target, field}
Done:
	return
}
