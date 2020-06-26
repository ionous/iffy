package qna

import (
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt/scope"
)

// NounScope validates requests for noun names,
// returning the name or an error if the noun doesnt exist.
type NounScope struct {
	scope.EmptyScope
	fields *Fields
}

func (ns *NounScope) GetVariable(name string) (interface{}, error) {
	return ns.fields.GetField(name, object.Name)
}
