package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
)

func (op *PatternRules) ImportPattern(k *Importer, patternName ephemera.Named) (err error) {
	if els := op.PatternRule; els != nil {
		for _, el := range *els {
			if e := el.ImportPattern(k, patternName); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

func (op *PatternRule) ImportPattern(k *Importer, patternName ephemera.Named) (err error) {
	if hook, e := op.Hook.ImportProgram(k); e != nil {
		err = e
	} else {
		name, rule := hook.NewRule(op.Guard)
		if patternProg, e := k.NewGob(name, rule); e != nil {
			err = e
		} else {
			k.NewPatternRule(patternName, patternProg)
		}
	}
	return
}

func (op *PatternLocals) ImportPattern(k *Importer, patternName ephemera.Named) (err error) {
	if els := op.LocalDecl; els != nil {
		for _, el := range *els {
			if e := el.ImportPattern(k, patternName); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

func (op *LocalDecl) ImportPattern(k *Importer, patternName ephemera.Named) (err error) {
	// fix: not implemented
	return
}

func (op *PatternedActivity) ImportActivity(k *Importer) (ret ephemera.Named, err error) {
	ret = k.NewName("execute", tables.NAMED_TYPE, op.At.String())
	return
}

func (op *PatternType) ImportType(k *Importer) (ret ephemera.Named, err error) {
	switch opt := op.Opt.(type) {
	case *PatternedActivity:
		ret, err = opt.ImportActivity(k)
	case *VariableType:
		ret, _, err = opt.ImportVariableType(k)
	default:
		err = ImportError(op, op.At, errutil.Fmt("%w for %T", UnhandledSwap, op.Opt))
	}
	return
}
