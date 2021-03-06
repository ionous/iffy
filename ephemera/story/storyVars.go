package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/tables"
)

type variableDecl struct {
	name, typeName ephemera.Named
	affinity       string
}

func (op *VariableDecl) ImportVariable(k *Importer, cat string) (ret variableDecl, err error) {
	if n, e := op.Name.NewName(k, cat); e != nil {
		err = e
	} else if t, aff, e := op.Type.ImportVariableType(k); e != nil {
		err = e
	} else {
		ret = variableDecl{n, t, aff}
	}
	return
}

// primitive type, object type, or ext
func (op *VariableType) ImportVariableType(k *Importer) (retType ephemera.Named, retAff string, err error) {
	type variableTypeImporter interface {
		ImportVariableType(*Importer) (ephemera.Named, string, error)
	}
	if opt, ok := op.Opt.(variableTypeImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	} else {
		retType, retAff, err = opt.ImportVariableType(k)
	}
	return
}

func (op *ObjectType) ImportVariableType(k *Importer) (retType ephemera.Named, retAff string, err error) {
	retType, err = op.Kind.NewName(k)
	retAff = affine.Object.String()
	return
}

func (op *PrimitiveType) ImportPrimType(k *Importer) (ret string, err error) {
	if str, ok := decode.FindChoice(op, op.Str); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w %q", InvalidValue, op.Str))
	} else {
		ret = str
	}
	return
}

// returns one of the evalType(s) as a "Named" value --
// we return a name to normalize references to object kinds which are also used as variables
func (op *PrimitiveType) ImportVariableType(k *Importer) (retType ephemera.Named, retAff string, err error) {
	// fix -- shouldnt this be a different type ??
	// ie. we should be able to use FindChoie here.
	var namedType string
	switch str := op.Str; str {
	case "$NUMBER":
		namedType = "number_eval"
	case "$TEXT":
		namedType = "text_eval"
	case "$BOOL":
		namedType = "bool_eval"
	default:
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", InvalidValue, str))
	}
	if err == nil {
		retType = k.NewName(namedType, tables.NAMED_TYPE, op.At.String())
	}
	return
}
