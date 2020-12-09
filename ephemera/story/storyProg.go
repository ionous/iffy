package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/rt"
)

type programImporter interface {
	ImportProgram(k *Importer) (ret programHook, err error)
}

func (op *ProgramHook) ImportProgram(k *Importer) (ret programHook, err error) {
	if opt, ok := op.Opt.(programImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	} else {
		ret, err = opt.ImportProgram(k)
	}
	return
}

// exists for formating in the composer; maps straight to program result.
func (op *ProgramReturn) ImportProgram(k *Importer) (ret programHook, err error) {
	return op.Result.ImportProgram(k)
}

// swaps b/t primitive and object functions
func (op *ProgramResult) ImportProgram(k *Importer) (ret programHook, err error) {
	if hook, ok := op.Opt.(programImporter); !ok {
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	} else {
		ret, err = hook.ImportProgram(k)
	}
	return
}

func (op *Activity) ImportProgram(k *Importer) (ret programHook, err error) {
	ret = &executeSlot{&core.Activity{op.Exe}}
	return
}

// swaps b/t different evals:
func (op *PrimitiveFunc) ImportProgram(k *Importer) (ret programHook, err error) {
	switch opt := op.Opt.(type) {
	case *NumberEval:
		ret = &numberSlot{rt.NumberEval(*opt)}
	case *TextEval:
		ret = &textSlot{rt.TextEval(*opt)}
	case *BoolEval:
		ret = &boolSlot{rt.BoolEval(*opt)}
	default:
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	}
	return
}

func (op *ObjectFunc) ImportProgram(k *Importer) (ret programHook, err error) {
	// FIX: we should wrap the text with a runtime "test the object matches" command
	// and, -- for simple text values -- add ephemera for "name, type"
	panic("not ported")
}
