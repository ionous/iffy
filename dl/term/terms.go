package term

import (
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/rt"
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/safe"
)

// Terms implements a Scope mapping names to specified parameters.
// The only current user is pattern.FromPattern::Stitch()
// It stores values from indexed and key name arguments ( originally specified as evals. )
// Its pushed into scope so the names can be used as a source of values for rt.Runtime::GetField().
// ( ex. For use with the commands Var{},  SimpleNoun{}, ProperNoun{}, ObjectName{}, ... )
type Terms []g.Field

// rather than copying etc here, can we use records --
// what is it anyway?
func (ps Terms) NewKind(kinds g.Kinds) *g.Kind {
	return g.NewKind(kinds, "", append([]g.Field{}, ps...))
}

// fix:in the declaration of a pattern, we allow the specification of objects and types
// and so probably? the specification should become a text filter for the pattern
func (ps Terms) ConvertTerm(run rt.Runtime, i int, v g.Value) (ret g.Value, err error) {
	ret = v          // provisionally
	if len(ps) > 0 { // the list can be nil meaning no conversion...
		if field := ps[i]; field.Type == affine.Object.String() { // our hack...
			switch aff := v.Affinity(); aff {
			// fix? templates do send us some object values right now...
			case affine.Object:
				ret = g.ObjectAsText(v)
			// its text into text, where the input is supposed to be an id...
			// lets make sure it *is* an id.
			case affine.Text:
				if v, e := safe.ObjectFromString(run, v.String()); e != nil {
					err = e
				} else {
					ret = g.ObjectAsText(v)
				}
			}
		}
	}
	return
}

func (ps *Terms) AddTerm(name string, affinity affine.Affinity, typeName string) {
	// see notes: ConvertTerm
	if affinity == affine.Object {
		affinity = affine.Text
		typeName = affine.Object.String()
	}
	(*ps) = append((*ps), g.Field{Name: name, Affinity: affinity, Type: typeName})
}
