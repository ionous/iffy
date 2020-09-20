package qna

import (
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rt/generic"
)

type keyType struct {
	target, field string
}

func (k *keyType) unknown() error {
	return rt.UnknownField{k.target, k.field}
}

func (k *keyType) dot(subField string) keyType {
	return keyType{k.target + "." + k.field, subField}
}

type valueMap map[keyType]*generic.Value

func makeKey(obj, field string) keyType {
	// FIX FIX FIX --
	// operations generating get field should be registering the field as a name
	// and, as best as possible, relating obj to field for property verification
	// name translation should be done there.
	if len(field) > 0 && field[0] != object.Prefix {
		field = lang.Camelize(field)
	}
	return keyType{obj, field}
}

func makeKeyForEval(obj, typeName string) keyType {
	return keyType{obj, typeName}
}
