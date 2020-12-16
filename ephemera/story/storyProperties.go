package story

import (
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

func (op *PropertyDecl) ImportProperty(k *Importer, kind ephemera.Named) (err error) {
	if prop, e := op.Property.NewName(k); e != nil {
		err = e
	} else {
		err = op.PropertyType.ImportPropertyType(k, kind, prop)
		// Comment      *Lines
	}
	return
}

func (op *PropertyType) ImportPropertyType(k *Importer, kind, prop ephemera.Named) (err error) {
	type propertyTypeImporter interface {
		ImportPropertyType(k *Importer, kind, prop ephemera.Named) error
	}
	if opt, ok := op.Opt.(propertyTypeImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	} else {
		err = opt.ImportPropertyType(k, kind, prop)
		// Comment      *Lines
	}
	return
}

// inform gives these the name "<noun> condition"
// we could only do that with an after the fact reduction, and with some additional mdl data.
// ( ex. in case the same aspect is assigned twice, or twice at difference depths )
// for now the name of the field is the name of the aspect
func (op *PropertyAspect) ImportPropertyType(k *Importer, kind, prop ephemera.Named) (err error) {
	// record the existence of an aspect with the same name as the property
	k.NewName(prop.String(), tables.NAMED_ASPECT, op.At.String())
	// record the use of that property and aspect.
	k.NewField(kind, prop, tables.PRIM_ASPECT)
	return
}

// "{a number%number}, {some text%text}, or {a true/false value%bool}");
// bool properties become implicit aspects
func (op *PrimitiveType) ImportPropertyType(k *Importer, kind, prop ephemera.Named) (err error) {
	if op.Str != "$BOOL" {
		if primType, e := op.ImportPrimType(k); e != nil {
			err = e
		} else {
			k.NewField(kind, prop, primType)
		}
	} else {
		// ex. innumerable, not innumerable, is innumerable
		// there is an aspect "innumerable
		aspect := prop.String()
		aspectUpper := strings.Title(aspect) // manually camel aspects
		k.NewImplicitAspect(aspect, kind.String(),
			"not"+aspectUpper, // false first
			"is"+aspectUpper,
		)
		k.NewField(kind, prop, tables.PRIM_ASPECT)
	}
	return
}

// number_list, text_list, record_type, record_list
func (op *ExtType) ImportPropertyType(k *Importer, kind, prop ephemera.Named) (err error) {
	type primTypeImporter interface {
		ImportPrimType(*Importer) (string, error)
	}
	if imp, ok := op.Opt.(primTypeImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	} else if primType, e := imp.ImportPrimType(k); e != nil {
		err = e
	} else {
		k.NewField(kind, prop, primType)
	}
	return
}

func (op *NumberList) ImportPrimType(k *Importer) (ret string, err error) {
	ret = affine.NumList.String()
	return
}
func (op *TextList) ImportPrimType(k *Importer) (ret string, err error) {
	ret = affine.TextList.String()
	return
}
func (op *RecordType) ImportPrimType(k *Importer) (ret string, err error) {
	ret = op.Kind.Str
	return
}
func (op *RecordList) ImportPrimType(k *Importer) (ret string, err error) {
	// fix: this is fake: probably NewField needs both affinity and type.
	ret = "*" + op.Kind.Str
	return
}
