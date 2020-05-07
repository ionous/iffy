package internal

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

// func init() {
// 	// "number_eval" string
// 	RegisterCallback((*core.DetermineNum)(nil), nil)
// }

// this is an imperative command. so it doesnt have normal parsing.
// alternatively, make imperative parsing return an interface{}
// or, if everything was a statemachine -- the functions would be like events
// and each would allow its own parameters, etc.
func imp_determine_num(k *Importer, m reader.Map) (ret interface{}, err error) {
	if from, e := fromPattern(k, m, "determine_num", "number_eval"); e != nil {
		err = e
	} else {
		ret = (*core.DetermineNum)(from)
	}
	return
}

func fromPattern(k *Importer, r reader.Map, typeName, evalType string) (ret *core.FromPattern, err error) {
	if m, e := reader.Slat(r, typeName); e != nil {
		err = e
	} else if pid, e := fromPatternName(k, m.MapOf("$NAME"), typeName); e != nil {
		err = e
	} else if ps, e := imp_parameters(k, pid, m.MapOf("parameters")); e != nil {
		err = e
	} else if n, e := importName(k, m, evalType, tables.NAMED_TYPE); e != nil {
		err = e
	} else {
		// fix: object type names will need adaption of some sort re plural_kinds
		k.NewPatternType(pid, n)
		ret = &core.FromPattern{pid.String(), ps}
	}
	return
}

func imp_parameters(k *Importer, pid ephemera.Named, r reader.Map) (ret *core.Parameters, err error) {
	if m, e := reader.Slat(r, "parameters"); e != nil {
		err = e
	} else {
		var ps []*core.Parameter
		// record the referenced parameters: names, types, pairings
		if e := reader.Repeats(m.SliceOf("$PARAMS"),
			func(m reader.Map) (err error) {
				if p, e := imp_parameter(k, pid, m); e != nil {
					err = e
				} else {
					ps = append(ps, p)
				}
				return
			}); e != nil {
			err = e
		} else {
			ret = &core.Parameters{ps}
		}
	}
	return
}

func imp_parameter(k *Importer, pid ephemera.Named, r reader.Map) (ret *core.Parameter, err error) {
	if m, e := reader.Slat(r, "parameter"); e != nil {
		err = e
	} else if n, e := imp_variable_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else if a, e := imp_assignment(k, m.MapOf("$FROM")); e != nil {
		err = e
	} else if evalType, e := assignmentType(a); e != nil {
		err = e
	} else {
		k.NewPatternParam(pid, n, k.Named(tables.NAMED_TYPE, evalType, reader.At(m)))
		ret = &core.Parameter{Name: n.String(), From: a}
	}
	return
}

// return the evaluation type name underlying the parameter assignment
func assignmentType(a core.Assignment) (ret string, err error) {
	switch a.(type) {
	case *core.FromBool:
		ret = "bool_eval"
	case *core.FromNum:
		ret = "num_eval"
	case *core.FromText:
		ret = "text_eval"
	case *core.FromNumList:
		ret = "num_list_eval"
	case *core.FromTextList:
		ret = "text_list_eval"
	default:
		err = errutil.Fmt("unknown parameter type %T", a)
	}
	return
}

func imp_assignment(k *Importer, m reader.Map) (ret core.Assignment, err error) {
	if a, e := k.expectProg(m, "assignment"); e != nil {
		err = e
	} else if a, ok := a.(core.Assignment); !ok {
		err = errutil.Fmt("unexpected assignment %T", a)
	} else {
		ret = a
	}
	return
}

func imp_text(k *Importer, m reader.Map) (ret string, err error) {
	return reader.String(m, "text")
}

// similar to imp_pattern_name, only we change the underlying type to indicate it's a reference.
func fromPatternName(k *Importer, m reader.Map, typeName string) (ret ephemera.Named, err error) {
	if id, e := reader.Type(m, "pattern_name"); e != nil {
		err = e
	} else {
		ret = k.Named(typeName, m.StrOf(reader.ItemValue), id)
	}
	return
}
