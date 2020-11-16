package term

import (
	"github.com/ionous/iffy/affine"
	g "github.com/ionous/iffy/rt/generic"
)

// Terms implements a Scope mapping names to specified parameters.
// The only current user is pattern.FromPattern::Stitch()
// It stores values from indexed and key name arguments ( originally specified as evals. )
// Its pushed into scope so the names can be used as a source of values for rt.Runtime::GetField().
// ( ex. For use with the commands GetVar{},  SimpleNoun{}, ProperNoun{}, ObjectName{}, ... )
type Terms struct {
	fields []g.Field
}

// rather than copying etc here, can we use records --
// what is it anyway?
func (ps *Terms) NewKind(kinds g.Kinds) *g.Kind {
	return g.NewKind(kinds, "", ps.fields)
}

func (ps *Terms) AddTerm(name string, affinity affine.Affinity, typeName string) {
	ps.fields = append(ps.fields, g.Field{Name: name, Affinity: affinity, Type: typeName})
}
