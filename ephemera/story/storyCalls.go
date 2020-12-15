package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

func (op *DetermineAct) ImportStub(k *Importer) (ret interface{}, err error) {
	if p, args, e := importCall(k, "execute", op.Name, op.Arguments); e != nil {
		err = ImportError(op, op.At, e)
	} else {
		ret = &pattern.DetermineAct{Pattern: p.String(), Arguments: args}
	}
	return
}
func (op *DetermineNum) ImportStub(k *Importer) (ret interface{}, err error) {
	if p, args, e := importCall(k, "number_eval", op.Name, op.Arguments); e != nil {
		err = ImportError(op, op.At, e)
	} else {
		ret = &pattern.DetermineNum{Pattern: p.String(), Arguments: args}
	}
	return
}
func (op *DetermineText) ImportStub(k *Importer) (ret interface{}, err error) {
	if p, args, e := importCall(k, "text_eval", op.Name, op.Arguments); e != nil {
		err = ImportError(op, op.At, e)
	} else {
		ret = &pattern.DetermineText{Pattern: p.String(), Arguments: args}
	}
	return
}
func (op *DetermineBool) ImportStub(k *Importer) (ret interface{}, err error) {
	if p, args, e := importCall(k, "bool_eval", op.Name, op.Arguments); e != nil {
		err = ImportError(op, op.At, e)
	} else {
		ret = &pattern.DetermineBool{Pattern: p.String(), Arguments: args}
	}
	return
}
func (op *DetermineNumList) ImportStub(k *Importer) (ret interface{}, err error) {
	if p, args, e := importCall(k, "num_list_eval", op.Name, op.Arguments); e != nil {
		err = ImportError(op, op.At, e)
	} else {
		ret = &pattern.DetermineNumList{Pattern: p.String(), Arguments: args}
	}
	return
}
func (op *DetermineTextList) ImportStub(k *Importer) (ret interface{}, err error) {
	if p, args, e := importCall(k, "text_list_eval", op.Name, op.Arguments); e != nil {
		err = ImportError(op, op.At, e)
	} else {
		ret = &pattern.DetermineTextList{Pattern: p.String(), Arguments: args}
	}
	return
}

func importCall(k *Importer, slot string, n PatternName, stubs *Arguments) (retName ephemera.Named, retArgs *core.Arguments, err error) {
	if p, e := n.NewName(k); e != nil {
		err = e
	} else if args, e := importArgs(k, p, stubs); e != nil {
		err = e
	} else {
		// fix: tests expect pattern type to be declared last :'(
		// fix: object type names will need adaption of some sort re plural_kinds
		patternType := k.NewName(slot, tables.NAMED_TYPE, n.At.String())
		k.NewPatternRef(p, p, patternType, "")
		retName, retArgs = p, args
	}
	return
}

func importArgs(k *Importer, p ephemera.Named, stubs *Arguments) (ret *core.Arguments, err error) {
	if stubs != nil {
		var argList []*core.Argument
		for _, stub := range stubs.Args {
			// FIX? GetEval should return affinity and type instead...
			if slotName, e := slotName(stub.From.GetEval()); e != nil {
				err = e
			} else if paramName, e := stub.Name.NewName(k, tables.NAMED_ARGUMENT); e != nil {
				err = errutil.Append(err, e)
			} else {
				if len(slotName) > 0 {
					paramType := k.NewName(slotName, tables.NAMED_TYPE, stub.At.String())
					k.NewPatternRef(p, paramName, paramType, "")
				}
				// after recording the "fact" of the parameter...
				// copy the stubbed argument data into the real argument list.
				newArg := &core.Argument{Name: paramName.String(), From: stub.From}
				argList = append(argList, newArg)
			}
		}
		if err == nil {
			ret = &core.Arguments{Args: argList}
		}
	}
	return
}
