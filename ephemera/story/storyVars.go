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

func (op *VariableType) ImportVariableType(k *Importer) (retType ephemera.Named, retAffinity string, err error) {
	switch opt := op.Opt.(type) {
	case *PrimitiveType:
		retType, err = opt.ImportEval(k)
	case *ObjectType:
		retType, err = opt.ImportType(k)
		retAffinity = affine.Object.String()
	default:
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	}
	return
}

func (op *ObjectType) ImportType(k *Importer) (ret ephemera.Named, err error) {
	// An   An
	return op.Kind.NewName(k)
}

func (op *PrimitiveType) ImportPrimType(k *Importer) (ret string, err error) {
	if str, ok := decode.FindChoice(op, op.Str); !ok || len(str) == 0 {
		err = ImportError(op, op.At, errutil.Fmt("%w %q", InvalidValue, op.Str))
	} else {
		ret = str
	}
	return
}

// returns one of the evalType(s) as a "Named" value --
// we return a name to normalize references to object kinds which are also used as variables
func (op *PrimitiveType) ImportEval(k *Importer) (ret ephemera.Named, err error) {
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
		ret = k.NewName(namedType, tables.NAMED_TYPE, op.At.String())
	}
	return
}
