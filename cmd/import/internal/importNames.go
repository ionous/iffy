package internal

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
)

// make.str("variable_name");
func imp_variable_name(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return k.namedType(r, "variable_name")
}

func imp_plural_kinds(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return k.namedType(r, "plural_kinds")
}

// make.str("trait");
func imp_trait(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return k.namedType(r, "trait")
}
