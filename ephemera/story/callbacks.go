package story

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

func imp_determine_act(k *imp.Porter, m reader.Map) (ret interface{}, err error) {
	var from core.DetermineAct
	if e := fromPattern(k, m, &from, (*core.FromPattern)(&from), "execute"); e != nil {
		err = e
	} else {
		ret = &from
	}
	return
}
func imp_determine_num(k *imp.Porter, m reader.Map) (ret interface{}, err error) {
	var from core.DetermineNum
	if e := fromPattern(k, m, &from, (*core.FromPattern)(&from), "number_eval"); e != nil {
		err = e
	} else {
		ret = &from
	}
	return
}
func imp_determine_text(k *imp.Porter, m reader.Map) (ret interface{}, err error) {
	var from core.DetermineText
	if e := fromPattern(k, m, &from, (*core.FromPattern)(&from), "text_eval"); e != nil {
		err = e
	} else {
		ret = &from
	}
	return
}
func imp_determine_bool(k *imp.Porter, m reader.Map) (ret interface{}, err error) {
	var from core.DetermineBool
	if e := fromPattern(k, m, &from, (*core.FromPattern)(&from), "bool_eval"); e != nil {
		err = e
	} else {
		ret = &from
	}
	return
}
func imp_determine_num_list(k *imp.Porter, m reader.Map) (ret interface{}, err error) {
	var from core.DetermineNumList
	if e := fromPattern(k, m, &from, (*core.FromPattern)(&from), "num_list_eval"); e != nil {
		err = e
	} else {
		ret = &from
	}
	return
}
func imp_determine_text_list(k *imp.Porter, m reader.Map) (ret interface{}, err error) {
	var from core.DetermineTextList
	if e := fromPattern(k, m, &from, (*core.FromPattern)(&from), "text_list_eval"); e != nil {
		err = e
	} else {
		ret = &from
	}
	return
}

// from and spec are the same object, but go-type unpacking from interfaces isnt always easy.
func fromPattern(k *imp.Porter,
	m reader.Map,
	spec composer.Specification,
	from *core.FromPattern,
	evalType string) (err error) {
	typeName := spec.Compose().Name
	if m, e := reader.Unpack(m, typeName); e != nil {
		err = e
	} else if pid, e := fromPatternName(k, m.MapOf("$NAME"), typeName); e != nil {
		err = e
	} else if ps, e := imp_parameters(k, pid, m.MapOf("$PARAMETERS")); e != nil {
		err = e
	} else {
		n := k.NewName(tables.NAMED_TYPE, evalType, reader.At(m))
		// fix: object type names will need adaption of some sort re plural_kinds
		k.NewPatternType(pid, n)
		// assign results
		from.Name = pid.String()
		from.Parameters = ps
	}
	return
}

func imp_parameters(k *imp.Porter, pid ephemera.Named, r reader.Map) (ret *core.Parameters, err error) {
	if m, e := reader.Unpack(r, "parameters"); e != nil {
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

func imp_parameter(k *imp.Porter, pid ephemera.Named, r reader.Map) (ret *core.Parameter, err error) {
	if m, e := reader.Unpack(r, "parameter"); e != nil {
		err = e
	} else if n, e := imp_variable_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else if a, e := imp_assignment(k, m.MapOf("$FROM")); e != nil {
		err = e
	} else if evalType, e := assignmentType(a); e != nil {
		err = e
	} else {
		k.NewPatternParam(pid, n, k.NewName(tables.NAMED_TYPE, evalType, reader.At(m)))
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
		ret = "number_eval"
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

func imp_assignment(k *imp.Porter, m reader.Map) (ret core.Assignment, err error) {
	if a, e := k.DecodeSlot(m, "assignment"); e != nil {
		err = e
	} else if a, ok := a.(core.Assignment); !ok {
		err = errutil.Fmt("unexpected assignment %T", a)
	} else {
		ret = a
	}
	return
}

func imp_text(k *imp.Porter, m reader.Map) (ret string, err error) {
	return reader.String(m, "text")
}

// similar to imp_pattern_name, only we change the underlying type to indicate it's a reference.
func fromPatternName(k *imp.Porter, m reader.Map, typeName string) (ret ephemera.Named, err error) {
	if id, e := reader.Type(m, "pattern_name"); e != nil {
		err = e
	} else {
		ret = k.NewName(typeName, m.StrOf(reader.ItemValue), id)
	}
	return
}
