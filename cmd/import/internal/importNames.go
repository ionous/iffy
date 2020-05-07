package internal

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
)

func importName(k *Importer, r reader.Map, readType, storeType string) (ret ephemera.Named, err error) {
	if len(storeType) == 0 {
		storeType = readType
	}
	if str, e := reader.String(r, readType); e != nil {
		err = e
	} else {
		ret = k.Named(storeType, str, reader.At(r))
	}
	return
}

func imp_variable_name(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return importName(k, r, "variable_name", "")
}

func imp_plural_kinds(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return importName(k, r, "plural_kinds", "")
}

func imp_singular_kind(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return importName(k, r, "singular_kind", "kind")
}

func imp_trait(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return importName(k, r, "trait", "")
}

func imp_relation(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return importName(k, r, "relation", "verb")
}

func imp_common_name(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return importName(k, r, "common_name", "noun")
}

func imp_proper_name(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return importName(k, r, "proper_name", "noun")
}

func imp_line_expr(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return importName(k, r, "lines", "expr")
}

func imp_property(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return importName(k, r, "field", "")
}

func imp_test_name(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return importName(k, r, "text", "test")
}
