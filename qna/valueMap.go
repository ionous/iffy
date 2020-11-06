package qna

import (
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/rt"
)

type keyType struct {
	target, field string
}

type valueMap map[keyType]rt.Value

func (k *keyType) unknown() error {
	return rt.UnknownField{k.target, k.field}
}

// subField should be one of the package object prefixes
func (k *keyType) dot(subField string) keyType {
	return keyType{subField, k.target + "." + k.field}
}

func makeKey(target, field string) keyType {
	// FIX?
	// operations generating get field should be registering the field as a name
	// and, as best as possible, relating obj to field for property verification
	// name translation should be done there.
	// we'd have to mark up things like text evaluations ( ex. HasTrait )
	if len(field) > 0 && field[0] != '#' {
		field = lang.Camelize(field)
	}
	return keyType{target, field}
}

func makeKeyForEval(obj, typeName string) keyType {
	return keyType{obj, typeName}
}
