package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

func (op *PropertyPhrase) ImportProperties(k *Importer, kind ephemera.Named) (err error) {
	if imp, ok := op.Opt.(PropertyImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	} else {
		err = imp.ImportProperty(k, kind)
	}
	return
}

type PropertyImporter interface {
	ImportProperty(*Importer, ephemera.Named) error
}

func (op *PrimitivePhrase) ImportProperty(k *Importer, kind ephemera.Named) (err error) {
	if prop, e := op.Property.NewName(k); e != nil {
		err = e
	} else if prim, e := op.PrimitiveType.ImportPrim(k); e != nil {
		err = e
	} else {
		k.NewField(kind, prop, prim)
	}
	return
}

func (op *AspectPhrase) ImportProperty(k *Importer, kind ephemera.Named) (err error) {
	if aspect, e := op.Aspect.NewName(k); e != nil {
		err = e
	} else {
		// inform gives these the name "<noun> condition"
		// we could only do that with an after the fact reduction, and with some additional mdl data.
		// ( ex. in case the same aspect is assigned twice, or twice at difference depths )
		// for now the name of the field is the name of the aspect
		if op.OptionalProperty != nil {
			err = ImportError(op, op.At, errutil.Fmt("%w optional property names aren't supported for aspects", InvalidValue))
		} else {
			k.NewField(kind, aspect, tables.PRIM_ASPECT)
		}
	}
	return
}
