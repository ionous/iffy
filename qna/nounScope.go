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

// when asking for a variable named x, check to see if x is a named noun.
// ex. maybe x is just a local loop counter, "i", or maybe x is "sam the dog".
func (ns *NounScope) GetVariable(name string) (interface{}, error) {
	return ns.fields.GetField(name, object.Id)
}
