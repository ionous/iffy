package story

import (
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/tables"
)

func (op *PatternActions) ImportPhrase(k *Importer) (err error) {
	if patternName, e := op.Name.NewName(k); e != nil {
		err = e
	} else if e := op.PatternRules.ImportPattern(k, patternName); e != nil {
		err = e
	} else if e := op.PatternReturn.ImportReturn(k, patternName); e != nil {
		err = e
	} else {
		// import each local if they exist
		if els := op.PatternLocals; els != nil {
			err = els.ImportPattern(k, patternName)
		}
	}
	return
}

func (op *PatternReturn) ImportReturn(k *Importer, patternName ephemera.Named) (err error) {
	if op != nil { // pattern returns are optional
		if val, e := op.Result.ImportVariable(k, tables.NAMED_RETURN); e != nil {
			err = errutil.Append(err, e)
		} else {
			k.NewPatternDecl(patternName, val.name, val.typeName, val.affinity)
		}
	}
	return
}

// Adds a new pattern declaration and optionally some associated pattern parameters.
func (op *PatternDecl) ImportPhrase(k *Importer) (err error) {
	if patternName, e := op.Name.NewName(k); e != nil {
		err = e
	} else if patternType, e := op.Type.ImportType(k); e != nil {
		err = e
	} else {
		k.NewPatternDecl(patternName, patternName, patternType, "")
		//
		if els := op.Optvars; els != nil {
			for _, el := range els.VariableDecl {
				if val, e := el.ImportVariable(k, tables.NAMED_PARAMETER); e != nil {
					err = errutil.Append(err, e)
				} else {
					k.NewPatternDecl(patternName, val.name, val.typeName, val.affinity)
				}
			}
		}
	}
	return
}

func (op *PatternVariablesDecl) ImportPhrase(k *Importer) (err error) {
	if patternName, e := op.PatternName.NewName(k); e != nil {
		err = e
	} else {
		// fix: shouldnt this be called pattern parameters?
		for _, el := range op.VariableDecl {
			if val, e := el.ImportVariable(k, tables.NAMED_PARAMETER); e != nil {
				err = errutil.Append(err, e)
			} else {
				k.NewPatternDecl(patternName, val.name, val.typeName, val.affinity)
			}
		}
	}
	return
}
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
	} else if flags, e := op.Flags.ReadFlags(); e != nil {
		err = e
	} else if slotType := hook.SlotType(); flags != pattern.Terminal && slotType != "execute" && !strings.HasSuffix(slotType, "_list") {
		err = errutil.New("didnt expect continuation flags for", slotType, "in", patternName.String())
	} else {
		guard := op.Guard
		if dom := k.Current.Domain.String(); len(dom) > 0 {
			guard = &core.AllTrue{[]rt.BoolEval{
				&core.HasDominion{dom},
				guard,
			}}
		}
		ruleType, rule := hook.NewRule(guard, flags)
		if patternProg, e := k.NewGob(ruleType, rule); e != nil {
			err = e
		} else {
			k.NewPatternRule(patternName, patternProg)
		}
	}
	return
}

func (op *PatternFlags) ReadFlags() (ret pattern.Flags, err error) {
	if op != nil {
		switch str := op.Str; str {
		case "$BEFORE":
			// run other matching patterns, and then run this pattern. other...this.
			ret = pattern.Postfix
		case "$AFTER":
			// keep going after running the current pattern. this...others.
			ret = pattern.Prefix
		case "$TERMINATE":
			ret = pattern.Terminal
		default:
			err = errutil.New("unknown pattern flags", str)
		}
	}
	return
}

func (op *PatternLocals) ImportPattern(k *Importer, patternName ephemera.Named) (err error) {
	// fix: shouldnt this be called pattern parameters?
	for _, el := range op.VariableDecl {
		if val, e := el.ImportVariable(k, tables.NAMED_LOCAL); e != nil {
			err = errutil.Append(err, e)
		} else {
			k.NewPatternDecl(patternName, val.name, val.typeName, val.affinity)
		}
	}
	return
}

// returns "prog" as the name of a type  ( eases the difference b/t user named kinds, and internally named types )
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
